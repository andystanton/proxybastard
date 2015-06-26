package proxy

import (
	"fmt"
	"regexp"
	"strings"
)

// SSHHost is an ssh host.
type SSHHost struct {
	Pattern    string
	Statements []string
}

// SSHFile is an ssh file.
type SSHFile struct {
	GlobalStatements []string
	Hosts            []SSHHost
}

func AddSSHConfig(config []string) []string {
	return config
}

// ParseSSHConfig parses a string slice into an SSHFile.
func ParseSSHFile(sshFile SSHFile) []string {
	output := sshFile.GlobalStatements
	output = append(output, "")
	for _, host := range sshFile.Hosts {
		output = append(output, fmt.Sprintf("Host %s", host.Pattern))
		for _, statement := range host.Statements {
			output = append(output, fmt.Sprintf("    %s", statement))
		}
		output = append(output, "")
	}
	return output
}

// ParseSSHConfig parses a string slice into an SSHFile.
func ParseSSHConfig(config []string) SSHFile {
	hostRegex := regexp.MustCompile("^\\s*Host\\s+(.*)$")

	inHost := false
	var currentHost SSHHost
	sshFile := SSHFile{}

	for _, line := range config {
		if hostRegex.MatchString(line) {
			inHost = true
			currentHost = SSHHost{
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
