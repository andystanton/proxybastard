package proxy

import (
	"fmt"
	"regexp"
	"strings"
)

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
	parsedOpts = append(parsedOpts, fmt.Sprintf("-Dhttp.proxyPort=%s", proxyPort))
	parsedOpts = append(parsedOpts, fmt.Sprintf("-Dhttps.proxyHost=%s", proxyHost))
	parsedOpts = append(parsedOpts, fmt.Sprintf("-Dhttps.proxyPort=%s", proxyPort))
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
