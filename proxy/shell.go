package proxy

import (
	"fmt"
	"regexp"
)

// ShellStatement is a line from a shell file. May contain multiple actual lines.
type ShellStatement struct {
	lines []string
}

// RemoveFromShell removes proxy entries from a shell file.
func RemoveFromShell(config Configuration) {
	for _, shellFile := range config.Targets.Shell.Files {
		sanitisedPath := TildeToUserHome(shellFile)
		writeSliceToFile(sanitisedPath, RemoveEnvVars(loadFileIntoSlice(sanitisedPath)))
	}
}

// AddToShell adds proxy entries to a shell file.
func AddToShell(config Configuration) {
	for _, shellFile := range config.Targets.Shell.Files {
		sanitisedPath := TildeToUserHome(shellFile)

		RemoveFromShell(config)

		shellContents := loadFileIntoSlice(sanitisedPath)

		shellContents = AddEnvVars(shellContents, config.ProxyHost, config.ProxyPort, config.NonProxyHosts)
		if config.Targets.Shell.JavaOpts {
			shellContents = AddJavaOpts(shellContents, config.ProxyHost, config.ProxyPort, config.NonProxyHosts)
		}

		writeSliceToFile(sanitisedPath, shellContents)
	}
}

// ParseShellStatements parses a slice of ShellStatements into a string slice.
func ParseShellStatements(shellStatements []ShellStatement) []string {
	contents := []string{}
	for _, statement := range shellStatements {
		for count, line := range statement.lines {
			if count == len(statement.lines)-1 {
				contents = append(contents, line)
			} else {
				contents = append(contents, fmt.Sprintf("%s \\", line))
			}
		}
	}
	return contents
}

// ParseShellContents parses a string slice into a list of ShellStatements.
func ParseShellContents(shellContents []string) []ShellStatement {
	shellLines := []ShellStatement{}

	trailingLineRegex := regexp.MustCompile(`^(.*)\s*\\$`)

	// State.
	processingMultipleLines := false
	var currentShellLine ShellStatement

	for _, shellLine := range shellContents {
		currentLine := shellLine

		currentLineTrails := trailingLineRegex.MatchString(currentLine)
		if currentLineTrails {
			currentLine = trailingLineRegex.FindStringSubmatch(currentLine)[1]
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
