package proxy

import (
	"fmt"
	"regexp"
	"strings"
)

// ShellLine is a line from a shell file. May contain multiple actual lines.
type ShellLine struct {
	lines []string
}

var proxyVars = []string{"http_proxy", "https_proxy", "ALL_PROXY"}
var nonProxyVars = []string{"NO_PROXY"}

// ParseShellFile parses a shell file.
func ParseShellFile(shellContents []string) []ShellLine {
	shellLines := []ShellLine{}

	multiLine := regexp.MustCompile(`^(.*)\\$`)

	inMultiLine := false
	var currentShellLine ShellLine

	for _, shellLine := range shellContents {
		newInMultiLine := multiLine.MatchString(shellLine)

		actualLine := shellLine
		if newInMultiLine {
			actualLine = multiLine.FindStringSubmatch(actualLine)[1]
		}

		if !inMultiLine {
			currentShellLine = ShellLine{[]string{actualLine}}
		} else {
			currentShellLine.lines = append(currentShellLine.lines, actualLine)
		}

		if !newInMultiLine {
			shellLines = append(shellLines, currentShellLine)
		}

		inMultiLine = newInMultiLine
	}

	return shellLines
}

// AddProxyVars adds shell vars to a file.
func AddProxyVars(shellContents []string, proxyHost string, proxyPort string, nonProxyHosts []string) []string {
	updated := shellContents

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

// RemoveProxyVars removes shell vars from a file.
func RemoveProxyVars(shellContents []string) []string {
	updated := []string{}

	for _, shellLine := range shellContents {
		matched := false

		for _, proxyVar := range proxyVars {
			proxyRegex := regexp.MustCompile(fmt.Sprintf("^export %s=[\\w:/.?&-]+$", proxyVar))
			matched = matched || proxyRegex.MatchString(shellLine)
		}

		for _, noProxyVar := range nonProxyVars {
			noProxyRegex := regexp.MustCompile(fmt.Sprintf("^export %s=[\\w:/,.?&-]+$", noProxyVar))
			matched = matched || noProxyRegex.MatchString(shellLine)
		}

		if !matched {
			updated = append(updated, shellLine)
		}
	}

	return updated
}
