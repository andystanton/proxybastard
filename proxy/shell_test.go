package proxy

import (
	"reflect"
	"testing"
)

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
