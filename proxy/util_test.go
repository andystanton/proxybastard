package proxy

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestLoadFileIntoSlice(t *testing.T) {
	cases := []struct {
		name     string
		filename string
		expected []string
	}{
		{
			"loadFileIntoSlice",
			"example_shell_file",
			createTextFile().
				addLine("#!/bin/bash").
				addLine("").
				addLine("export foo=bar").
				addLine("").
				toSlice(),
		},
	}

	for _, c := range cases {
		relativeFilename := fmt.Sprintf("%s/%s", "_testdata", c.filename)
		actual := loadFileIntoSlice(relativeFilename)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`%s
Call:
LoadFileIntoSlice(%s) != {{expected}}

Expected:
===============
%s
===============

Actual:
===============
%s
===============`,
				c.name,
				relativeFilename,
				strings.Join(c.expected, "\n"),
				strings.Join(actual, "\n"))
		}
	}
}
