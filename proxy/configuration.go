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
	Shell []string `json:"shell"`
	Maven []string `json:"maven"`
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
