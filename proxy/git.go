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
		suggestedProxy, suggestedPort := extractProxyFromGit()
		return &Configuration{
			ProxyHost: suggestedProxy,
			ProxyPort: suggestedPort,
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

func extractProxyFromGit() (string, string) {
	var suggestedProxy string
	var suggestedPort string
	current, err := util.ShellOut("git", []string{"config", "--global", "http.proxy"})

	if err == nil && current != "" {
		hostRegexp := regexp.MustCompile("(.+):(.+)")
		hostMatches := hostRegexp.FindStringSubmatch(current)
		if len(hostMatches) > 0 {
			suggestedProxy = hostMatches[1]
			suggestedPort = hostMatches[2]
		} else {
			suggestedProxy = current
		}
	}

	return suggestedProxy, suggestedPort
}
