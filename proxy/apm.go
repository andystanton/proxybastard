package proxy

import (
	"fmt"

	"github.com/andystanton/proxybastard/util"
)

// addToAPM adds to apm.
func addToAPM(config Configuration) {
	if config.Targets.APM.Enabled {
		util.ShellOut("apm", []string{"config", "set", "http-proxy", fmt.Sprintf("%s:%s", config.ProxyHost, config.ProxyPort)})
		util.ShellOut("apm", []string{"config", "set", "https-proxy", fmt.Sprintf("%s:%s", config.ProxyHost, config.ProxyPort)})
		util.ShellOut("apm", []string{"config", "set", "strict-ssl", "false"})
	}
}

// removeFromAPM removes from apm.
func removeFromAPM(config Configuration) {
	if config.Targets.APM.Enabled {
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
}
