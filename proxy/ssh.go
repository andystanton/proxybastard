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

// AddToSSH adds proxy settings to ssh file targets.
func AddToSSH(config Configuration) {
	if config.Targets.SSH.Enabled {
		RemoveFromSSH(config)

		for _, sshConfig := range config.Targets.SSH.Files {
			sanitisedPath := util.SanitisePath(sshConfig)
			contents := util.LoadFileIntoSlice(sanitisedPath)
			util.WriteSliceToFile(sanitisedPath, addSSHConfig(contents, config.SocksProxyHost, config.SocksProxyPort))
		}
	}
}

// RemoveFromSSH removes proxy settings from ssh file targets.
func RemoveFromSSH(config Configuration) {
	if config.Targets.SSH.Enabled {
		for _, sshConfig := range config.Targets.SSH.Files {
			sanitisedPath := util.SanitisePath(sshConfig)
			contents := util.LoadFileIntoSlice(sanitisedPath)
			util.WriteSliceToFile(sanitisedPath, removeSSHConfig(contents))
		}
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

func addToStatements(original []string, socksProxyHost string, socksProxyPort string) []string {

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
	statements = append(statements, fmt.Sprintf("ProxyCommand nc -x %s:%s %%h %%p", socksProxyHost, socksProxyPort))
	return statements
}

func addSSHConfig(config []string, socksProxyHost string, socksProxyPort string) []string {
	sshFile := parseSSHConfig(config)

	hosts := []sshHost{}
	for _, host := range sshFile.Hosts {
		host.Statements = addToStatements(host.Statements, socksProxyHost, socksProxyPort)
		hosts = append(hosts, host)
	}
	sshFile.Hosts = hosts

	if len(hosts) == 0 {
		sshFile.GlobalStatements = addToStatements(sshFile.GlobalStatements, socksProxyHost, socksProxyPort)
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
