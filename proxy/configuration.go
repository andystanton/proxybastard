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
	Shell       ShellConfiguration       `json:"shell"`
	Maven       MavenConfiguration       `json:"maven"`
	SSH         SSHConfiguration         `json:"ssh"`
	Git         GitConfiguration         `json:"git"`
	NPM         NPMConfiguration         `json:"npm"`
	APM         APMConfiguration         `json:"apm"`
	Subversion  SubversionConfiguration  `json:"subversion"`
	Boot2Docker Boot2DockerConfiguration `json:"boot2docker"`
	Stunnel     StunnelConfiguration     `json:"stunnel"`
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

// SubversionConfiguration struct.
type SubversionConfiguration struct {
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

// APMConfiguration struct.
type APMConfiguration struct {
	Enabled bool `json:"enabled"`
}

// Boot2DockerConfiguration struct.
type Boot2DockerConfiguration struct {
	Enabled bool   `json:"enabled"`
	SSHHost string `json:"ssh_host"`
	SSHPort string `json:"ssh_port"`
	SSHKey  string `json:"ssh_key"`
}

// StunnelConfiguration struct.
type StunnelConfiguration struct {
	Enabled bool     `json:"enabled"`
	Files   []string `json:"files"`
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
