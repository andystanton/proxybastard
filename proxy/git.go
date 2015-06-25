package proxy

import "github.com/andystanton/proxybastard/util"

// AddToGit adds to Git.
func AddToGit(proxyHost string, proxyPort string) {
	util.ShellOut("git", []string{"config", "--global", "http.proxy"})
}

// RemoveFromGit removes from Git.
func RemoveFromGit(proxyHost string, proxyPort string) {

}
