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
			"AddJavaOpts with port and no existing JAVA_OPTS",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			createTextFile().
				addLine("#!/bin/bash").
				toSlice(),
			createTextFile().
				addLine("#!/bin/bash").
				addLine("export JAVA_OPTS=\" \\").
				addLine(fmt.Sprintf("-Dhttp.proxyHost=%s \\", proxyHost)).
				addLine(fmt.Sprintf("-Dhttp.proxyPort=%s \\", proxyPort)).
				addLine(fmt.Sprintf("-Dhttps.proxyHost=%s \\", proxyHost)).
				addLine(fmt.Sprintf("-Dhttps.proxyPort=%s \\", proxyPort)).
				addLine(fmt.Sprintf("-Dhttp.nonProxyHosts=%s\"", strings.Join(nonProxyHosts, "|"))).
				toSlice(),
		},
		{
			"AddJavaOpts with port and existing JAVA_OPTS",
			proxyHost,
			proxyPort,
			nonProxyHosts,
			createTextFile().
				addLine("#!/bin/bash").
				addLine("export JAVA_OPTS=\" \\").
				addLine("-Djavax.net.ssl.trustStore=/etc/pki/truststore.jks \\").
				addLine("-Djavax.net.ssl.keyStore=/etc/pki/private/cert.p12 \\").
				addLine("-Djavax.net.ssl.keyStoreType=PKCS12\"").
				toSlice(),
			createTextFile().
				addLine("#!/bin/bash").
				addLine("export JAVA_OPTS=\" \\").
				addLine("-Djavax.net.ssl.trustStore=/etc/pki/truststore.jks \\").
				addLine("-Djavax.net.ssl.keyStore=/etc/pki/private/cert.p12 \\").
				addLine("-Djavax.net.ssl.keyStoreType=PKCS12 \\").
				addLine(fmt.Sprintf("-Dhttp.proxyHost=%s \\", proxyHost)).
				addLine(fmt.Sprintf("-Dhttp.proxyPort=%s \\", proxyPort)).
				addLine(fmt.Sprintf("-Dhttps.proxyHost=%s \\", proxyHost)).
				addLine(fmt.Sprintf("-Dhttps.proxyPort=%s \\", proxyPort)).
				addLine(fmt.Sprintf("-Dhttp.nonProxyHosts=%s\"", strings.Join(nonProxyHosts, "|"))).
				toSlice(),
		},
		{
			"AddJavaOpts without port and existing JAVA_OPTS",
			proxyHost,
			"",
			nonProxyHosts,
			createTextFile().
				addLine("#!/bin/bash").
				addLine("export JAVA_OPTS=\" \\").
				addLine("-Djavax.net.ssl.trustStore=/etc/pki/truststore.jks \\").
				addLine("-Djavax.net.ssl.keyStore=/etc/pki/private/cert.p12 \\").
				addLine("-Djavax.net.ssl.keyStoreType=PKCS12").
				toSlice(),
			createTextFile().
				addLine("#!/bin/bash").
				addLine("export JAVA_OPTS=\" \\").
				addLine("-Djavax.net.ssl.trustStore=/etc/pki/truststore.jks \\").
				addLine("-Djavax.net.ssl.keyStore=/etc/pki/private/cert.p12 \\").
				addLine("-Djavax.net.ssl.keyStoreType=PKCS12 \\").
				addLine(fmt.Sprintf("-Dhttp.proxyHost=%s \\", proxyHost)).
				addLine(fmt.Sprintf("-Dhttps.proxyHost=%s \\", proxyHost)).
				addLine(fmt.Sprintf("-Dhttp.nonProxyHosts=%s\"", strings.Join(nonProxyHosts, "|"))).
				toSlice(),
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

func TestParseShellStatements(t *testing.T) {
	cases := []struct {
		statements   []ShellStatement
		expectedFile string
	}{
		{
			createShellFile().
				addShellStatement(ShellStatement{[]string{"#!/bin/bash"}}).
				addShellStatement(ShellStatement{[]string{""}}).
				addShellStatement(ShellStatement{[]string{"export foo=bar"}}).
				addShellStatement(ShellStatement{
				[]string{
					"export multiline_foo=\"",
					"foo",
					"bar",
					"baz\"",
				}}).
				addShellStatement(ShellStatement{[]string{""}}).
				toSlice(),
			"_testdata/util/example_shell_file",
		},
	}
	for _, c := range cases {
		expectedFileContents := loadFileIntoSlice(c.expectedFile)
		actual := ParseShellStatements(c.statements)
		if !reflect.DeepEqual(actual, expectedFileContents) {
			t.Errorf(
				`Call:
ParseShellStatements() != {{expected}}

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
				c.statements,
				expectedFileContents,
				actual)
		}
	}
}

func TestParseShellContents(t *testing.T) {
	cases := []struct {
		inputFile string
		expected  []ShellStatement
	}{
		{
			"_testdata/util/example_shell_file",
			createShellFile().
				addShellStatement(ShellStatement{[]string{"#!/bin/bash"}}).
				addShellStatement(ShellStatement{[]string{""}}).
				addShellStatement(ShellStatement{[]string{"export foo=bar"}}).
				addShellStatement(ShellStatement{
				[]string{
					"export multiline_foo=\" ",
					"foo ",
					"bar ",
					"baz\"",
				}}).
				addShellStatement(ShellStatement{[]string{""}}).
				toSlice(),
		},
	}
	for _, c := range cases {
		fileContents := loadFileIntoSlice(c.inputFile)
		actual := ParseShellContents(fileContents)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`Call:
ParseShellContents() != {{expected}}

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
				fileContents,
				c.expected,
				actual)
		}
	}
}
