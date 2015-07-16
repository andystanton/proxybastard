package main

import (
	"log"
	"os"
	"regexp"

	"github.com/andystanton/proxybastard/proxy"
	"github.com/andystanton/proxybastard/util"
)

func getMode(args []string) string {
	if len(os.Args) != 2 {
		log.Fatalf("Incorrect args supplied: %s\n", args)
	}

	modeRegex := regexp.MustCompile("^(on|off|setup|backup|restore)$")
	modeMatch := modeRegex.FindStringSubmatch(args[1])
	if len(modeMatch) != 2 {
		log.Fatalf("Incorrect args supplied: %s\n", args)
	}

	return modeMatch[0]
}

func main() {
	util.ConfigureLog("bastard.log", true, false)

	if mode := getMode(os.Args); mode == "setup" {
		proxy.Scan()
	} else {
		config := proxy.LoadConfigurationFromFile("~/.proxybastard.json")
		switch mode {
		case "on":
			proxy.ToggleProxies(config, proxy.Enable)
		case "off":
			proxy.ToggleProxies(config, proxy.Disable)
		case "backup":
			proxy.DirtyBackupOperation(config, proxy.Backup)
		case "restore":
			proxy.DirtyBackupOperation(config, proxy.Restore)
		}
	}
}
