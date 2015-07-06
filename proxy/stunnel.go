package proxy

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/andystanton/proxybastard/util"
)

func addToStunnel(config Configuration) {
	if config.Targets.Stunnel.Enabled {
		for _, file := range config.Targets.Stunnel.Files {
			sanitisedPath := util.SanitisePath(file)
			lines := util.LoadFileIntoSlice(sanitisedPath)
			util.WriteSliceToFile(sanitisedPath, removeStunnelProxies(lines))
		}
		for _, file := range config.Targets.Stunnel.Files {
			sanitisedPath := util.SanitisePath(file)
			lines := util.LoadFileIntoSlice(sanitisedPath)
			util.WriteSliceToFile(sanitisedPath, addStunnelProxies(lines, config.SocksProxyHost, config.SocksProxyPort))
		}
		if config.Targets.Stunnel.KillProcess {
			restartStunnel()
		}
	}
}

func removeFromStunnel(config Configuration) {
	if config.Targets.Stunnel.Enabled {
		for _, file := range config.Targets.Stunnel.Files {
			sanitisedPath := util.SanitisePath(file)
			lines := util.LoadFileIntoSlice(sanitisedPath)
			util.WriteSliceToFile(sanitisedPath, removeStunnelProxies(lines))
		}
		if config.Targets.Stunnel.KillProcess {
			restartStunnel()
		}
	}
}

func restartStunnel() {
	util.ShellOut("pkill", []string{"-15", "stunnel"})
}

func addStunnelProxies(contents []string, socksProxyHost string, socksProxyPort string) []string {
	output := []string{}

	socksRegex := regexp.MustCompile("execargs\\s*=.*")

	for _, line := range contents {
		if socksRegex.MatchString(line) {
			output = append(output, fmt.Sprintf("%s -S %s:%s", line, socksProxyHost, socksProxyPort))
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
