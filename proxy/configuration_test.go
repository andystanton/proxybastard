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
        "shell": {
			"files": [
				"~/.zshrc",
				"~/.bashrc"
			]
		},
        "maven": {
			"files": [
				"~/.m2/settings.xml"
			]
		}
    }
}
            `,
			Configuration{
				ProxyHost:     "http://www.proxy-bastard.com",
				ProxyPort:     "80",
				NonProxyHosts: []string{"localhost", "127.0.0.1"},
				Targets: Targets{
					Shell: Shell{
						Files: []string{"~/.zshrc", "~/.bashrc"},
					},
					Maven: Maven{
						Files: []string{"~/.m2/settings.xml"},
					},
				},
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
