package proxy

import (
	"fmt"
	"regexp"

	"github.com/andystanton/proxybastard/util"
)

func (gitConfiguration GitConfiguration) validate() error {
	return nil
}

func (gitConfiguration GitConfiguration) isEnabled() bool {
	return gitConfiguration.Enabled
}

func (gitConfiguration GitConfiguration) suggestConfiguration() *Configuration {
	gitExecutable := "git"
	_, err := util.ShellOut("which", []string{gitExecutable})
	if err == nil {
		return &Configuration{
			Targets: &TargetsConfiguration{
				Git: &GitConfiguration{
					Enabled: true,
				},
			},
		}
	}
	return nil
}

func (gitConfiguration GitConfiguration) addProxySettings(proxyHost string, proxyPort string, nonProxyHosts []string) {
	util.ShellOut("git", []string{"config", "--global", "http.proxy", fmt.Sprintf("%s:%s", proxyHost, proxyPort)})
}

func (gitConfiguration GitConfiguration) removeProxySettings() {
	current, err := util.ShellOut("git", []string{"config", "--global", "http.proxy"})

	if err == nil && current != "" {
		util.ShellOut("git", []string{"config", "--global", "--remove-section", "http"})
	}
}

func extractProxyFromGit(contents []string) (string, string) {
	var suggestestedProxy string
	var suggestedPort string
	current, err := util.ShellOut("git", []string{"config", "--global", "http.proxy"})

	if err == nil && current != "" {
		hostRegexp := regexp.MustCompile("(.+):(.+)")
		hostMatches := hostRegexp.FindStringSubmatch(current)
		if len(hostMatches) > 0 {
			suggestestedProxy = hostMatches[1]
			suggestedPort = hostMatches[2]
		} else {
			suggestestedProxy = current
		}
	}

	return suggestestedProxy, suggestedPort
}
