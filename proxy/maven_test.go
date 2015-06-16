package proxy

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAddProxyVarsMaven(t *testing.T) {
	proxyHost := "http://www.proxy-bastard.com"
	proxyPort := 80
	nonProxyHosts := []string{"localhost", "127.0.0.1"}

	cases := []struct {
		name          string
		proxyHost     string
		proxyPort     *int
		nonProxyHosts []string
		inputFile     string
		expectedFile  string
	}{
		{
			"AddProxyVarsMaven with existing proxy settings disabled",
			proxyHost,
			&proxyPort,
			nonProxyHosts,
			"_testdata/maven/maven_settings_disabled.xml",
			"_testdata/maven/maven_settings_enabled.xml",
		},
		{
			"AddProxyVarsMaven with no existing proxy settings",
			proxyHost,
			&proxyPort,
			nonProxyHosts,
			"_testdata/maven/maven_settings_missing.xml",
			"_testdata/maven/maven_settings_enabled.xml",
		},
	}
	for _, c := range cases {
		actual := AddProxyVarsMaven(loadXML(c.inputFile), c.proxyHost, c.proxyPort, c.nonProxyHosts)
		if !reflect.DeepEqual(actual, loadXML(c.expectedFile)) {
			resolvedPort := "nil"
			if c.proxyPort != nil {
				resolvedPort = fmt.Sprintf("%d", *c.proxyPort)
			}
			t.Errorf(
				`%s
Call:
AddProxyVarsMaven({{input}}, %s, %s) != {{expected}}

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
				resolvedPort,
				loadXML(c.inputFile),
				loadXML(c.expectedFile),
				actual)
		}
	}
}

func TestRemoveProxyVarsMaven(t *testing.T) {
	proxyHost := "http://www.proxy-bastard.com"
	proxyPort := 80
	nonProxyHosts := []string{"localhost", "127.0.0.1"}

	cases := []struct {
		name          string
		proxyHost     string
		proxyPort     *int
		nonProxyHosts []string
		inputFile     string
		expectedFile  string
	}{
		{
			"RemoveProxyVarsMaven with existing proxy settings enabled",
			proxyHost,
			&proxyPort,
			nonProxyHosts,
			"_testdata/maven/maven_settings_disabled.xml",
			"_testdata/maven/maven_settings_disabled.xml",
		},
		{
			"RemoveProxyVarsMaven with no existing proxy settings",
			proxyHost,
			&proxyPort,
			nonProxyHosts,
			"_testdata/maven/maven_settings_missing.xml",
			"_testdata/maven/maven_settings_disabled.xml",
		},
	}
	for _, c := range cases {
		actual := RemoveProxyVarsMaven(loadXML(c.inputFile), c.proxyHost, c.proxyPort, c.nonProxyHosts)
		if !reflect.DeepEqual(actual, loadXML(c.expectedFile)) {
			resolvedPort := "nil"
			if c.proxyPort != nil {
				resolvedPort = fmt.Sprintf("%d", *c.proxyPort)
			}
			t.Errorf(
				`%s
Call:
RemoveProxyVarsMaven({{input}}, %s, %s) != {{expected}}

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
				resolvedPort,
				loadXML(c.inputFile),
				loadXML(c.expectedFile),
				actual)
		}
	}
}
