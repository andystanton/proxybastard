package proxy

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRemoveStunnelProxies(t *testing.T) {
	cases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			"removeStunnelProxies",
			[]string{
				"; stunnel pid",
				"pid = /path/to/stunnel.pid ",
				"",
				"; Certs",
				"cert = /path/to/client.p12",
				"CAfile = /path/to/ca.pem",
				"",
				"; Protocol version (all, SSLv2, SSLv3, TLSv1)",
				"sslVersion = TLSv1",
				"verify = 2",
				"",
				"[group]",
				"accept  = 127.0.0.1:1234",
				"exec = /usr/local/bin/connect",
				"execargs = /usr/local/bin/connect -w 5 -S socks.proxy-bastard.com:1085 internal.server 4321",
				"client = yes",
			},
			[]string{
				"; stunnel pid",
				"pid = /path/to/stunnel.pid ",
				"",
				"; Certs",
				"cert = /path/to/client.p12",
				"CAfile = /path/to/ca.pem",
				"",
				"; Protocol version (all, SSLv2, SSLv3, TLSv1)",
				"sslVersion = TLSv1",
				"verify = 2",
				"",
				"[group]",
				"accept  = 127.0.0.1:1234",
				"exec = /usr/local/bin/connect",
				"execargs = /usr/local/bin/connect -w 5 internal.server 4321",
				"client = yes",
			},
		},
	}
	for _, c := range cases {
		actual := removeStunnelProxies(c.input)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`%s
Call:
removeStunnelProxies({{input}}) != {{expected}}

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
				c.name,
				c.input,
				c.expected,
				actual)
		}
	}
}

func TestAddStunnelProxies(t *testing.T) {
	socksProxyHost := "socks.proxy-bastard.com"
	socksProxyPort := "1085"

	cases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			"removeStunnelProxies",
			[]string{
				"; stunnel pid",
				"pid = /Users/stanta01/.stunnel/.stunnel.pid ",
				"",
				"; Certs",
				"cert = /path/to/client.p12",
				"CAfile = /path/to/ca.pem",
				"",
				"; Protocol version (all, SSLv2, SSLv3, TLSv1)",
				"sslVersion = TLSv1",
				"verify = 2",
				"",
				"[group]",
				"accept  = 127.0.0.1:1234",
				"exec = /usr/local/bin/connect",
				"execargs = /usr/local/bin/connect -w 5 internal.server 4321",
				"client = yes",
			},
			[]string{
				"; stunnel pid",
				"pid = /Users/stanta01/.stunnel/.stunnel.pid ",
				"",
				"; Certs",
				"cert = /path/to/client.p12",
				"CAfile = /path/to/ca.pem",
				"",
				"; Protocol version (all, SSLv2, SSLv3, TLSv1)",
				"sslVersion = TLSv1",
				"verify = 2",
				"",
				"[group]",
				"accept  = 127.0.0.1:1234",
				"exec = /usr/local/bin/connect",
				fmt.Sprintf("execargs = /usr/local/bin/connect -w 5 internal.server 4321 -S %s:%s", socksProxyHost, socksProxyPort),
				"client = yes",
			},
		},
	}
	for _, c := range cases {
		actual := addStunnelProxies(c.input, socksProxyHost, socksProxyPort)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`%s
Call:
addStunnelProxies({{input}}) != {{expected}}

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
				c.name,
				c.input,
				c.expected,
				actual)
		}
	}
}
