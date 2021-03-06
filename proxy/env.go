package proxy

import (
	"fmt"
	"log"
	"strings"

	"github.com/andystanton/proxybastard/util"
)

// PrintEnv prints a list of the environment settings.
func PrintEnv(config Configuration) {
	envSettings := []string{}
	shellConfiguration := config.Targets.Shell
	if shellConfiguration.Enabled && len(shellConfiguration.Files) > 0 {
		shellFile := util.SanitisePath(shellConfiguration.Files[0])
		if util.FileExists(shellFile) {
			contents, err := util.LoadFileIntoSlice(shellFile)
			if err != nil {
				log.Fatal(err)
			}
			proxyHost, proxyPort, nonProxyHosts, _ := extractProxyFromShellContents(contents)
			if len(proxyHost) > 0 {
				for _, proxyVar := range proxyVars {
					envSettings = append(envSettings, fmt.Sprintf("export %s=%s:%s", proxyVar, proxyHost, proxyPort))
				}
				if len(nonProxyHosts) > 0 {
					for _, noProxyVar := range nonProxyVars {
						envSettings = append(envSettings, fmt.Sprintf("export %s=%s", noProxyVar, strings.Join(nonProxyHosts, ",")))
					}
				}
			} else {
				for _, proxyVar := range proxyVars {
					envSettings = append(envSettings, fmt.Sprintf("unset %s", proxyVar))
				}
				for _, noProxyVar := range nonProxyVars {
					envSettings = append(envSettings, fmt.Sprintf("unset %s", noProxyVar))
				}
			}
		}

	}
	for _, setting := range envSettings {
		fmt.Println(setting)
	}
}
