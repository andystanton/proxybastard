package main

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/andystanton/proxybastard/proxy"
	"github.com/andystanton/proxybastard/util"
)

func main() {
	log.SetOutput(util.TeeLogger{
		ToStdout: true,
		ToFile:   false,
		Filename: "bastard.log",
	})

	log.SetFlags(log.Flags() ^ log.Ldate)
	log.SetFlags(log.Flags() ^ log.Ltime)

	var mode string

	if len(os.Args) != 2 {
		log.Fatalf("Incorrect args supplied: %s\n", os.Args)
	} else {
		modeParam := os.Args[1]
		modeRegex := regexp.MustCompile("^(on|off|setup|backup|restore)$")
		if len(modeRegex.FindStringSubmatch(modeParam)) != 2 {
			log.Fatalf("Incorrect args supplied: %s\n", os.Args)
		}
		mode = modeRegex.FindStringSubmatch(modeParam)[0]
	}

	switch mode {
	case "on":
		proxy.ToggleProxies(getConfig(), proxy.Enable)
	case "off":
		proxy.ToggleProxies(getConfig(), proxy.Disable)
	case "setup":
		proxy.Scan()
	case "backup":
		proxy.DirtyBackupOperation(getConfig(), proxy.Backup)
	case "restore":
		proxy.DirtyBackupOperation(getConfig(), proxy.Restore)
	}
}

func getConfig() proxy.Configuration {
	configBytes, err := ioutil.ReadFile(util.SanitisePath("~/.proxybastard.json"))
	if err != nil {
		log.Fatal(err)
	}
	return proxy.ParseConfigurationJSON(configBytes)
}
