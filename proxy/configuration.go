package proxy

import (
	"encoding/json"
	"log"
)

// Configuration represents the Proxybastard configuration.
type Configuration struct {
	ProxyHost     string   `json:"proxyHost"`
	ProxyPort     string   `json:"proxyPort"`
	NonProxyHosts []string `json:"nonProxyHosts"`
	Targets       Targets  `json:"targets"`
}

type Targets struct {
	Shell Shell `json:"shell"`
	Maven Maven `json:"maven"`
}

type Shell struct {
	JavaOpts bool     `json:javaOpts`
	Files    []string `json:files`
}

type Maven struct {
	Files []string `json:files`
}

// ParseConfigurationJSON parses configuration json.
func ParseConfigurationJSON(jsoncontent []byte) Configuration {
	res := Configuration{}
	err := json.Unmarshal(jsoncontent, &res)
	if err != nil {
		log.Fatal(err)
	}
	return res
}
