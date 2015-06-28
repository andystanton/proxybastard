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
	SSH   SSHConfiguration   `json:"ssh"`
	Git   GitConfiguration   `json:"git"`
	NPM   NPMConfiguration   `json:"npm"`
}

// ShellConfiguration struct.
type ShellConfiguration struct {
	Enabled  bool     `json:"enabled"`
	JavaOpts bool     `json:"javaOpts"`
	Files    []string `json:"files"`
}

// MavenConfiguration struct.
type MavenConfiguration struct {
	Enabled bool     `json:"enabled"`
	Files   []string `json:"files"`
}

// SSHConfiguration struct.
type SSHConfiguration struct {
	Enabled bool     `json:"enabled"`
	Files   []string `json:"files"`
}

// GitConfiguration struct.
type GitConfiguration struct {
	Enabled bool `json:"enabled"`
}

// NPMConfiguration struct.
type NPMConfiguration struct {
	Enabled bool `json:"enabled"`
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
