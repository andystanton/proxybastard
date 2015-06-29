package proxy

import (
	"fmt"
	"regexp"
	"strings"
)

func parseJavaOpts(javaOpts []string) []string {
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

func addJavaOpts(shellContents []string, proxyHost string, proxyPort string, nonProxyHosts []string) []string {
	shellStatements := parseShellContents(shellContents)

	javaOptRegex := regexp.MustCompile("^\\s*export\\s*JAVA_OPTS=.*")

	existingOpts := false
	var javaOptStatment shellStatement
	var javaOptIndex int

	for index, statement := range shellStatements {
		if javaOptRegex.MatchString(statement.lines[0]) {
			existingOpts = true
			javaOptStatment = statement
			javaOptIndex = index
		}
	}

	if !existingOpts {
		javaOptStatment = shellStatement{
			lines: []string{"export JAVA_OPTS=\"\""},
		}
	}

	parsedOpts := parseJavaOpts(javaOptStatment.lines)

	parsedOpts = append(parsedOpts, fmt.Sprintf("-Dhttp.proxyHost=%s", proxyHost))
	parsedOpts = append(parsedOpts, fmt.Sprintf("-Dhttp.proxyPort=%s", proxyPort))
	parsedOpts = append(parsedOpts, fmt.Sprintf("-Dhttps.proxyHost=%s", proxyHost))
	parsedOpts = append(parsedOpts, fmt.Sprintf("-Dhttps.proxyPort=%s", proxyPort))
	parsedOpts = append(parsedOpts, fmt.Sprintf("-Dhttp.nonProxyHosts=%s", strings.Join(nonProxyHosts, "|")))

	outputLines := []string{"export JAVA_OPTS=\""}
	outputLines = append(outputLines, parsedOpts...)
	outputLines[len(outputLines)-1] = regexp.MustCompile("$").ReplaceAllString(outputLines[len(outputLines)-1], "\"")

	javaOptStatment.lines = outputLines

	if existingOpts {
		shellStatements[javaOptIndex] = javaOptStatment
	} else {
		shellStatements = append(shellStatements, javaOptStatment)
	}

	return parseShellStatements(shellStatements)
}

func removeJavaOpts(shellContents []string) []string {
	shellStatements := parseShellContents(shellContents)

	javaOptRegex := regexp.MustCompile("^\\s*export\\s*JAVA_OPTS=.*")

	existingOpts := false
	var javaOptStatment shellStatement
	var javaOptIndex int

	for index, statement := range shellStatements {
		if javaOptRegex.MatchString(statement.lines[0]) {
			existingOpts = true
			javaOptStatment = statement
			javaOptIndex = index
		}
	}

	if existingOpts {
		parsedOpts := parseJavaOpts(javaOptStatment.lines)
		javaOptStatment.lines = []string{"export JAVA_OPTS=\""}

		proxyRegex := regexp.MustCompile("-Dhttps?.(proxyHost|proxyPort|nonProxyHosts)=.*")
		for _, opt := range parsedOpts {
			if !proxyRegex.MatchString(opt) {
				javaOptStatment.lines = append(javaOptStatment.lines, opt)
			}
		}
		javaOptStatment.lines[len(javaOptStatment.lines)-1] = regexp.MustCompile("$").ReplaceAllString(javaOptStatment.lines[len(javaOptStatment.lines)-1], "\"")
		shellStatements[javaOptIndex] = javaOptStatment
	}

	return parseShellStatements(shellStatements)
}
