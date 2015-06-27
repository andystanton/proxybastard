package proxy

import (
	"fmt"
	"regexp"
	"strings"
)

type sshHost struct {
	Pattern    string
	Statements []string
}

type sshFile struct {
	GlobalStatements []string
	Hosts            []sshHost
}

func removeSSHConfig(config []string) []string {
	sshFile := parseSSHConfig(config)

	proxyRegex := regexp.MustCompile("^ProxyCommand nc -x .+:\\d+$")

	hosts := []sshHost{}
	for _, host := range sshFile.Hosts {
		statements := []string{}
		for _, statement := range host.Statements {
			if !proxyRegex.MatchString(statement) {
				statements = append(statements, statement)
			}
		}
		host.Statements = statements
		hosts = append(hosts, host)
	}
	sshFile.Hosts = hosts

	if len(hosts) == 0 {
		statements := []string{}
		for _, statement := range sshFile.GlobalStatements {
			if !proxyRegex.MatchString(statement) {
				statements = append(statements, statement)
			}
		}
		sshFile.GlobalStatements = statements
	}

	return parseSSHFile(sshFile)
}

func addSSHConfig(config []string, socksProxyHost string, socksProxyPort string) []string {
	sshFile := parseSSHConfig(config)

	hosts := []sshHost{}
	for _, host := range sshFile.Hosts {
		host.Statements = append(host.Statements, fmt.Sprintf("ProxyCommand nc -x %s:%s", socksProxyHost, socksProxyPort))
		hosts = append(hosts, host)
	}
	sshFile.Hosts = hosts

	if len(hosts) == 0 {
		sshFile.GlobalStatements = append(sshFile.GlobalStatements, fmt.Sprintf("ProxyCommand nc -x %s:%s", socksProxyHost, socksProxyPort))
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
