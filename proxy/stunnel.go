package proxy

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/andystanton/proxybastard/util"
)

func (stunnelConfiguration StunnelConfiguration) validate() error {
	return nil
}

func (stunnelConfiguration StunnelConfiguration) isEnabled() bool {
	return stunnelConfiguration.Enabled
}

func (stunnelConfiguration StunnelConfiguration) suggestConfiguration() interface{} {
	return nil
}

func (stunnelConfiguration StunnelConfiguration) addSOCKSProxySettings(socksProxyHost string, socksProxyPort string) {
	for _, file := range stunnelConfiguration.Files {
		sanitisedPath := util.SanitisePath(file)
		lines, _ := util.LoadFileIntoSlice(sanitisedPath)
		util.WriteSliceToFile(sanitisedPath, removeStunnelProxies(lines))
	}
	for _, file := range stunnelConfiguration.Files {
		sanitisedPath := util.SanitisePath(file)
		lines, _ := util.LoadFileIntoSlice(sanitisedPath)
		util.WriteSliceToFile(sanitisedPath, addStunnelProxies(lines, socksProxyHost, socksProxyPort))
	}
	if stunnelConfiguration.KillProcess {
		restartStunnel()
	}
}

func (stunnelConfiguration StunnelConfiguration) removeSOCKSProxySettings() {
	for _, file := range stunnelConfiguration.Files {
		sanitisedPath := util.SanitisePath(file)
		lines, _ := util.LoadFileIntoSlice(sanitisedPath)
		util.WriteSliceToFile(sanitisedPath, removeStunnelProxies(lines))
	}
	if stunnelConfiguration.KillProcess {
		restartStunnel()
	}
}

func restartStunnel() {
	util.ShellOut("pkill", []string{"-15", "stunnel"})
}

func addStunnelProxies(contents []string, SOCKSProxyHost string, SOCKSProxyPort string) []string {
	output := []string{}

	socksRegex := regexp.MustCompile("^(execargs\\s*=.*)\\s+(.+)\\s+(\\d+)$")

	for _, line := range contents {
		if socksRegex.MatchString(line) {
			matches := socksRegex.FindStringSubmatch(line)
			output = append(output, fmt.Sprintf("%s -S %s:%s %s %s", matches[1], SOCKSProxyHost, SOCKSProxyPort, matches[2], matches[3]))
		} else {
			output = append(output, line)
		}
	}
	return output
}

func removeStunnelProxies(contents []string) []string {
	output := []string{}

	socksRegex := regexp.MustCompile("(execargs\\s*=.*)-S [\\w.:-]+(.*)")

	for _, line := range contents {
		if socksRegex.MatchString(line) {
			match := socksRegex.FindStringSubmatch(line)
			output = append(output, strings.TrimSpace(fmt.Sprintf("%s %s", strings.TrimSpace(match[1]), strings.TrimSpace(match[2]))))
		} else {
			output = append(output, line)
		}
	}
	return output
}
