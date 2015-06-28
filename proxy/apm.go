package proxy

import (
	"fmt"

	"github.com/andystanton/proxybastard/util"
)

// AddToAPM adds to apm.
func AddToAPM(config Configuration) {
	if config.Targets.APM.Enabled {
		util.ShellOut("apm", []string{"config", "set", "proxy", fmt.Sprintf("%s:%s", config.ProxyHost, config.ProxyPort)})
		util.ShellOut("apm", []string{"config", "set", "https-proxy", fmt.Sprintf("%s:%s", config.ProxyHost, config.ProxyPort)})
	}
}

// RemoveFromAPM removes from apm.
func RemoveFromAPM(config Configuration) {
	if config.Targets.APM.Enabled {
		currentHTTPProxy, err := util.ShellOut("apm", []string{"config", "get", "proxy"})

		if err == nil && currentHTTPProxy != "null" {
			util.ShellOut("apm", []string{"config", "delete", "proxy"})
		}

		currentHTTPSProxy, err := util.ShellOut("apm", []string{"config", "get", "https-proxy"})

		if err == nil && currentHTTPSProxy != "null" {
			util.ShellOut("apm", []string{"config", "delete", "https-proxy"})
		}
	}
}
