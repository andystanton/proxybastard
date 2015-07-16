package proxy

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/andystanton/proxybastard/util"
)

type sshHost struct {
	Pattern    string
	Statements []string
}

type sshFile struct {
	GlobalStatements []string
	Hosts            []sshHost
}

func (sshConfiguration SSHConfiguration) validate() error {
	return nil
}

func (sshConfiguration SSHConfiguration) isEnabled() bool {
	return sshConfiguration.Enabled
}

func (sshConfiguration SSHConfiguration) bringMeTheThing() interface{} {
	return nil
}

func (sshConfiguration SSHConfiguration) addSOCKSProxySettings(socksProxyHost string, socksProxyPort string) {
	sshConfiguration.removeSOCKSProxySettings()

	for _, sshConfig := range sshConfiguration.Files {
		sanitisedPath := util.SanitisePath(sshConfig)
		contents, _ := util.LoadFileIntoSlice(sanitisedPath)
		util.WriteSliceToFile(sanitisedPath, addSSHConfig(contents, socksProxyHost, socksProxyPort))
	}
}

func (sshConfiguration SSHConfiguration) removeSOCKSProxySettings() {
	for _, sshConfig := range sshConfiguration.Files {
		sanitisedPath := util.SanitisePath(sshConfig)
		contents, _ := util.LoadFileIntoSlice(sanitisedPath)
		util.WriteSliceToFile(sanitisedPath, removeSSHConfig(contents))
	}
}

func removeFromStatements(original []string) []string {
	proxyRegex := regexp.MustCompile("^ProxyCommand nc -x .+:\\d+( -X \\d+)? %h %p$")
	proxySSHRegex := regexp.MustCompile("^#(ProxyCommand ssh .+)")

	statements := []string{}
	for _, statement := range original {
		if !proxyRegex.MatchString(statement) {
			if proxySSHRegex.MatchString(statement) {
				statements = append(statements, proxySSHRegex.FindStringSubmatch(statement)[1])
			} else {
				statements = append(statements, statement)
			}
		}
	}
	return statements
}

func removeSSHConfig(config []string) []string {
	sshFile := parseSSHConfig(config)

	hosts := []sshHost{}
	for _, host := range sshFile.Hosts {
		host.Statements = removeFromStatements(host.Statements)
		hosts = append(hosts, host)
	}
	sshFile.Hosts = hosts

	if len(hosts) == 0 {
		sshFile.GlobalStatements = removeFromStatements(sshFile.GlobalStatements)
	}

	return parseSSHFile(sshFile)
}

func addToStatements(original []string, SOCKSProxyHost string, SOCKSProxyPort string) []string {

	proxyRegex := regexp.MustCompile("^ProxyCommand")
	proxySSHRegex := regexp.MustCompile("^ProxyCommand ssh .+")
	proxyOtherRegex := regexp.MustCompile("^ProxyCommand (ssh|nc)")

	statements := []string{}

	for _, statement := range original {

		if proxyRegex.MatchString(statement) && !proxyOtherRegex.MatchString(statement) {
			// don't change anything in the list of statements if this is an 'other' ProxyCommand.
			return original
		}

		if proxySSHRegex.MatchString(statement) {
			statements = append(statements, fmt.Sprintf("#%s", proxySSHRegex.FindStringSubmatch(statement)[0]))
		} else {
			statements = append(statements, statement)
		}
	}
	statements = append(statements, fmt.Sprintf("ProxyCommand nc -x %s:%s %%h %%p", SOCKSProxyHost, SOCKSProxyPort))
	return statements
}

func addSSHConfig(config []string, SOCKSProxyHost string, SOCKSProxyPort string) []string {
	sshFile := parseSSHConfig(config)

	hosts := []sshHost{}
	for _, host := range sshFile.Hosts {
		host.Statements = addToStatements(host.Statements, SOCKSProxyHost, SOCKSProxyPort)
		hosts = append(hosts, host)
	}
	sshFile.Hosts = hosts

	if len(hosts) == 0 {
		sshFile.GlobalStatements = addToStatements(sshFile.GlobalStatements, SOCKSProxyHost, SOCKSProxyPort)
	}

	return parseSSHFile(sshFile)
}

func parseSSHFile(file sshFile) []string {
	output := file.GlobalStatements
	output = append(output, "")
	for _, host := range file.Hosts {
		output = append(output, fmt.Sprintf("Host %s", host.Pattern))
		for _, statement := range host.Statements {
			output = append(output, fmt.Sprintf("    %s", statement))
		}
		output = append(output, "")
	}
	return output
}

func parseSSHConfig(config []string) sshFile {
	hostRegex := regexp.MustCompile("^\\s*Host\\s+(.*)$")

	inHost := false
	var currentHost sshHost
	sshFile := sshFile{}

	for _, line := range config {
		if hostRegex.MatchString(line) {
			if inHost {
				sshFile.Hosts = append(sshFile.Hosts, currentHost)
			}
			inHost = true
			currentHost = sshHost{
				Pattern: hostRegex.FindStringSubmatch(line)[1],
			}
		} else if inHost {
			if len(line) == 0 {
				inHost = false
				sshFile.Hosts = append(sshFile.Hosts, currentHost)
			} else {
				currentHost.Statements = append(currentHost.Statements, strings.TrimSpace(line))
			}
		} else {
			if len(line) != 0 {
				sshFile.GlobalStatements = append(sshFile.GlobalStatements, line)
			}
		}

	}
	if inHost {
		sshFile.Hosts = append(sshFile.Hosts, currentHost)
	}

	return sshFile
}
