package proxy

import (
	"fmt"
	"regexp"
)

var proxyVars = []string{"http_proxy", "https_proxy", "ALL_PROXY"}

// AddProxyVars adds shell vars to a file.
func AddProxyVars(shellContents []string, proxyHost string, proxyPort *int) []string {
	updated := shellContents

	for _, proxyVar := range proxyVars {
		if proxyPort != nil {
			updated = append(updated, fmt.Sprintf("export %s=%s:%d", proxyVar, proxyHost, *proxyPort))
		} else {
			updated = append(updated, fmt.Sprintf("export %s=%s", proxyVar, proxyHost))
		}
	}

	return updated
}

// RemoveProxyVars removes shell vars from a file.
func RemoveProxyVars(shellContents []string) []string {
	updated := []string{}

	for _, shellLine := range shellContents {
		matched := false

		for _, proxyVar := range proxyVars {
			proxyRegex := regexp.MustCompile(fmt.Sprintf("^export %s=[\\w:/.?&-]+$", proxyVar))
			matched = matched || proxyRegex.MatchString(shellLine)
		}

		if !matched {
			updated = append(updated, shellLine)
		}
	}

	return updated
}
