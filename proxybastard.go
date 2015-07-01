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
	log.SetOutput(bastardLogger{
		toStdout: true,
		toFile:   false,
		filename: "bastard.log",
	})

	var enableProxies bool

	if len(os.Args) != 2 {
		log.Fatalf("Incorrect args supplied: %s\n", os.Args)
	} else {
		onOffParam := os.Args[1]
		onOffRegexp := regexp.MustCompile("^(on|off)$")
		if len(onOffRegexp.FindStringSubmatch(onOffParam)) != 2 {
			log.Fatalf("Incorrect args supplied: %s\n", os.Args)
		}
		enableProxies = onOffParam == "on"
	}

	configBytes, err := ioutil.ReadFile(util.SanitisePath("~/.proxybastard.json"))
	if err != nil {
		log.Fatal(err)
	}
	config := proxy.ParseConfigurationJSON(configBytes)

	if enableProxies {
		proxy.EnableProxies(config)
	} else {
		proxy.DisableProxies(config)
	}

}
