package proxy

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestAddBowerProxySettings(t *testing.T) {
	proxyHost := "http://www.proxy-bastard.com"
	proxyPort := "80"
	nonProxyHosts := []string{"localhost", "127.0.0.1", "::1"}

	cases := []struct {
		name          string
		proxyHost     string
		proxyPort     string
		nonProxyHosts []string
		input         string
		expected      string
	}{
		{
			"addBowerProxySettings",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			`
			{
				"foo": "bar"
			}
			`,
			`
			{
				"foo": "bar",
			` +
				fmt.Sprintf("\"proxy\": \"%s:%s\",", proxyHost, proxyPort) +
				fmt.Sprintf("\"https-proxy\": \"%s:%s\"", proxyHost, proxyPort) +
				`
			}
			`,
		},
	}
	for _, c := range cases {
		inMap := make(map[string]interface{})
		expectedMap := make(map[string]interface{})
		json.Unmarshal([]byte(c.input), &inMap)
		json.Unmarshal([]byte(c.expected), &expectedMap)
		actual := addBowerProxySettings(inMap, c.proxyHost, c.proxyPort)
		if !reflect.DeepEqual(actual, expectedMap) {
			t.Errorf(
				`%s
Call:
addBowerProxySettings({{input}}, %s, %s) != {{expected}}

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
				inMap,
				expectedMap,
				actual)
		}
	}
}

func TestRemoveBowerProxySettings(t *testing.T) {
	proxyHost := "http://www.proxy-bastard.com"
	proxyPort := "80"

	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"removeBowerProxySettings",
			`
			{
				"foo": "bar",
			` +
				fmt.Sprintf("\"proxy\": \"%s:%s\",", proxyHost, proxyPort) +
				fmt.Sprintf("\"https-proxy\": \"%s:%s\",", proxyHost, proxyPort) +
				`
			"bar": "foo"
			}
			`,
			`
			{
				"foo": "bar",
				"bar": "foo"
			}
			`,
		},
	}
	for _, c := range cases {
		inMap := make(map[string]interface{})
		expectedMap := make(map[string]interface{})
		json.Unmarshal([]byte(c.input), &inMap)
		json.Unmarshal([]byte(c.expected), &expectedMap)
		actual := removeBowerProxySettings(inMap)
		if !reflect.DeepEqual(actual, expectedMap) {
			t.Errorf(
				`%s
Call:
removeBowerProxySettings({{input}}) != {{expected}}

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
				inMap,
				expectedMap,
				actual)
		}
	}
}
