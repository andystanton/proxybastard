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

	var enableProxies bool
	var scan bool
	var backup bool

	if len(os.Args) != 2 {
		log.Fatalf("Incorrect args supplied: %s\n", os.Args)
	} else {
		onOffParam := os.Args[1]
		onOffRegexp := regexp.MustCompile("^(on|off|setup|backup|restore)$")
		if len(onOffRegexp.FindStringSubmatch(onOffParam)) != 2 {
			log.Fatalf("Incorrect args supplied: %s\n", os.Args)
		}
		enableProxies = onOffParam == "on"
		scan = onOffParam == "setup"
		backup = onOffParam == "backup"
	}

	if scan {
		util.SafeWriteSliceToFile("/Users/stanta01/.go/src/github.com/andystanton/proxybastard/blah", []string{})
		proxy.Scan()
	} else {
		configBytes, err := ioutil.ReadFile(util.SanitisePath("~/.proxybastard.json"))
		if err != nil {
			log.Fatal(err)
		}
		config := proxy.ParseConfigurationJSON(configBytes)

		if backup {
			proxy.DirtyBackup(config)
		} else {
			if enableProxies {
				proxy.ToggleProxies(config, proxy.Enable)
			} else {
				proxy.ToggleProxies(config, proxy.Disable)
			}
		}
	}
}
