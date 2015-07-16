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

func (subversionConfiguration SubversionConfiguration) validate() error {
	return nil
}

func (subversionConfiguration SubversionConfiguration) isEnabled() bool {
	return subversionConfiguration.Enabled
}

func (subversionConfiguration SubversionConfiguration) suggestConfiguration() *Configuration {
	return nil
}

func (subversionConfiguration SubversionConfiguration) addProxySettings(proxyHost string, proxyPort string, nonProxyHosts []string) {
	subversionConfiguration.removeProxySettings()
	for _, svnFile := range subversionConfiguration.Files {
		sanitisedPath := util.SanitisePath(svnFile)
		contents, _ := util.LoadFileIntoSlice(sanitisedPath)
		updated := addSubversionProxies(contents, proxyHost, proxyPort, nonProxyHosts)
		util.WriteSliceToFile(sanitisedPath, updated)
	}
}

func (subversionConfiguration SubversionConfiguration) removeProxySettings() {
	for _, svnFile := range subversionConfiguration.Files {
		sanitisedPath := util.SanitisePath(svnFile)
		contents, _ := util.LoadFileIntoSlice(sanitisedPath)
		updated := removeSubversionProxies(contents)
		util.WriteSliceToFile(sanitisedPath, updated)
	}
}

func removeSubversionProxies(contents []string) []string {
	updated := []SvnStatement{}

	proxyRegex := regexp.MustCompile("^http-proxy-.+$")

	for _, statement := range parseSubversionContents(contents) {
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

	return parseSubversionStatements(updated)
}

func addSubversionProxies(contents []string, proxyHost string, proxyPort string, nonProxyHosts []string) []string {
	updated := []SvnStatement{}

	proxyProtocolMatch := regexp.MustCompile("^(?:https?://)?(.+)$")
	proxyHostWithoutProtocol := proxyHost
	if proxyProtocolMatch.MatchString(proxyHost) {
		proxyHostWithoutProtocol = proxyProtocolMatch.FindStringSubmatch(proxyHostWithoutProtocol)[1]
	}

	proxyLines := []string{
		fmt.Sprintf("http-proxy-host=%s", proxyHostWithoutProtocol),
		fmt.Sprintf("http-proxy-port=%s", proxyPort),
		fmt.Sprintf("http-proxy-exceptions=%s", strings.Join(nonProxyHosts, ",")),
	}

	found := false
	for _, statement := range parseSubversionContents(contents) {
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

	return parseSubversionStatements(updated)
}

func parseSubversionContents(contents []string) []SvnStatement {
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

func parseSubversionStatements(statements []SvnStatement) []string {
	contents := []string{}

	for _, statement := range statements {
		contents = append(contents, fmt.Sprintf("[%s]", statement.name))
		contents = append(contents, statement.lines...)
		contents = append(contents, "")
	}

	return contents
}
