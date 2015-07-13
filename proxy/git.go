package proxy

import (
	"fmt"

	"github.com/andystanton/proxybastard/util"
)

// addToGit adds to Git.
func (gitConfiguration GitConfiguration) addProxySettings(proxyHost string, proxyPort string, nonProxyHosts []string) {
	if gitConfiguration.Enabled {
		util.ShellOut("git", []string{"config", "--global", "http.proxy", fmt.Sprintf("%s:%s", proxyHost, proxyPort)})
	}
}

// removeFromGit removes from Git.
func (gitConfiguration GitConfiguration) removeProxySettings() {
	if gitConfiguration.Enabled {
		current, err := util.ShellOut("git", []string{"config", "--global", "http.proxy"})

		if err == nil && current != "" {
			util.ShellOut("git", []string{"config", "--global", "--remove-section", "http"})
		}
	}
}
