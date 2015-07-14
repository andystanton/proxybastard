package proxy

import (
	"fmt"

	"github.com/andystanton/proxybastard/util"
)

func (apmConfiguration APMConfiguration) isEnabled() bool {
	return apmConfiguration.Enabled
}

func (apmConfiguration APMConfiguration) addProxySettings(proxyHost string, proxyPort string, nonProxyHosts []string) {
	util.ShellOut("apm", []string{"config", "set", "http-proxy", fmt.Sprintf("%s:%s", proxyHost, proxyPort)})
	util.ShellOut("apm", []string{"config", "set", "https-proxy", fmt.Sprintf("%s:%s", proxyHost, proxyPort)})
	util.ShellOut("apm", []string{"config", "set", "strict-ssl", "false"})
}

func (apmConfiguration APMConfiguration) removeProxySettings() {
	currentHTTPProxy, err := util.ShellOut("apm", []string{"config", "get", "http-proxy"})

	if err == nil && currentHTTPProxy != "null" {
		util.ShellOut("apm", []string{"config", "delete", "http-proxy"})
	}

	currentHTTPSProxy, err := util.ShellOut("apm", []string{"config", "get", "https-proxy"})

	if err == nil && currentHTTPSProxy != "null" {
		util.ShellOut("apm", []string{"config", "delete", "https-proxy"})
	}

	currentStrictSSL, err := util.ShellOut("apm", []string{"config", "get", "strict-ssl"})

	if err == nil && currentStrictSSL != "null" {
		util.ShellOut("apm", []string{"config", "delete", "strict-ssl"})
	}
}
