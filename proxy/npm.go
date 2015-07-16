package proxy

import (
	"fmt"
	"regexp"

	"github.com/andystanton/proxybastard/util"
)

func (npmConfiguration NPMConfiguration) validate() error {
	return nil
}

func (npmConfiguration NPMConfiguration) isEnabled() bool {
	return npmConfiguration.Enabled
}

func (npmConfiguration NPMConfiguration) suggestConfiguration() *Configuration {
	npmExecutable := "npm"
	npmFile := "~/.npmrc"
	npmFileSanitised := util.SanitisePath(npmFile)

	_, err := util.ShellOut("which", []string{npmExecutable})
	hasNPM := err == nil
	hasNPMRC := util.FileExists(npmFileSanitised)

	if hasNPM && hasNPMRC {

		contents, _ := util.LoadFileIntoSlice(npmFileSanitised)
		suggestedProxy, suggestedPort := extractProxyFromNPMContents(contents)

		return &Configuration{
			ProxyHost: suggestedProxy,
			ProxyPort: suggestedPort,
			Targets: &TargetsConfiguration{
				NPM: &NPMConfiguration{
					Enabled: true,
					Files:   []string{npmFile},
				},
			},
		}
	}
	return nil
}

func (npmConfiguration NPMConfiguration) addProxySettings(proxyHost string, proxyPort string, nonProxyHosts []string) {
	npmConfiguration.removeProxySettings()
	for _, file := range npmConfiguration.Files {
		sanitisedPath := util.SanitisePath(file)
		contents, _ := util.LoadFileIntoSlice(sanitisedPath)
		util.WriteSliceToFile(sanitisedPath, addNPMProxySettings(contents, proxyHost, proxyPort))
	}
}

func (npmConfiguration NPMConfiguration) removeProxySettings() {
	for _, file := range npmConfiguration.Files {
		sanitisedPath := util.SanitisePath(file)
		contents, _ := util.LoadFileIntoSlice(sanitisedPath)
		util.WriteSliceToFile(sanitisedPath, removeNPMProxySettings(contents))
	}
}

func removeNPMProxySettings(contents []string) []string {
	proxyRegexp := regexp.MustCompile("^(https-)?proxy=.*$")
	outLines := []string{}
	for _, line := range contents {
		if !proxyRegexp.MatchString(line) {
			outLines = append(outLines, line)
		}
	}
	return outLines
}

func addNPMProxySettings(contents []string, proxyHost string, proxyPort string) []string {
	return append(contents, []string{
		fmt.Sprintf("proxy=%s:%s", proxyHost, proxyPort),
		fmt.Sprintf("https-proxy=%s:%s", proxyHost, proxyPort),
	}...)
}

func extractProxyFromNPMContents(contents []string) (string, string) {
	proxyRegexp := regexp.MustCompile("^(?:https-)?proxy=(.*)$")

	var suggestestedProxy string
	var suggestedPort string

	for _, line := range contents {
		matches := proxyRegexp.FindStringSubmatch(line)
		if len(matches) > 0 {
			hostRegexp := regexp.MustCompile("(.+):(.+)")
			hostMatches := hostRegexp.FindStringSubmatch(matches[1])
			if len(hostMatches) > 0 {
				suggestestedProxy = hostMatches[1]
				suggestedPort = hostMatches[2]
			} else {
				suggestestedProxy = matches[1]
			}
			break
		}
	}
	return suggestestedProxy, suggestedPort
}
