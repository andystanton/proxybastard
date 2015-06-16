package main

import (
	"log"
	"os"
	"regexp"

	"github.com/andystanton/proxybastard/proxy"
)

func main() {
	log.SetOutput(bastardLogger{
		toStdout: true,
		toFile:   true,
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

	proxyHost := "http://www.proxy-bastard.com"
	proxyPort := 80

	config := proxy.Configuration{
		ProxyHost:  proxyHost,
		ProxyPort:  &proxyPort,
		ShellFiles: []string{"~/.zshrc"},
		MavenFiles: []string{"sagdsg"},
	}

	proxy.Bastardise(config, enableProxies)

}
