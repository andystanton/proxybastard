package proxy

import (
	"fmt"

	"github.com/andystanton/proxybastard/util"
)

// addToNPM adds to npm.
func addToNPM(config Configuration) {
	if config.Targets.NPM.Enabled {
		util.ShellOut("npm", []string{"config", "set", "proxy", fmt.Sprintf("%s:%s", config.ProxyHost, config.ProxyPort)})
		util.ShellOut("npm", []string{"config", "set", "https-proxy", fmt.Sprintf("%s:%s", config.ProxyHost, config.ProxyPort)})
	}
}

// removeFromNPM removes from npm.
func removeFromNPM(config Configuration) {
	if config.Targets.NPM.Enabled {
		currentHTTPProxy, err := util.ShellOut("npm", []string{"config", "get", "proxy"})

		if err == nil && currentHTTPProxy != "undefined" {
			util.ShellOut("npm", []string{"config", "delete", "proxy"})
		}

		currentHTTPSProxy, err := util.ShellOut("npm", []string{"config", "get", "https-proxy"})

		if err == nil && currentHTTPSProxy != "undefined" {
			util.ShellOut("npm", []string{"config", "delete", "https-proxy"})
		}
	}
}
