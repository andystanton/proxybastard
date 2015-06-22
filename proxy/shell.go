package proxy

import (
	"fmt"
	"regexp"
	"strings"
)

// ShellStatement is a line from a shell file. May contain multiple actual lines.
type ShellStatement struct {
	lines []string
}

var proxyVars = []string{"http_proxy", "https_proxy", "ALL_PROXY"}
var nonProxyVars = []string{"NO_PROXY"}

// ParseShellContents parses a string slice into a list of ShellStatements.
func ParseShellContents(shellContents []string) []ShellStatement {
	shellLines := []ShellStatement{}

	multiLine := regexp.MustCompile(`^(.*)\\$`)

	processingMultipleLines := false
	var currentShellLine ShellStatement

	for _, shellLine := range shellContents {
		currentLine := shellLine

		currentLineTrails := multiLine.MatchString(currentLine)
		if currentLineTrails {
			currentLine = multiLine.FindStringSubmatch(currentLine)[1]
		}

		if !processingMultipleLines {
			currentShellLine = ShellStatement{[]string{currentLine}}
		} else {
			currentShellLine.lines = append(currentShellLine.lines, currentLine)
		}

		if !currentLineTrails {
			shellLines = append(shellLines, currentShellLine)
		}

		processingMultipleLines = currentLineTrails
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
