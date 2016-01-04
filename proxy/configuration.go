package proxy

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/andystanton/proxybastard/util"
)

// WithProxy is a thing that can have proxy settings added or removed.
type WithProxy interface {
	addProxySettings(string, string, []string)
	removeProxySettings()
}

// WithSOCKSProxy is a thing that can have SOCKS proxy settings added or removed.
type WithSOCKSProxy interface {
	addSOCKSProxySettings(string, string)
	removeSOCKSProxySettings()
}

// WithConfig is a thing with config.
type WithConfig interface {
	validate() error
	isEnabled() bool
	suggestConfiguration() *Configuration
}

// Configuration represents the Proxybastard configuration.
type Configuration struct {
	Version        string                `json:"version"`
	ProxyHost      string                `json:"proxyHost,omitempty"`
	ProxyPort      string                `json:"proxyPort,omitempty"`
	SOCKSProxyHost string                `json:"socksProxyHost,omitempty"`
	SOCKSProxyPort string                `json:"socksProxyPort,omitempty"`
	NonProxyHosts  []string              `json:"nonProxyHosts,omitempty"`
	Targets        *TargetsConfiguration `json:"targets,omitempty"`
}

// TargetsConfiguration struct.
type TargetsConfiguration struct {
	APM           *APMConfiguration           `json:"apm,omitempty"`
	Boot2Docker   *Boot2DockerConfiguration   `json:"boot2docker,omitempty"`
	Bower         *BowerConfiguration         `json:"bower,omitempty"`
	DockerMachine *DockerMachineConfiguration `json:"docker-machine,omitempty"`
	Git           *GitConfiguration           `json:"git,omitempty"`
	Maven         *MavenConfiguration         `json:"maven,omitempty"`
	NPM           *NPMConfiguration           `json:"npm,omitempty"`
	Shell         *ShellConfiguration         `json:"shell,omitempty"`
	SSH           *SSHConfiguration           `json:"ssh,omitempty"`
	Stunnel       *StunnelConfiguration       `json:"stunnel,omitempty"`
	Subversion    *SubversionConfiguration    `json:"subversion,omitempty"`
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

// BowerConfiguration struct.
type BowerConfiguration struct {
	Enabled bool     `json:"enabled,omitempty"`
	Files   []string `json:"files,omitempty"`
}

func parseConfigurationJSON(jsoncontent []byte) Configuration {
	configuration := Configuration{}
	err := json.Unmarshal(jsoncontent, &configuration)
	if err != nil {
		log.Fatal(err)
	}
	configuration.ProxyHost = util.SanitiseHTTPProxyURL(configuration.ProxyHost)
	return configuration
}

// LoadConfigurationFromFile loads a configuration from a file.
func LoadConfigurationFromFile(filename string) (Configuration, error) {
	configBytes, err := ioutil.ReadFile(util.SanitisePath(filename))
	if err != nil {
		return Configuration{}, err
	}
	return parseConfigurationJSON(configBytes), nil
}
