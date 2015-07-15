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

func (apmConfiguration APMConfiguration) addProxySettings(proxyHost string, proxyPort string, nonProxyHosts []string) {
	apmConfiguration.removeProxySettings()
	for _, file := range apmConfiguration.Files {
		sanitisedPath := util.SanitisePath(file)
		contents, _ := util.LoadFileIntoSlice(sanitisedPath)
		util.SafeWriteSliceToFile(sanitisedPath, addAPMProxySettings(contents, proxyHost, proxyPort))
	}
}

func (apmConfiguration APMConfiguration) removeProxySettings() {
	for _, file := range apmConfiguration.Files {
		sanitisedPath := util.SanitisePath(file)
		contents, _ := util.LoadFileIntoSlice(sanitisedPath)
		util.SafeWriteSliceToFile(sanitisedPath, removeAPMProxySettings(contents))
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
