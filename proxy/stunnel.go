package proxy

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/andystanton/proxybastard/util"
)

func (stunnelConfiguration StunnelConfiguration) addSocksProxySettings(socksProxyHost string, socksProxyPort string) {
	if stunnelConfiguration.Enabled {
		for _, file := range stunnelConfiguration.Files {
			sanitisedPath := util.SanitisePath(file)
			lines := util.LoadFileIntoSlice(sanitisedPath)
			util.WriteSliceToFile(sanitisedPath, removeStunnelProxies(lines))
		}
		for _, file := range stunnelConfiguration.Files {
			sanitisedPath := util.SanitisePath(file)
			lines := util.LoadFileIntoSlice(sanitisedPath)
			util.WriteSliceToFile(sanitisedPath, addStunnelProxies(lines, socksProxyHost, socksProxyPort))
		}
		if stunnelConfiguration.KillProcess {
			restartStunnel()
		}
	}
}

func (stunnelConfiguration StunnelConfiguration) removeSocksProxySettings() {
	if stunnelConfiguration.Enabled {
		for _, file := range stunnelConfiguration.Files {
			sanitisedPath := util.SanitisePath(file)
			lines := util.LoadFileIntoSlice(sanitisedPath)
			util.WriteSliceToFile(sanitisedPath, removeStunnelProxies(lines))
		}
		if stunnelConfiguration.KillProcess {
			restartStunnel()
		}
	}
}

func restartStunnel() {
	util.ShellOut("pkill", []string{"-15", "stunnel"})
}

func addStunnelProxies(contents []string, SOCKSProxyHost string, SOCKSProxyPort string) []string {
	output := []string{}

	socksRegex := regexp.MustCompile("execargs\\s*=.*")

	for _, line := range contents {
		if socksRegex.MatchString(line) {
			output = append(output, fmt.Sprintf("%s -S %s:%s", line, SOCKSProxyHost, SOCKSProxyPort))
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
