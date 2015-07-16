package proxy

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestAddNPMProxySettings(t *testing.T) {
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
			"addNPMProxySettings",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			[]string{
				"foo=bar",
			},
			[]string{
				"foo=bar",
				fmt.Sprintf("proxy=%s:%s", proxyHost, proxyPort),
				fmt.Sprintf("https-proxy=%s:%s", proxyHost, proxyPort),
			},
		},
	}
	for _, c := range cases {
		actual := addNPMProxySettings(c.input, c.proxyHost, c.proxyPort)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`%s
Call:
addNPMProxySettings({{input}}, %s, %s) != {{expected}}

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

func TestRemoveNPMProxySettings(t *testing.T) {
	proxyHost := "http://www.proxy-bastard.com"
	proxyPort := "80"

	cases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			"removeNPMProxySettings",
			[]string{
				"foo=bar",
				fmt.Sprintf("proxy=%s:%s", proxyHost, proxyPort),
				fmt.Sprintf("https-proxy=%s:%s", proxyHost, proxyPort),
				"bar=foo",
			},
			[]string{
				"foo=bar",
				"bar=foo",
			},
		},
	}
	for _, c := range cases {
		actual := removeNPMProxySettings(c.input)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`%s
Call:
removeNPMProxySettings({{input}}) != {{expected}}

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
