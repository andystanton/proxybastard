package proxy

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestAddShellEnvVars(t *testing.T) {
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
			"addShellEnvVars",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			[]string{
				"#!/bin/bash",
			},
			[]string{
				"#!/bin/bash",
				fmt.Sprintf("export http_proxy=%s:%s", proxyHost, proxyPort),
				fmt.Sprintf("export https_proxy=%s:%s", proxyHost, proxyPort),
				fmt.Sprintf("export ALL_PROXY=%s:%s", proxyHost, proxyPort),
				fmt.Sprintf("export NO_PROXY=%s", strings.Join(nonProxyHosts, ",")),
			},
		},
	}
	for _, c := range cases {
		actual := addShellEnvVars(c.input, c.proxyHost, c.proxyPort, c.nonProxyHosts)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`%s
Call:
addShellEnvVars({{input}}, %s, %s) != {{expected}}

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

func TestRemoveShellEnvVars(t *testing.T) {
	proxyHost := "http://www.proxy-bastard.com"

	cases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			"removeShellEnvVars",
			[]string{
				"#!/bin/bash",
				fmt.Sprintf("export http_proxy=%s", proxyHost),
				fmt.Sprintf("export https_proxy=%s", proxyHost),
				fmt.Sprintf("export ALL_PROXY=%s", proxyHost),
			},
			[]string{
				"#!/bin/bash",
			},
		},
	}
	for _, c := range cases {
		actual := removeShellEnvVars(c.input)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`%s
Call:
removeShellEnvVars({{input}}) != {{expected}}

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
