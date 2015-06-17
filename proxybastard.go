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

	proxyHost := "http://www-cache.reith.bbc.co.uk"
	proxyPort := 80

	config := proxy.Configuration{
		ProxyHost: proxyHost,
		ProxyPort: &proxyPort,
		NonProxyHosts: []string{
			"localhost",
			"127.0.0.1",
			"127.0.0.0/8",
			"::1",
			"192.168.59.103",
			"sandbox.dev.bbc.co.uk",
		},
		ShellFiles: []string{"~/.zshrc"},
		MavenFiles: []string{"~/.m2/settings.xml"},
	}

	proxy.Bastardise(config, enableProxies)

}
