package proxy

import (
	"reflect"
	"testing"
)

func TestAddEnvVarsMaven(t *testing.T) {
	proxyHost := "http://www.proxy-bastard.com"
	proxyPort := "80"
	nonProxyHosts := []string{"localhost", "127.0.0.1"}

	cases := []struct {
		name          string
		proxyHost     string
		proxyPort     string
		nonProxyHosts []string
		inputFile     string
		expectedFile  string
	}{
		{
			"AddEnvVarsMaven with existing proxy settings disabled",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			"_testdata/maven/maven_settings_disabled.xml",
			"_testdata/maven/maven_settings_enabled.xml",
		},
		{
			"AddEnvVarsMaven with no existing proxy settings",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			"_testdata/maven/maven_settings_missing.xml",
			"_testdata/maven/maven_settings_enabled.xml",
		},
	}
	for _, c := range cases {
		actual := AddEnvVarsMaven(loadXML(c.inputFile), c.proxyHost, c.proxyPort, c.nonProxyHosts)
		if !reflect.DeepEqual(actual, loadXML(c.expectedFile)) {
			t.Errorf(
				`%s
Call:
AddEnvVarsMaven({{input}}, %s, %s) != {{expected}}

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
				loadXML(c.inputFile),
				loadXML(c.expectedFile),
				actual)
		}
	}
}

func TestRemoveEnvVarsMaven(t *testing.T) {
	proxyHost := "http://www.proxy-bastard.com"
	proxyPort := "80"
	nonProxyHosts := []string{"localhost", "127.0.0.1"}

	cases := []struct {
		name          string
		proxyHost     string
		proxyPort     string
		nonProxyHosts []string
		inputFile     string
		expectedFile  string
	}{
		{
			"RemoveEnvVarsMaven with existing proxy settings enabled",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			"_testdata/maven/maven_settings_disabled.xml",
			"_testdata/maven/maven_settings_disabled.xml",
		},
		{
			"RemoveEnvVarsMaven with no existing proxy settings",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			"_testdata/maven/maven_settings_missing.xml",
			"_testdata/maven/maven_settings_disabled.xml",
		},
	}
	for _, c := range cases {
		actual := RemoveEnvVarsMaven(loadXML(c.inputFile), c.proxyHost, c.proxyPort, c.nonProxyHosts)
		if !reflect.DeepEqual(actual, loadXML(c.expectedFile)) {
			t.Errorf(
				`%s
Call:
RemoveEnvVarsMaven({{input}}, %s, %s) != {{expected}}

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
				loadXML(c.inputFile),
				loadXML(c.expectedFile),
				actual)
		}
	}
}
