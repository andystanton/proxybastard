package proxy

import (
	"fmt"
	"regexp"

	"github.com/andystanton/proxybastard/util"
)

type shellStatement struct {
	lines []string
}

func removeFromShell(config Configuration) {
	if config.Targets.Shell.Enabled {
		for _, shellFile := range config.Targets.Shell.Files {
			sanitisedPath := util.SanitisePath(shellFile)
			shellContents := util.LoadFileIntoSlice(sanitisedPath)

			shellContents = removeShellEnvVars(shellContents)
			if config.Targets.Shell.JavaOpts {
				shellContents = removeJavaOpts(shellContents)
			}
			util.WriteSliceToFile(sanitisedPath, shellContents)
		}
	}
}

func addToShell(config Configuration) {
	if config.Targets.Shell.Enabled {
		for _, shellFile := range config.Targets.Shell.Files {
			removeFromShell(config)

			sanitisedPath := util.SanitisePath(shellFile)
			shellContents := util.LoadFileIntoSlice(sanitisedPath)

			shellContents = addShellEnvVars(shellContents, config.ProxyHost, config.ProxyPort, config.NonProxyHosts)
			if config.Targets.Shell.JavaOpts {
				shellContents = addJavaOpts(shellContents, config.ProxyHost, config.ProxyPort, config.NonProxyHosts)
			}

			util.WriteSliceToFile(sanitisedPath, shellContents)
		}
	}
}

func parseShellStatements(shellStatements []shellStatement) []string {
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

func parseShellContents(shellContents []string) []shellStatement {
	shellLines := []shellStatement{}

	trailingLineRegex := regexp.MustCompile(`^(.*)\s*\\$`)

	// State.
	processingMultipleLines := false
	var currentShellLine shellStatement

	for _, shellLine := range shellContents {
		currentLine := shellLine

		currentLineTrails := trailingLineRegex.MatchString(currentLine)
		if currentLineTrails {
			currentLine = trailingLineRegex.FindStringSubmatch(currentLine)[1]
		}

		if !processingMultipleLines {
			currentShellLine = shellStatement{[]string{currentLine}}
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
