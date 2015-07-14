package proxy

import (
	"fmt"

	"github.com/andystanton/proxybastard/util"
)

func (npmConfiguration NPMConfiguration) validate() error {
	return nil
}

func (npmConfiguration NPMConfiguration) isEnabled() bool {
	return npmConfiguration.Enabled
}

func (npmConfiguration NPMConfiguration) addProxySettings(proxyHost string, proxyPort string, nonProxyHosts []string) {
	util.ShellOut("npm", []string{"config", "set", "proxy", fmt.Sprintf("%s:%s", proxyHost, proxyPort)})
	util.ShellOut("npm", []string{"config", "set", "https-proxy", fmt.Sprintf("%s:%s", proxyHost, proxyPort)})
}

func (npmConfiguration NPMConfiguration) removeProxySettings() {
	currentHTTPProxy, err := util.ShellOut("npm", []string{"config", "get", "proxy"})

	if err == nil && currentHTTPProxy != "undefined" {
		util.ShellOut("npm", []string{"config", "delete", "proxy"})
	}

	currentHTTPSProxy, err := util.ShellOut("npm", []string{"config", "get", "https-proxy"})

	if err == nil && currentHTTPSProxy != "undefined" {
		util.ShellOut("npm", []string{"config", "delete", "https-proxy"})
	}
}
