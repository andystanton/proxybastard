package proxy

import (
	"reflect"
	"testing"
)

func TestAddSSHConfig(t *testing.T) {
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
			},
		},
	}

	for _, c := range cases {
		actual := AddSSHConfig(c.config)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`
Call:
AddSSHConfig() != {{expected}}

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
		sshFile SSHFile
		config  []string
	}{
		{
			SSHFile{
				GlobalStatements: []string{
					"VisualHostKey yes",
				},
				Hosts: []SSHHost{
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
		actual := ParseSSHFile(c.sshFile)
		if !reflect.DeepEqual(actual, c.config) {
			t.Errorf(
				`
Call:
ParseSSHFile() != {{expected}}

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

				c.sshFile,
				c.config,
				actual)
		}
	}
}

func TestParseSSHConfig(t *testing.T) {
	cases := []struct {
		config  []string
		sshFile SSHFile
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
			SSHFile{
				GlobalStatements: []string{
					"VisualHostKey yes",
				},
				Hosts: []SSHHost{
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
		actual := ParseSSHConfig(c.config)
		if !reflect.DeepEqual(actual, c.sshFile) {
			t.Errorf(
				`
Call:
ParseSSHConfig() != {{expected}}

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
				c.sshFile,
				actual)
		}
	}
}
