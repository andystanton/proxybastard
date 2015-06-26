package proxy

import (
	"encoding/json"
	"log"
)

// Configuration represents the Proxybastard configuration.
type Configuration struct {
	ProxyHost      string               `json:"proxyHost"`
	ProxyPort      string               `json:"proxyPort"`
	SocksProxyHost string               `json:"socksProxyHost"`
	SocksProxyPort string               `json:"socksProxyPort"`
	NonProxyHosts  []string             `json:"nonProxyHosts"`
	Targets        TargetsConfiguration `json:"targets"`
}

// TargetsConfiguration struct.
type TargetsConfiguration struct {
	Shell ShellConfiguration `json:"shell"`
	Maven MavenConfiguration `json:"maven"`
}

// ShellConfiguration struct.
type ShellConfiguration struct {
	JavaOpts bool     `json:"javaOpts"`
	Files    []string `json:"files"`
}

// MavenConfiguration struct.
type MavenConfiguration struct {
	Files []string `json:"files"`
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
