package proxy

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestAddProxyVars(t *testing.T) {
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
			"AddProxyVars with port",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			createTextFile().
				addLine("#!/bin/bash").
				toSlice(),
			createTextFile().
				addLine("#!/bin/bash").
				addLine(fmt.Sprintf("export http_proxy=%s:%s", proxyHost, proxyPort)).
				addLine(fmt.Sprintf("export https_proxy=%s:%s", proxyHost, proxyPort)).
				addLine(fmt.Sprintf("export ALL_PROXY=%s:%s", proxyHost, proxyPort)).
				addLine(fmt.Sprintf("export NO_PROXY=%s", strings.Join(nonProxyHosts, ","))).
				toSlice(),
		},
		{
			"AddProxyVars without port",
			proxyHost,
			"",
			nonProxyHosts,
			createTextFile().
				addLine("#!/bin/bash").
				toSlice(),
			createTextFile().
				addLine("#!/bin/bash").
				addLine(fmt.Sprintf("export http_proxy=%s", proxyHost)).
				addLine(fmt.Sprintf("export https_proxy=%s", proxyHost)).
				addLine(fmt.Sprintf("export ALL_PROXY=%s", proxyHost)).
				addLine(fmt.Sprintf("export NO_PROXY=%s", strings.Join(nonProxyHosts, ","))).
				toSlice(),
		},
	}
	for _, c := range cases {
		actual := AddProxyVars(c.input, c.proxyHost, c.proxyPort, c.nonProxyHosts)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`%s
Call:
AddProxyVars({{input}}, %s, %s) != {{expected}}

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

func TestRemoveProxyVars(t *testing.T) {
	proxyHost := "http://www.proxy-bastard.com"

	cases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			"RemoveProxyVars",
			createTextFile().
				addLine("#!/bin/bash").
				addLine(fmt.Sprintf("export http_proxy=%s", proxyHost)).
				addLine(fmt.Sprintf("export https_proxy=%s", proxyHost)).
				addLine(fmt.Sprintf("export ALL_PROXY=%s", proxyHost)).
				toSlice(),
			createTextFile().
				addLine("#!/bin/bash").
				toSlice(),
		},
	}
	for _, c := range cases {
		actual := RemoveProxyVars(c.input)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`%s
Call:
RemoveProxyVars({{input}}) != {{expected}}

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
