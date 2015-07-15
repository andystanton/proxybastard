package proxy

import (
	"encoding/json"
	"log"
)

// WithProxy is a thing that can have proxy settings added or removed.
type WithProxy interface {
	addProxySettings(string, string, []string)
	removeProxySettings()
	validate() error
	isEnabled() bool
}

// WithSOCKSProxy is a thing that can have SOCKS proxy settings added or removed.
type WithSOCKSProxy interface {
	addSOCKSProxySettings(string, string)
	removeSOCKSProxySettings()
	validate() error
	isEnabled() bool
}

// Configuration represents the Proxybastard configuration.
type Configuration struct {
	ProxyHost      string               `json:"proxyHost,omitempty"`
	ProxyPort      string               `json:"proxyPort,omitempty"`
	SOCKSProxyHost string               `json:"socksProxyHost,omitempty"`
	SOCKSProxyPort string               `json:"socksProxyPort,omitempty"`
	NonProxyHosts  []string             `json:"nonProxyHosts,omitempty"`
	Targets        TargetsConfiguration `json:"targets,omitempty"`
}

// TargetsConfiguration struct.
type TargetsConfiguration struct {
	Shell         *ShellConfiguration         `json:"shell,omitempty"`
	Maven         *MavenConfiguration         `json:"maven,omitempty"`
	SSH           *SSHConfiguration           `json:"ssh,omitempty"`
	Git           *GitConfiguration           `json:"git,omitempty"`
	NPM           *NPMConfiguration           `json:"npm,omitempty"`
	APM           *APMConfiguration           `json:"apm,omitempty"`
	Subversion    *SubversionConfiguration    `json:"subversion,omitempty"`
	Boot2Docker   *Boot2DockerConfiguration   `json:"boot2docker,omitempty"`
	Stunnel       *StunnelConfiguration       `json:"stunnel,omitempty"`
	DockerMachine *DockerMachineConfiguration `json:"docker-machine,omitempty"`
}

// ShellConfiguration struct.
type ShellConfiguration struct {
	Enabled  bool     `json:"enabled,omitempty"`
	JavaOpts bool     `json:"javaOpts,omitempty"`
	Files    []string `json:"files,omitempty"`
}

// MavenConfiguration struct.
type MavenConfiguration struct {
	Enabled bool     `json:"enabled,omitempty"`
	Files   []string `json:"files,omitempty"`
}

// SSHConfiguration struct.
type SSHConfiguration struct {
	Enabled bool     `json:"enabled,omitempty"`
	Files   []string `json:"files,omitempty"`
}

// SubversionConfiguration struct.
type SubversionConfiguration struct {
	Enabled bool     `json:"enabled,omitempty"`
	Files   []string `json:"files,omitempty"`
}

// GitConfiguration struct.
type GitConfiguration struct {
	Enabled bool `json:"enabled,omitempty"`
}

// NPMConfiguration struct.
type NPMConfiguration struct {
	Enabled bool     `json:"enabled,omitempty"`
	Files   []string `json:"files,omitempty"`
}

// APMConfiguration struct.
type APMConfiguration struct {
	Enabled bool     `json:"enabled,omitempty"`
	Files   []string `json:"files,omitempty"`
}

// Boot2DockerConfiguration struct.
type Boot2DockerConfiguration struct {
	Enabled bool   `json:"enabled,omitempty"`
	SSHHost string `json:"ssh_host,omitempty"`
	SSHPort string `json:"ssh_port,omitempty"`
	SSHKey  string `json:"ssh_key,omitempty"`
}

// StunnelConfiguration struct.
type StunnelConfiguration struct {
	Enabled     bool     `json:"enabled,omitempty"`
	KillProcess bool     `json:"kill_process,omitempty"`
	Files       []string `json:"files,omitempty"`
}

// DockerMachineConfiguration struct.
type DockerMachineConfiguration struct {
	Enabled bool     `json:"enabled,omitempty"`
	Hosts   []string `json:"hosts,omitempty"`
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
