package proxy

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRemoveSSHConfig(t *testing.T) {
	cases := []struct {
		config   []string
		expected []string
	}{
		{
			[]string{
				"VisualHostKey yes",
				"ProxyCommand nc -x socks-proxy.proxybastard.com:1085",
				"",
			},
			[]string{
				"VisualHostKey yes",
				"",
			},
		},
		{
			[]string{
				"VisualHostKey yes",
				"",
				"Host foo",
				"    Hostname 123.132.213.231",
				"    Port 1234",
				"    User foo",
				"    IdentityFile ~/.ssh/foo",
				"    ProxyCommand nc -x socks-proxy.proxybastard.com:1085",
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
				"",
				"Host bar",
				"    IdentityFile ~/.ssh/bar",
				"",
			},
		},
		{
			[]string{
				"VisualHostKey yes",
				"",
				"Host foo",
				"    Hostname 123.132.213.231",
				"    Port 1234",
				"    User foo",
				"    IdentityFile ~/.ssh/foo",
				"    ProxyCommand nc -x socks-proxy.proxybastard.com:1085",
				"",
				"Host bar",
				"    IdentityFile ~/.ssh/bar",
				"    ProxyCommand nc -x socks-proxy.proxybastard.com:1085",
				"",
				"Host access.*.proxybastard.com",
				"    IdentityFile ~/.ssh/foobar",
				"    ProxyCommand nc -x socks-proxy.proxybastard.com:1085",
				"    #ProxyCommand ssh -q -p 2222 user@bastion nc %h %p",
				"",
				"Host *,*",
				"    ProxyCommand >&1; h=\"%h\"; exec ssh -q bastion@access.${h##*,}.proxybastard.com nc ${h%%,*} %p",
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
				"",
				"Host bar",
				"    IdentityFile ~/.ssh/bar",
				"",
				"Host access.*.proxybastard.com",
				"    IdentityFile ~/.ssh/foobar",
				"    ProxyCommand ssh -q -p 2222 user@bastion nc %h %p",
				"",
				"Host *,*",
				"    ProxyCommand >&1; h=\"%h\"; exec ssh -q bastion@access.${h##*,}.proxybastard.com nc ${h%%,*} %p",
				"",
			},
		},
	}

	for _, c := range cases {
		actual := removeSSHConfig(c.config)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`
Call:
removeSSHConfig() != {{expected}}

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
			},
			[]string{
				"VisualHostKey yes",
				fmt.Sprintf("ProxyCommand nc -x %s:%s", socksProxyHost, socksProxyPort),
				"",
			},
		},
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
				"Host access.*.proxybastard.com",
				"    IdentityFile ~/.ssh/foobar",
				"    ProxyCommand ssh -q -p 2222 user@bastion nc %h %p",
				"",
				"Host *,*",
				"    ProxyCommand >&1; h=\"%h\"; exec ssh -q bastion@access.${h##*,}.proxybastard.com nc ${h%%,*} %p",
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
				"Host access.*.proxybastard.com",
				"    IdentityFile ~/.ssh/foobar",
				"    #ProxyCommand ssh -q -p 2222 user@bastion nc %h %p",
				fmt.Sprintf("    ProxyCommand nc -x %s:%s", socksProxyHost, socksProxyPort),
				"",
				"Host *,*",
				"    ProxyCommand >&1; h=\"%h\"; exec ssh -q bastion@access.${h##*,}.proxybastard.com nc ${h%%,*} %p",
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
