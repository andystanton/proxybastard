package proxy

import (
	"fmt"
	"regexp"
	"strings"
)

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
