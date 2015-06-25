package proxy

import (
	"fmt"

	"github.com/andystanton/proxybastard/util"
)

// AddToGit adds to Git.
func AddToGit(config Configuration) {
	util.ShellOut("git", []string{"config", "--global", "http.proxy", fmt.Sprintf("%s:%s", config.ProxyHost, config.ProxyPort)})
}

// RemoveFromGit removes from Git.
func RemoveFromGit(config Configuration) {
	current, err := util.ShellOut("git", []string{"config", "--global", "http.proxy"})

	if err == nil && current != "" {
		util.ShellOut("git", []string{"config", "--global", "--remove-section", "http"})
	}
}
