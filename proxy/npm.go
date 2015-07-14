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

func (npmConfiguration NPMConfiguration) addProxySettings(proxyHost string, proxyPort string, nonProxyHosts []string) {
	npmConfiguration.removeProxySettings()
	for _, file := range npmConfiguration.Files {
		sanitisedPath := util.SanitisePath(file)
		contents := util.LoadFileIntoSlice(sanitisedPath)
		util.WriteSliceToFile(sanitisedPath, addNPMProxySettings(contents, proxyHost, proxyPort))
	}
}

func (npmConfiguration NPMConfiguration) removeProxySettings() {
	for _, file := range npmConfiguration.Files {
		sanitisedPath := util.SanitisePath(file)
		contents := util.LoadFileIntoSlice(sanitisedPath)
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
