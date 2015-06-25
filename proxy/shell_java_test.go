package proxy

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParseJavaOpts(t *testing.T) {
	cases := []struct {
		name       string
		javaOpts   []string
		parsedOpts []string
	}{
		{
			"Single line JAVA_OPTS",
			[]string{"export JAVA_OPTS=\"-Dfoo=bar -Dfoo2=bar2\""},
			[]string{
				"-Dfoo=bar",
				"-Dfoo2=bar2",
			},
		},
		{
			"Multi line JAVA_OPTS",
			[]string{
				"export JAVA_OPTS=\"\\",
				"-Dfoo=bar \\",
				"-Dfoo2=bar2 -Dfoo3=bar3\"",
			},
			[]string{
				"-Dfoo=bar",
				"-Dfoo2=bar2",
				"-Dfoo3=bar3",
			},
		},
		{
			"Another multi line JAVA_OPTS",
			[]string{
				"export JAVA_OPTS=\"-Dfoo=bar \\",
				"-Dfoo2=bar2 -Dfoo3=bar3 \\",
				"-Dfoo4=bar4\"",
			},
			[]string{
				"-Dfoo=bar",
				"-Dfoo2=bar2",
				"-Dfoo3=bar3",
				"-Dfoo4=bar4",
			},
		},
	}
	for _, c := range cases {
		actual := ParseJavaOpts(c.javaOpts)
		if !reflect.DeepEqual(actual, c.parsedOpts) {
			t.Errorf(
				`%s
Call:
ParseJavaOpts() != {{expected}}

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
				c.javaOpts,
				c.parsedOpts,
				actual)
		}
	}
}

func TestAddJavaOpts(t *testing.T) {
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
			"AddJavaOpts with no existing JAVA_OPTS",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			[]string{
				"#!/bin/bash",
			},
			[]string{
				"#!/bin/bash",
				"export JAVA_OPTS=\" \\",
				fmt.Sprintf("-Dhttp.proxyHost=%s \\", proxyHost),
				fmt.Sprintf("-Dhttp.proxyPort=%s \\", proxyPort),
				fmt.Sprintf("-Dhttps.proxyHost=%s \\", proxyHost),
				fmt.Sprintf("-Dhttps.proxyPort=%s \\", proxyPort),
				fmt.Sprintf("-Dhttp.nonProxyHosts=%s\"", strings.Join(nonProxyHosts, "|")),
			},
		},
		{
			"AddJavaOpts with existing JAVA_OPTS",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			[]string{
				"#!/bin/bash",
				"export JAVA_OPTS=\" \\",
				"-Djavax.net.ssl.trustStore=/etc/pki/truststore.jks \\",
				"-Djavax.net.ssl.keyStore=/etc/pki/private/cert.p12 \\",
				"-Djavax.net.ssl.keyStoreType=PKCS12\"",
			},
			[]string{
				"#!/bin/bash",
				"export JAVA_OPTS=\" \\",
				"-Djavax.net.ssl.trustStore=/etc/pki/truststore.jks \\",
				"-Djavax.net.ssl.keyStore=/etc/pki/private/cert.p12 \\",
				"-Djavax.net.ssl.keyStoreType=PKCS12 \\",
				fmt.Sprintf("-Dhttp.proxyHost=%s \\", proxyHost),
				fmt.Sprintf("-Dhttp.proxyPort=%s \\", proxyPort),
				fmt.Sprintf("-Dhttps.proxyHost=%s \\", proxyHost),
				fmt.Sprintf("-Dhttps.proxyPort=%s \\", proxyPort),
				fmt.Sprintf("-Dhttp.nonProxyHosts=%s\"", strings.Join(nonProxyHosts, "|")),
			},
		},
	}
	for _, c := range cases {
		actual := AddJavaOpts(c.input, c.proxyHost, c.proxyPort, c.nonProxyHosts)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`%s
Call:
AddJavaOpts() != {{expected}}

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
