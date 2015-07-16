package proxy

import (
	"fmt"

	"github.com/andystanton/proxybastard/util"
)

func (gitConfiguration GitConfiguration) validate() error {
	return nil
}

func (gitConfiguration GitConfiguration) isEnabled() bool {
	return gitConfiguration.Enabled
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
