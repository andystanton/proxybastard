package proxy

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/andystanton/proxybastard/util"
)

// SvnStatement represents a block in an svn file.
type SvnStatement struct {
	name  string
	lines []string
}

// AddToSubversion adds proxy settings to subversion files.
func AddToSubversion(config Configuration) {
	if config.Targets.Subversion.Enabled {
		RemoveFromSubversion(config)
		for _, svnFile := range config.Targets.Subversion.Files {
			sanitisedPath := util.SanitisePath(svnFile)
			contents := util.LoadFileIntoSlice(sanitisedPath)
			updated := addSubversionProxies(contents, config.ProxyHost, config.ProxyPort, config.NonProxyHosts)
			util.WriteSliceToFile(sanitisedPath, updated)
		}
	}
}

// RemoveFromSubversion removes proxy settings from subversion files.
func RemoveFromSubversion(config Configuration) {
	if config.Targets.Subversion.Enabled {
		for _, svnFile := range config.Targets.Subversion.Files {
			sanitisedPath := util.SanitisePath(svnFile)
			contents := util.LoadFileIntoSlice(sanitisedPath)
			updated := removeSubversionProxies(contents)
			util.WriteSliceToFile(sanitisedPath, updated)
		}
	}
}

func removeSubversionProxies(contents []string) []string {
	updated := []SvnStatement{}

	proxyRegex := regexp.MustCompile("^http-proxy-.+$")

	for _, statement := range ParseSubversionContents(contents) {
		if statement.name == "global" {
			newLines := []string{}
			for _, line := range statement.lines {
				if !proxyRegex.MatchString(line) {
					newLines = append(newLines, line)
				}
			}
			statement.lines = newLines
		}
		updated = append(updated, statement)
	}

	return ParseSubversionStatements(updated)
}

func addSubversionProxies(contents []string, proxyHost string, proxyPort string, nonProxyHosts []string) []string {
	updated := []SvnStatement{}

	proxyLines := []string{
		fmt.Sprintf("http-proxy-host=%s", proxyHost),
		fmt.Sprintf("http-proxy-port=%s", proxyPort),
		fmt.Sprintf("http-proxy-exceptions=%s", strings.Join(nonProxyHosts, ",")),
	}

	found := false
	for _, statement := range ParseSubversionContents(contents) {
		if statement.name == "global" {
			found = true
			statement.lines = append(statement.lines, proxyLines...)
		}
		updated = append(updated, statement)
	}
	if !found {
		updated = append(updated, SvnStatement{
			name:  "global",
			lines: proxyLines,
		})
	}

	return ParseSubversionStatements(updated)
}

// ParseSubversionContents parses the contents of a subversion file into a series of SvnStatements.
func ParseSubversionContents(contents []string) []SvnStatement {
	headingRegexp := regexp.MustCompile("^\\[(.+)\\]$")
	statements := []SvnStatement{}

	// state
	var currentSvnStatement SvnStatement
	inStatement := false

	for _, line := range contents {
		if headingRegexp.MatchString(line) {
			if inStatement {
				statements = append(statements, currentSvnStatement)
			}
			currentSvnStatement = SvnStatement{
				name:  headingRegexp.FindStringSubmatch(line)[1],
				lines: []string{},
			}
			inStatement = true
		} else if inStatement {
			if len(line) > 0 && line[0] != '#' {
				currentSvnStatement.lines = append(currentSvnStatement.lines, line)
			}
		}
	}

	if inStatement {
		statements = append(statements, currentSvnStatement)
	}

	return statements
}

// ParseSubversionStatements parses the contents of a subversion file into a series of SvnStatements.
func ParseSubversionStatements(statements []SvnStatement) []string {
	contents := []string{}

	for _, statement := range statements {
		contents = append(contents, fmt.Sprintf("[%s]", statement.name))
		contents = append(contents, statement.lines...)
		contents = append(contents, "")
	}

	return contents
}
