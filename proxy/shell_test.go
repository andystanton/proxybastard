package proxy

import (
	"reflect"
	"testing"
)

func TestparseShellStatements(t *testing.T) {
	cases := []struct {
		statements []shellStatement
		expected   []string
	}{
		{
			[]shellStatement{
				shellStatement{[]string{"#!/bin/bash"}},
				shellStatement{[]string{""}},
				shellStatement{[]string{"export foo=bar"}},
				shellStatement{
					[]string{
						"export multiline_foo=\"",
						"foo",
						"bar",
						"baz\"",
					}},
				shellStatement{[]string{""}},
			},
			[]string{
				"#!/bin/bash",
				"",
				"export foo=bar",
				"export multiline_foo=\"\\",
				"foo \\",
				"bar \\",
				"baz\"",
				"",
			},
		},
	}
	for _, c := range cases {
		actual := parseShellStatements(c.statements)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`Call:
parseShellStatements() != {{expected}}

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

func TestparseShellContents(t *testing.T) {
	cases := []struct {
		input    []string
		expected []shellStatement
	}{
		{
			[]string{
				"#!/bin/bash",
				"",
				"export foo=bar",
				"export multiline_foo=\"\\",
				"foo \\",
				"bar \\",
				"baz\"",
				"",
			},
			[]shellStatement{
				shellStatement{[]string{"#!/bin/bash"}},
				shellStatement{[]string{""}},
				shellStatement{[]string{"export foo=bar"}},
				shellStatement{
					[]string{
						"export multiline_foo=\"",
						"foo ",
						"bar ",
						"baz\"",
					}},
				shellStatement{[]string{""}},
			},
		},
	}
	for _, c := range cases {
		actual := parseShellContents(c.input)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`Call:
parseShellContents() != {{expected}}

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
