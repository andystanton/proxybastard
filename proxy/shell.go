package proxy

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/andystanton/proxybastard/util"
)

func (shellConfiguration ShellConfiguration) CustomPrompt(test string) string {
	return "shell configuration"
}

type shellStatement struct {
	lines []string
}

func (shellConfiguration ShellConfiguration) validate() error {
	return nil
}

func (shellConfiguration ShellConfiguration) isEnabled() bool {
	return shellConfiguration.Enabled
}

func (shellConfiguration ShellConfiguration) suggestConfiguration() *Configuration {
	shellFiles := []string{"~/.zshrc", "~/.bash_profile", "~/.bashrc", "~/.profile"}

	var shellFile string
	var hasShellFile bool
	for _, file := range shellFiles {
		if !hasShellFile && util.FileExists(util.SanitisePath(file)) {
			hasShellFile = true
			shellFile = file
		}
	}

	if hasShellFile {
		shellFileSanitised := util.SanitisePath(shellFile)
		contents, _ := util.LoadFileIntoSlice(shellFileSanitised)
		suggestedProxy, suggestedPort, suggestedNonProxyHosts, hasJavaOpts := extractProxyFromShellContents(contents)

		return &Configuration{
			ProxyHost:     suggestedProxy,
			ProxyPort:     suggestedPort,
			NonProxyHosts: suggestedNonProxyHosts,
			Targets: &TargetsConfiguration{
				Shell: &ShellConfiguration{
					Enabled:  true,
					JavaOpts: hasJavaOpts,
					Files:    []string{shellFile},
				},
			},
		}
	}
	return nil
}

func (shellConfiguration ShellConfiguration) addProxySettings(proxyHost string, proxyPort string, nonProxyHosts []string) {
	shellConfiguration.removeProxySettings()

	for _, shellFile := range shellConfiguration.Files {
		sanitisedPath := util.SanitisePath(shellFile)
		shellContents, _ := util.LoadFileIntoSlice(sanitisedPath)

		shellContents = addShellEnvVars(shellContents, proxyHost, proxyPort, nonProxyHosts)
		if shellConfiguration.JavaOpts {
			shellContents = addJavaOpts(shellContents, proxyHost, proxyPort, nonProxyHosts)
		}

		util.WriteSliceToFile(sanitisedPath, shellContents)
	}
}

func (shellConfiguration ShellConfiguration) removeProxySettings() {
	for _, shellFile := range shellConfiguration.Files {
		sanitisedPath := util.SanitisePath(shellFile)
		shellContents, _ := util.LoadFileIntoSlice(sanitisedPath)

		shellContents = removeShellEnvVars(shellContents)
		if shellConfiguration.JavaOpts {
			shellContents = removeJavaOpts(shellContents)
		}
		util.WriteSliceToFile(sanitisedPath, shellContents)
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

func extractProxyFromShellContents(contents []string) (string, string, []string, bool) {
	proxyRegexp := regexp.MustCompile("^export http_proxy=(.+)$")
	nphRegexp := regexp.MustCompile("^export NO_PROXY=(.+)")
	javaOptsRegexp := regexp.MustCompile("^export JAVA_OPTS=.*")

	var suggestedProxy string
	var suggestedPort string
	var suggestedNonProxyHosts []string
	var hasJavaOpts bool

	for _, line := range contents {
		proxyMatches := proxyRegexp.FindStringSubmatch(line)
		nphMatches := nphRegexp.FindStringSubmatch(line)
		hasJavaOpts = hasJavaOpts || javaOptsRegexp.MatchString(line)
		if len(proxyMatches) > 0 {
			hostRegexp := regexp.MustCompile("(.+):(.+)")
			hostMatches := hostRegexp.FindStringSubmatch(proxyMatches[1])
			if len(hostMatches) > 0 {
				suggestedProxy = hostMatches[1]
				suggestedPort = hostMatches[2]
			} else {
				suggestedProxy = proxyMatches[1]
			}
		} else if len(nphMatches) > 0 {
			suggestedNonProxyHosts = strings.Split(nphMatches[1], ",")
		}
	}
	return suggestedProxy, suggestedPort, suggestedNonProxyHosts, hasJavaOpts
}
