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

// RemoveFromShell removes proxy entries from a shell file.
func RemoveFromShell(filename string) {
	sanitisedPath := TildeToUserHome(filename)
	writeSliceToFile(sanitisedPath, RemoveEnvVars(loadFileIntoSlice(sanitisedPath)))
}

// AddToShell adds proxy entries to a shell file.
func AddToShell(filename string, config Configuration) {
	RemoveFromShell(filename)

	for _, shellFile := range config.Targets.Shell.Files {
		sanitisedPath := TildeToUserHome(shellFile)

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

// ParseJavaOpts tokenises a java_opts string.
func ParseJavaOpts(javaOpts []string) []string {
	firstLine := regexp.MustCompile("^\\s*export\\s*JAVA_OPTS=\"(.*)\\\\?$")
	interLine := regexp.MustCompile("^\\s*(.*)\\s*\\\\$")
	finalLine := regexp.MustCompile("(.*)\"$")

	optsMap := []string{}

	for _, statement := range javaOpts {
		lineContent := statement

		isFirstLine := firstLine.MatchString(lineContent)
		if isFirstLine {
			lineContent = firstLine.FindStringSubmatch(lineContent)[1]
		}

		isFinalLine := finalLine.MatchString(lineContent)
		if isFinalLine {
			lineContent = finalLine.FindStringSubmatch(lineContent)[1]
		}

		isInterLine := interLine.MatchString(lineContent)
		if isInterLine {
			lineContent = interLine.FindStringSubmatch(lineContent)[1]
		}

		sanitisedLineContent := strings.TrimSpace(lineContent)

		if len(sanitisedLineContent) > 0 {
			for _, token := range strings.Split(sanitisedLineContent, " ") {
				optsMap = append(optsMap, token)
			}
		}
	}

	return optsMap
}

// AddJavaOpts adds proxy settings to a JAVA_OPTS declaration in a shell file.
func AddJavaOpts(shellContents []string, proxyHost string, proxyPort string, nonProxyHosts []string) []string {
	shellStatements := ParseShellContents(shellContents)

	javaOptRegex := regexp.MustCompile("^\\s*export\\s*JAVA_OPTS=.*")

	existingOpts := false
	var javaOptStatment ShellStatement
	var javaOptIndex int

	for index, statement := range shellStatements {
		if javaOptRegex.MatchString(statement.lines[0]) {
			existingOpts = true
			javaOptStatment = statement
			javaOptIndex = index
		}
	}

	if !existingOpts {
		javaOptStatment = ShellStatement{
			lines: []string{"export JAVA_OPTS=\"\""},
		}
	}

	parsedOpts := ParseJavaOpts(javaOptStatment.lines)

	parsedOpts = append(parsedOpts, fmt.Sprintf("-Dhttp.proxyHost=%s", proxyHost))
	if proxyPort != "" {
		parsedOpts = append(parsedOpts, fmt.Sprintf("-Dhttp.proxyPort=%s", proxyPort))
	}
	parsedOpts = append(parsedOpts, fmt.Sprintf("-Dhttps.proxyHost=%s", proxyHost))
	if proxyPort != "" {
		parsedOpts = append(parsedOpts, fmt.Sprintf("-Dhttps.proxyPort=%s", proxyPort))
	}
	parsedOpts = append(parsedOpts, fmt.Sprintf("-Dhttp.nonProxyHosts=%s", strings.Join(nonProxyHosts, "|")))

	outputLines := []string{"export JAVA_OPTS=\""}
	outputLines = append(outputLines, parsedOpts...)
	outputLines[len(outputLines)-1] = regexp.MustCompile("$").ReplaceAllString(outputLines[len(outputLines)-1], "\"")
	for _, blah := range outputLines {
		fmt.Println(blah)
	}

	javaOptStatment.lines = outputLines

	if existingOpts {
		shellStatements[javaOptIndex] = javaOptStatment
	} else {
		shellStatements = append(shellStatements, javaOptStatment)
	}

	return ParseShellStatements(shellStatements)
}

// RemoveJavaOpts removes proxy settings to a JAVA_OPTS declaration in a shell file.
func RemoveJavaOpts(shellContents []string, proxyHost string, proxyPort string, nonProxyHosts []string) []string {
	shellStatements := ParseShellContents(shellContents)
	return ParseShellStatements(shellStatements)
}

// AddEnvVars adds shell vars to a file.
func AddEnvVars(shellContents []string, proxyHost string, proxyPort string, nonProxyHosts []string) []string {
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

// RemoveEnvVars removes shell vars from a file.
func RemoveEnvVars(shellContents []string) []string {
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
