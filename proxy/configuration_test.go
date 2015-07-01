package proxy

import (
	"reflect"
	"testing"
)

func TestParseConfigurationJSON(t *testing.T) {
	cases := []struct {
		json     string
		expected Configuration
	}{
		{
			`{}`,
			Configuration{},
		},
		{
			`
{
    "proxyHost": "http://www.proxy-bastard.com",
    "proxyPort": "80",
    "nonProxyHosts": [
		"localhost",
		"127.0.0.1"
	],
    "targets": {
		"git": {
			"enabled": true
		},
		"npm": {
			"enabled": true
		},
		"apm": {
			"enabled": true
		},
        "shell": {
			"enabled": true,
			"javaOpts": true,
			"files": [
				"~/.zshrc",
				"~/.bashrc"
			]
		},
        "maven": {
			"enabled": true,
			"files": [
				"~/.m2/settings.xml"
			]
		},
		"ssh": {
			"enabled": true,
			"files": [
				"~/.ssh/config"
			]
		},
		"subversion": {
			"enabled": true,
			"files": [
				"~/.subversion/servers"
			]
		},
		"boot2docker": {
			"enabled": true,
			"ssh_host": "192.168.59.103",
			"ssh_port": "22",
			"ssh_key": "~/.ssh/id_boot2docker"
		}
    }
}
            `,
			Configuration{
				ProxyHost:     "http://www.proxy-bastard.com",
				ProxyPort:     "80",
				NonProxyHosts: []string{"localhost", "127.0.0.1"},
				Targets: TargetsConfiguration{
					Git: GitConfiguration{
						Enabled: true,
					},
					NPM: NPMConfiguration{
						Enabled: true,
					},
					APM: APMConfiguration{
						Enabled: true,
					},
					Shell: ShellConfiguration{
						Enabled:  true,
						JavaOpts: true,
						Files:    []string{"~/.zshrc", "~/.bashrc"},
					},
					Maven: MavenConfiguration{
						Enabled: true,
						Files:   []string{"~/.m2/settings.xml"},
					},
					SSH: SSHConfiguration{
						Enabled: true,
						Files:   []string{"~/.ssh/config"},
					},
					Subversion: SubversionConfiguration{
						Enabled: true,
						Files:   []string{"~/.subversion/servers"},
					},
					Boot2Docker: Boot2DockerConfiguration{
						Enabled: true,
						SSHHost: "192.168.59.103",
						SSHPort: "22",
						SSHKey:  "~/.ssh/id_boot2docker",
					},
				},
			},
		},
		{
			`
{
	"proxyHost": "http://www.proxy-bastard.com",
	"proxyPort": "80",
	"socksProxyHost": "socks.proxy-bastard.com",
	"socksProxyPort": "1085",
	"nonProxyHosts": [
		"localhost",
		"127.0.0.1"
	]
}
			`,
			Configuration{
				ProxyHost:      "http://www.proxy-bastard.com",
				ProxyPort:      "80",
				SocksProxyHost: "socks.proxy-bastard.com",
				SocksProxyPort: "1085",
				NonProxyHosts:  []string{"localhost", "127.0.0.1"},
			},
		},
	}
	for _, c := range cases {
		actual := ParseConfigurationJSON([]byte(c.json))
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`
Call:
ParseConfigurationJson({{input}}) != {{expected}}

Input:
===============
%s
===============

Expected:
===============
%s
===============

Actual:
===============
%s
===============`,
				c.json,
				c.expected,
				actual)
		}
	}
}
