package proxy

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestRemoveSubversionProxies(t *testing.T) {
	cases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			"removeSubversionProxies with existing global group",
			[]string{
				"[global]",
				"http-proxy-host=http://www.proxy-bastard.com",
				"http-proxy-port=80",
				"http-proxy-exceptions=foo,bar,127.0.0.1",
				"",
				"[something else]",
				"foo=bar",
				"",
			},
			[]string{
				"[global]",
				"",
				"[something else]",
				"foo=bar",
				"",
			},
		},
	}
	for _, c := range cases {
		actual := removeSubversionProxies(c.input)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`%s
Call:
removeSubversionProxies() != {{expected}}

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

func TestAddSubversionProxies(t *testing.T) {
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
			"addSubversionProxies with existing global group",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			[]string{
				"[global]",
				"",
				"[something else]",
				"foo=bar",
				"",
			},
			[]string{
				"[global]",
				fmt.Sprintf("http-proxy-host=%s", proxyHost),
				fmt.Sprintf("http-proxy-port=%s", proxyPort),
				fmt.Sprintf("http-proxy-exceptions=%s", strings.Join(nonProxyHosts, ",")),
				"",
				"[something else]",
				"foo=bar",
				"",
			},
		},
		{
			"addSubversionProxies without existing global group",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			[]string{
				"[something else]",
				"foo=bar",
				"",
			},
			[]string{
				"[something else]",
				"foo=bar",
				"",
				"[global]",
				fmt.Sprintf("http-proxy-host=%s", proxyHost),
				fmt.Sprintf("http-proxy-port=%s", proxyPort),
				fmt.Sprintf("http-proxy-exceptions=%s", strings.Join(nonProxyHosts, ",")),
				"",
			},
		},
	}
	for _, c := range cases {
		actual := addSubversionProxies(c.input, c.proxyHost, c.proxyPort, c.nonProxyHosts)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`%s
Call:
addSubversionProxies() != {{expected}}

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

func TestparseSubversionContents(t *testing.T) {
	cases := []struct {
		input    []string
		expected []SvnStatement
	}{
		{
			[]string{
				"[global]",
				"foo=bar",
				"",
				"[something else]",
				"foo=bar",
				"bar=foo",
				"",
			},
			[]SvnStatement{
				SvnStatement{
					"global",
					[]string{
						"foo=bar",
					},
				},
				SvnStatement{
					"something else",
					[]string{
						"foo=bar",
						"bar=foo",
					},
				},
			},
		},
		{
			[]string{
				"#",
				"##",
				"[section]",
				"###",
			},
			[]SvnStatement{
				SvnStatement{
					"section",
					[]string{},
				},
			},
		},
	}
	for _, c := range cases {
		actual := parseSubversionContents(c.input)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`Call:
parseSubversionContents() != {{expected}}

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
				c.input,
				c.expected,
				actual)
		}
	}
}

func TestParseSubversionStatements(t *testing.T) {
	cases := []struct {
		input    []SvnStatement
		expected []string
	}{
		{
			[]SvnStatement{
				SvnStatement{
					"global",
					[]string{
						"foo=bar",
					},
				},
				SvnStatement{
					"something else",
					[]string{
						"foo=bar",
						"bar=foo",
					},
				},
			},
			[]string{
				"[global]",
				"foo=bar",
				"",
				"[something else]",
				"foo=bar",
				"bar=foo",
				"",
			},
		},
	}
	for _, c := range cases {
		actual := parseSubversionStatements(c.input)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`Call:
ParseSubversionStatements() != {{expected}}

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
				c.input,
				c.expected,
				actual)
		}
	}
}
