package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/andystanton/proxybastard/proxy"
	"github.com/andystanton/proxybastard/util"
)

var version = "1.0.0-alpha"

func getMode(args []string) string {
	if len(os.Args) != 2 {
		printHelp()
		os.Exit(1)
	}

	modeRegex := regexp.MustCompile("^(on|off|setup|backup|restore|env|version)$")
	modeMatch := modeRegex.FindStringSubmatch(args[1])
	if len(modeMatch) != 2 {
		printHelp()
		os.Exit(1)
	}

	return modeMatch[0]
}

func printHelp() {
	fmt.Printf("%s\n", "Usage:")
	fmt.Printf("\t%s\t%s\n", "proxybastard", "<operation>")
	fmt.Printf("\n")
	fmt.Printf("%s\n", "Operations:")
	fmt.Printf("\t%s\t\t%s\n", "on", "Enables proxy settings.")
	fmt.Printf("\t%s\t\t%s\n", "off", "Disables proxy settings.")
	fmt.Printf("\t%s\t\t%s\n", "env", "Outputs export/unset statements based on the shell profile/rc.")
	fmt.Printf("\t%s\t\t%s\n", "setup", "Run the setup helper.")
	fmt.Printf("\t%s\t\t%s\n", "backup", "Backs up existing configuration.")
	fmt.Printf("\t%s\t\t%s\n", "restore", "Restores previous configuration.")
	fmt.Printf("\t%s\t\t%s\n", "version", "Displays the application version.")
}

func printVersion() {
	fmt.Printf("%s\n", version)
}

func main() {
	util.ConfigureLog("bastard.log", true, false)
	config, err := proxy.LoadConfigurationFromFile("~/.proxybastard/config.json")

	if mode := getMode(os.Args); mode == "version" {
		printVersion()
	} else if mode == "setup" || err != nil {
		if err != nil {
			fmt.Println("It looks like you don't have a proxybastard configuration yet. Let's get started!")
		}
		proxy.Setup()

	} else {
		switch mode {
		case "on":
			proxy.ToggleProxies(config, proxy.Enable)
		case "off":
			proxy.ToggleProxies(config, proxy.Disable)
		case "env":
			proxy.PrintEnv(config)
		case "backup":
			proxy.DirtyBackupOperation(config, proxy.Backup)
		case "restore":
			proxy.DirtyBackupOperation(config, proxy.Restore)
		}
	}
}
