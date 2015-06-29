package proxy

import (
	"fmt"

	"github.com/andystanton/proxybastard/util"
)

// addToGit adds to Git.
func addToGit(config Configuration) {
	if config.Targets.Git.Enabled {
		util.ShellOut("git", []string{"config", "--global", "http.proxy", fmt.Sprintf("%s:%s", config.ProxyHost, config.ProxyPort)})
	}
}

// removeFromGit removes from Git.
func removeFromGit(config Configuration) {
	if config.Targets.Git.Enabled {
		current, err := util.ShellOut("git", []string{"config", "--global", "http.proxy"})

		if err == nil && current != "" {
			util.ShellOut("git", []string{"config", "--global", "--remove-section", "http"})
		}
	}
}
