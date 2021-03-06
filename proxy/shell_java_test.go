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
		actual := parseJavaOpts(c.javaOpts)
		if !reflect.DeepEqual(actual, c.parsedOpts) {
			t.Errorf(
				`%s
Call:
parseJavaOpts() != {{expected}}

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
	proxyHostNoProtocol := "www.proxy-bastard.com"
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
			"addJavaOpts with no existing JAVA_OPTS",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			[]string{
				"#!/bin/bash",
			},
			[]string{
				"#!/bin/bash",
				"export JAVA_OPTS=\"\\",
				fmt.Sprintf("-Dhttp.proxyHost=%s \\", proxyHostNoProtocol),
				fmt.Sprintf("-Dhttp.proxyPort=%s \\", proxyPort),
				fmt.Sprintf("-Dhttps.proxyHost=%s \\", proxyHostNoProtocol),
				fmt.Sprintf("-Dhttps.proxyPort=%s \\", proxyPort),
				fmt.Sprintf("-Dhttp.nonProxyHosts=\\\"%s\\\"\"", strings.Join(nonProxyHosts, "|")),
			},
		},
		{
			"addJavaOpts with existing JAVA_OPTS",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			[]string{
				"#!/bin/bash",
				"export JAVA_OPTS=\"\\",
				"-Djavax.net.ssl.trustStore=/etc/pki/truststore.jks \\",
				"-Djavax.net.ssl.keyStore=/etc/pki/private/cert.p12 \\",
				"-Djavax.net.ssl.keyStoreType=PKCS12\"",
			},
			[]string{
				"#!/bin/bash",
				"export JAVA_OPTS=\"\\",
				"-Djavax.net.ssl.trustStore=/etc/pki/truststore.jks \\",
				"-Djavax.net.ssl.keyStore=/etc/pki/private/cert.p12 \\",
				"-Djavax.net.ssl.keyStoreType=PKCS12 \\",
				fmt.Sprintf("-Dhttp.proxyHost=%s \\", proxyHostNoProtocol),
				fmt.Sprintf("-Dhttp.proxyPort=%s \\", proxyPort),
				fmt.Sprintf("-Dhttps.proxyHost=%s \\", proxyHostNoProtocol),
				fmt.Sprintf("-Dhttps.proxyPort=%s \\", proxyPort),
				fmt.Sprintf("-Dhttp.nonProxyHosts=\\\"%s\\\"\"", strings.Join(nonProxyHosts, "|")),
			},
		},
	}
	for _, c := range cases {
		actual := addJavaOpts(c.input, c.proxyHost, c.proxyPort, c.nonProxyHosts)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`%s
Call:
addJavaOpts() != {{expected}}

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

func TestRemoveJavaOpts(t *testing.T) {
	cases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			"removeJavaOpts with no existing JAVA_OPTS",
			[]string{
				"#!/bin/bash",
			},
			[]string{
				"#!/bin/bash",
			},
		},
		{
			"removeJavaOpts with existing JAVA_OPTS not containing proxy",
			[]string{
				"#!/bin/bash",
				"export JAVA_OPTS=\"\\",
				"-Dfoo=bar\"",
			},
			[]string{
				"#!/bin/bash",
				"export JAVA_OPTS=\"\\",
				"-Dfoo=bar\"",
			},
		},
		{
			"removeJavaOpts with existing JAVA_OPTS not containing proxy",
			[]string{
				"#!/bin/bash",
				"export JAVA_OPTS=\"\\",
				"-Dhttp.proxyHost=foo \\",
				"-Dhttp.proxyPort=bar \\",
				"-Dhttps.proxyHost=foo \\",
				"-Dhttps.proxyPort=bar \\",
				"-Dhttp.nonProxyHosts=\\\"any|thing|at|all\\\" \\",
				"-Dfoo=bar\"",
			},
			[]string{
				"#!/bin/bash",
				"export JAVA_OPTS=\"\\",
				"-Dfoo=bar\"",
			},
		},
	}
	for _, c := range cases {
		actual := removeJavaOpts(c.input)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`%s
Call:
removeJavaOpts() != {{expected}}

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

func TestExtractJavaOpts(t *testing.T) {
	cases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			"extractJavaOpts with no existing JAVA_OPTS",
			[]string{
				"#!/bin/bash",
			},
			[]string{},
		},
		{
			"extractJavaOpts with existing JAVA_OPTS not containing proxy",
			[]string{
				"#!/bin/bash",
				"export JAVA_OPTS=\"\\",
				"-Dfoo=bar\"",
			},
			[]string{
				"export JAVA_OPTS=\"\\",
				"-Dfoo=bar\"",
			},
		},
		{
			"extractJavaOpts with existing JAVA_OPTS not containing proxy",
			[]string{
				"#!/bin/bash",
				"export JAVA_OPTS=\"\\",
				"-Dhttp.proxyHost=foo \\",
				"-Dhttp.proxyPort=bar \\",
				"-Dhttps.proxyHost=foo \\",
				"-Dhttps.proxyPort=bar \\",
				"-Dhttp.nonProxyHosts=\\\"any|thing|at|all\\\" \\",
				"-Dfoo=bar\"",
			},
			[]string{
				"export JAVA_OPTS=\"\\",
				"-Dhttp.proxyHost=foo \\",
				"-Dhttp.proxyPort=bar \\",
				"-Dhttps.proxyHost=foo \\",
				"-Dhttps.proxyPort=bar \\",
				"-Dhttp.nonProxyHosts=\\\"any|thing|at|all\\\" \\",
				"-Dfoo=bar\"",
			},
		},
	}
	for _, c := range cases {
		actual := extractJavaOpts(c.input)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`%s
Call:
extractJavaOpts() != {{expected}}

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
