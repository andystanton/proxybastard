package proxy

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAddSSHConfig(t *testing.T) {
	socksProxyHost := "socks-proxy.proxybastard.com"
	socksProxyPort := "1085"

	cases := []struct {
		config   []string
		expected []string
	}{
		{
			[]string{
				"VisualHostKey yes",
				"",
				"Host foo",
				"    Hostname 123.132.213.231",
				"    Port 1234",
				"    User foo",
				"    IdentityFile ~/.ssh/foo",
				"",
				"Host bar",
				"    IdentityFile ~/.ssh/bar",
				"",
			},
			[]string{
				"VisualHostKey yes",
				"",
				"Host foo",
				"    Hostname 123.132.213.231",
				"    Port 1234",
				"    User foo",
				"    IdentityFile ~/.ssh/foo",
				fmt.Sprintf("    ProxyCommand nc -x %s:%s", socksProxyHost, socksProxyPort),
				"",
				"Host bar",
				"    IdentityFile ~/.ssh/bar",
				fmt.Sprintf("    ProxyCommand nc -x %s:%s", socksProxyHost, socksProxyPort),
				"",
			},
		},
	}

	for _, c := range cases {
		actual := addSSHConfig(c.config, socksProxyHost, socksProxyPort)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`
Call:
addSSHConfig() != {{expected}}

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
				c.config,
				c.expected,
				actual)
		}
	}
}

func TestParseSSHFile(t *testing.T) {
	cases := []struct {
		file   sshFile
		config []string
	}{
		{
			sshFile{
				GlobalStatements: []string{
					"VisualHostKey yes",
				},
				Hosts: []sshHost{
					{
						Pattern: "foo",
						Statements: []string{
							"Hostname 123.132.213.231",
							"Port 1234",
							"User foo",
							"IdentityFile ~/.ssh/foo",
						},
					},
					{
						Pattern: "bar",
						Statements: []string{
							"IdentityFile ~/.ssh/bar",
						},
					},
				},
			},
			[]string{
				"VisualHostKey yes",
				"",
				"Host foo",
				"    Hostname 123.132.213.231",
				"    Port 1234",
				"    User foo",
				"    IdentityFile ~/.ssh/foo",
				"",
				"Host bar",
				"    IdentityFile ~/.ssh/bar",
				"",
			},
		},
	}

	for _, c := range cases {
		actual := parseSSHFile(c.file)
		if !reflect.DeepEqual(actual, c.config) {
			t.Errorf(
				`
Call:
parseSSHFile() != {{expected}}

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

				c.file,
				c.config,
				actual)
		}
	}
}

func TestParseSSHConfig(t *testing.T) {
	cases := []struct {
		config []string
		file   sshFile
	}{
		{
			[]string{
				"VisualHostKey yes",
				"",
				"Host foo",
				"    Hostname 123.132.213.231",
				"    Port 1234",
				"    User foo",
				"    IdentityFile ~/.ssh/foo",
				"",
				"Host bar",
				"    IdentityFile ~/.ssh/bar",
			},
			sshFile{
				GlobalStatements: []string{
					"VisualHostKey yes",
				},
				Hosts: []sshHost{
					{
						Pattern: "foo",
						Statements: []string{
							"Hostname 123.132.213.231",
							"Port 1234",
							"User foo",
							"IdentityFile ~/.ssh/foo",
						},
					},
					{
						Pattern: "bar",
						Statements: []string{
							"IdentityFile ~/.ssh/bar",
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		actual := parseSSHConfig(c.config)
		if !reflect.DeepEqual(actual, c.file) {
			t.Errorf(
				`
Call:
parseSSHConfig() != {{expected}}

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
				c.config,
				c.file,
				actual)
		}
	}
}
