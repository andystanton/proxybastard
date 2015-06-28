package proxy

import (
	"fmt"
	"regexp"

	"github.com/andystanton/proxybastard/util"
)

// ShellStatement is a line from a shell file. May contain multiple actual lines.
type ShellStatement struct {
	lines []string
}

// RemoveFromShell removes proxy entries from a shell file.
func RemoveFromShell(config Configuration) {
	if config.Targets.Shell.Enabled {
		for _, shellFile := range config.Targets.Shell.Files {
			sanitisedPath := util.SanitisePath(shellFile)
			shellContents := util.LoadFileIntoSlice(sanitisedPath)

			shellContents = RemoveEnvVars(shellContents)
			if config.Targets.Shell.JavaOpts {
				shellContents = RemoveJavaOpts(shellContents)
			}
			util.WriteSliceToFile(sanitisedPath, shellContents)
		}
	}
}

// AddToShell adds proxy entries to a shell file.
func AddToShell(config Configuration) {
	if config.Targets.Shell.Enabled {
		for _, shellFile := range config.Targets.Shell.Files {
			RemoveFromShell(config)

			sanitisedPath := util.SanitisePath(shellFile)
			shellContents := util.LoadFileIntoSlice(sanitisedPath)

			shellContents = AddEnvVars(shellContents, config.ProxyHost, config.ProxyPort, config.NonProxyHosts)
			if config.Targets.Shell.JavaOpts {
				shellContents = AddJavaOpts(shellContents, config.ProxyHost, config.ProxyPort, config.NonProxyHosts)
			}

			util.WriteSliceToFile(sanitisedPath, shellContents)
		}
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
				newLine := regexp.MustCompile("\" \\\\").ReplaceAllString(fmt.Sprintf("%s \\", line), "\"\\")
				contents = append(contents, newLine)
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
