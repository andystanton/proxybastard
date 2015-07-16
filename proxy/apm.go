package proxy

import (
	"fmt"
	"regexp"

	"github.com/andystanton/proxybastard/util"
)

func (apmConfiguration APMConfiguration) isEnabled() bool {
	return apmConfiguration.Enabled
}

func (apmConfiguration APMConfiguration) validate() error {
	return nil
}

func (apmConfiguration APMConfiguration) suggestConfiguration() *Configuration {
	apmExecutable := "apm"
	apmFile := "~/.atom/.apmrc"
	apmFileSanitised := util.SanitisePath(apmFile)

	_, err := util.ShellOut("which", []string{apmExecutable})
	hasAPM := err == nil
	hasAPMRC := util.FileExists(apmFileSanitised)

	if hasAPM && hasAPMRC {

		contents, _ := util.LoadFileIntoSlice(apmFileSanitised)
		suggestedProxy, suggestedPort := extractProxyFromAPMContents(contents)

		return &Configuration{
			ProxyHost: suggestedProxy,
			ProxyPort: suggestedPort,
			Targets: &TargetsConfiguration{
				APM: &APMConfiguration{
					Enabled: true,
					Files:   []string{apmFile},
				},
			},
		}
	}
	return nil
}

func (apmConfiguration APMConfiguration) addProxySettings(proxyHost string, proxyPort string, nonProxyHosts []string) {
	apmConfiguration.removeProxySettings()
	for _, file := range apmConfiguration.Files {
		sanitisedPath := util.SanitisePath(file)
		contents, _ := util.LoadFileIntoSlice(sanitisedPath)
		util.WriteSliceToFile(sanitisedPath, addAPMProxySettings(contents, proxyHost, proxyPort))
	}
}

func (apmConfiguration APMConfiguration) removeProxySettings() {
	for _, file := range apmConfiguration.Files {
		sanitisedPath := util.SanitisePath(file)
		contents, _ := util.LoadFileIntoSlice(sanitisedPath)
		util.WriteSliceToFile(sanitisedPath, removeAPMProxySettings(contents))
	}
}

func removeAPMProxySettings(contents []string) []string {
	proxyRegexp := regexp.MustCompile("^(https?-proxy|strict-ssl)=.*$")
	outLines := []string{}
	for _, line := range contents {
		if !proxyRegexp.MatchString(line) {
			outLines = append(outLines, line)
		}
	}
	return outLines
}

func addAPMProxySettings(contents []string, proxyHost string, proxyPort string) []string {
	return append(contents, []string{
		fmt.Sprintf("http-proxy=%s:%s", proxyHost, proxyPort),
		fmt.Sprintf("https-proxy=%s:%s", proxyHost, proxyPort),
		"strict-ssl=false",
	}...)
}

func extractProxyFromAPMContents(contents []string) (string, string) {
	proxyRegexp := regexp.MustCompile("^https?-proxy=(.*)$")

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
