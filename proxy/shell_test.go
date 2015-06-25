package proxy

import (
	"reflect"
	"testing"
)

func TestParseShellStatements(t *testing.T) {
	cases := []struct {
		statements []ShellStatement
		expected   []string
	}{
		{
			[]ShellStatement{
				ShellStatement{[]string{"#!/bin/bash"}},
				ShellStatement{[]string{""}},
				ShellStatement{[]string{"export foo=bar"}},
				ShellStatement{
					[]string{
						"export multiline_foo=\"",
						"foo",
						"bar",
						"baz\"",
					}},
				ShellStatement{[]string{""}},
			},
			[]string{
				"#!/bin/bash",
				"",
				"export foo=bar",
				"export multiline_foo=\" \\",
				"foo \\",
				"bar \\",
				"baz\"",
				"",
			},
		},
	}
	for _, c := range cases {
		actual := ParseShellStatements(c.statements)
		if !reflect.DeepEqual(actual, c.expected) {
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
				c.expected,
				actual)
		}
	}
}

func TestParseShellContents(t *testing.T) {
	cases := []struct {
		input    []string
		expected []ShellStatement
	}{
		{
			[]string{
				"#!/bin/bash",
				"",
				"export foo=bar",
				"export multiline_foo=\" \\",
				"foo \\",
				"bar \\",
				"baz\"",
				"",
			},
			[]ShellStatement{
				ShellStatement{[]string{"#!/bin/bash"}},
				ShellStatement{[]string{""}},
				ShellStatement{[]string{"export foo=bar"}},
				ShellStatement{
					[]string{
						"export multiline_foo=\" ",
						"foo ",
						"bar ",
						"baz\"",
					}},
				ShellStatement{[]string{""}},
			},
		},
	}
	for _, c := range cases {
		actual := ParseShellContents(c.input)
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
				c.input,
				c.expected,
				actual)
		}
	}
}
