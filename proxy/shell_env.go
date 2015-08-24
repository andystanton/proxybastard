package proxy

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/andystanton/proxybastard/util"
)

var proxyVars = []string{
	"HTTP_PROXY", "http_proxy",
	"HTTPS_PROXY", "https_proxy",
	"ALL_PROXY", "all_proxy",
}

var nonProxyVars = []string{
	"NO_PROXY", "no_proxy",
}

func addShellEnvVars(shellContents []string, proxyHost string, proxyPort string, nonProxyHosts []string) []string {
	updated := []string{}

	for _, shellLine := range shellContents {
		matched := false

		for _, proxyVar := range proxyVars {
			proxyRegex := regexp.MustCompile(fmt.Sprintf("^unset %s", proxyVar))
			matched = matched || proxyRegex.MatchString(shellLine)
		}

		for _, noProxyVar := range nonProxyVars {
			noProxyRegex := regexp.MustCompile(fmt.Sprintf("^unset %s", noProxyVar))
			matched = matched || noProxyRegex.MatchString(shellLine)
		}

		if !matched {
			updated = append(updated, shellLine)
		}
	}

	for _, proxyVar := range proxyVars {
		if proxyPort != "" {
			updated = append(updated, fmt.Sprintf("export %s=%s:%s", proxyVar, proxyHost, proxyPort))
		} else {
			updated = append(updated, fmt.Sprintf("export %s=%s", proxyVar, proxyHost))
		}
	}

	if len(nonProxyHosts) > 0 {
		for _, noProxyVar := range nonProxyVars {
			updated = append(updated, fmt.Sprintf("export %s=%s", noProxyVar, strings.Join(nonProxyHosts, ",")))
		}
	}

	return updated
}

func removeShellEnvVars(shellContents []string) []string {
	updated := []string{}

	unsetProxyVars := []string{}
	unsetNoProxyVars := []string{}

	for _, shellLine := range shellContents {
		matched := false

		for _, proxyVar := range proxyVars {
			proxyRegex := regexp.MustCompile(fmt.Sprintf("^export %s=[\\w:/.?&-]+$", proxyVar))
			matched = matched || proxyRegex.MatchString(shellLine)

			unsetRegex := regexp.MustCompile(fmt.Sprintf("^unset %s$", proxyVar))
			if unsetRegex.MatchString(shellLine) {
				unsetProxyVars = append(unsetProxyVars, proxyVar)
			}
		}

		for _, noProxyVar := range nonProxyVars {
			noProxyRegex := regexp.MustCompile(fmt.Sprintf("^export %s=[\\w:/,.?&-]+$", noProxyVar))
			matched = matched || noProxyRegex.MatchString(shellLine)

			unsetRegex := regexp.MustCompile(fmt.Sprintf("^unset %s$", noProxyVar))
			if unsetRegex.MatchString(shellLine) {
				unsetNoProxyVars = append(unsetNoProxyVars, noProxyVar)
			}
		}

		if !matched {
			updated = append(updated, shellLine)
		}
	}

	for _, proxyVar := range proxyVars {
		if !util.ContainsString(unsetProxyVars, proxyVar) {
			updated = append(updated, fmt.Sprintf("unset %s", proxyVar))
		}
	}

	for _, noProxyVar := range nonProxyVars {
		if !util.ContainsString(unsetNoProxyVars, noProxyVar) {
			updated = append(updated, fmt.Sprintf("unset %s", noProxyVar))
		}
	}

	return updated
}
