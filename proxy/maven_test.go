package proxy

import (
	"reflect"
	"testing"

	"github.com/andystanton/proxybastard/util"
)

func TestAddToMavenXML(t *testing.T) {
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
			"addToMavenXML with existing proxy settings disabled",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			"_testdata/maven/maven_settings_disabled.xml",
			"_testdata/maven/maven_settings_enabled.xml",
		},
		{
			"addToMavenXML with no existing proxy settings",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			"_testdata/maven/maven_settings_missing.xml",
			"_testdata/maven/maven_settings_enabled.xml",
		},
	}
	for _, c := range cases {
		actual := addToMavenXML(util.LoadXML(c.inputFile), c.proxyHost, c.proxyPort, c.nonProxyHosts)
		if !reflect.DeepEqual(actual, util.LoadXML(c.expectedFile)) {
			t.Errorf(
				`%s
Call:
addToMavenXML({{input}}, %s, %s) != {{expected}}

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
				util.LoadXML(c.inputFile),
				util.LoadXML(c.expectedFile),
				actual)
		}
	}
}

func TestRemoveFromMavenXML(t *testing.T) {
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
			"removeFromMavenXML with existing proxy settings enabled",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			"_testdata/maven/maven_settings_enabled.xml",
			"_testdata/maven/maven_settings_missing.xml",
		},
		{
			"removeFromMavenXML with no existing proxy settings",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			"_testdata/maven/maven_settings_missing.xml",
			"_testdata/maven/maven_settings_missing.xml",
		},
	}
	for _, c := range cases {
		actual := removeFromMavenXML(util.LoadXML(c.inputFile))
		if !reflect.DeepEqual(actual, util.LoadXML(c.expectedFile)) {
			t.Errorf(
				`%s
Call:
removeFromMavenXML({{input}}, %s, %s) != {{expected}}

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
				util.LoadXML(c.inputFile),
				util.LoadXML(c.expectedFile),
				actual)
		}
	}
}
