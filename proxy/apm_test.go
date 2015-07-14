package proxy

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestAddAPMProxySettings(t *testing.T) {
	proxyHost := "http://www.proxy-bastard.com"
	proxyPort := "80"
	nonProxyHosts := []string{"localhost", "127.0.0.1", "::1"}

	cases := []struct {
		name          string
		proxyHost     string
		proxyPort     string
		nonProxyHosts []string
		input         []string
		expected      []string
	}{
		{
			"addAPMProxySettings",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			[]string{
				"foo=bar",
			},
			[]string{
				"foo=bar",
				fmt.Sprintf("http-proxy=%s:%s", proxyHost, proxyPort),
				fmt.Sprintf("https-proxy=%s:%s", proxyHost, proxyPort),
				"strict-ssl=false",
			},
		},
	}
	for _, c := range cases {
		actual := addAPMProxySettings(c.input, c.proxyHost, c.proxyPort)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`%s
Call:
addAPMProxySettings({{input}}, %s, %s) != {{expected}}

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
				c.proxyHost,
				c.proxyPort,
				strings.Join(c.input, "\n"),
				strings.Join(c.expected, "\n"),
				strings.Join(actual, "\n"))
		}
	}
}

func TestRemoveAPMProxySettings(t *testing.T) {
	proxyHost := "http://www.proxy-bastard.com"
	proxyPort := "80"

	cases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			"removeAPMProxySettings",
			[]string{
				"foo=bar",
				fmt.Sprintf("http-proxy=%s:%s", proxyHost, proxyPort),
				fmt.Sprintf("https-proxy=%s:%s", proxyHost, proxyPort),
				"strict-ssl=false",
				"bar=foo",
			},
			[]string{
				"foo=bar",
				"bar=foo",
			},
		},
	}
	for _, c := range cases {
		actual := removeAPMProxySettings(c.input)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`%s
Call:
removeAPMProxySettings({{input}}) != {{expected}}

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
				strings.Join(c.input, "\n"),
				strings.Join(c.expected, "\n"),
				strings.Join(actual, "\n"))
		}
	}
}
