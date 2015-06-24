package proxy

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/clbanning/mxj"
)

func TestLoadXML(t *testing.T) {
	cases := []struct {
		name     string
		filename string
		expected mxj.Map
	}{
		{
			"loadXML",
			"example_xml_file.xml",
			mxj.Map(
				map[string]interface{}{
					"rootElement": map[string]interface{}{
						"childElement": []interface{}{
							map[string]interface{}{
								"-attribute1": "attribute 1 value",
								"-attribute2": "attribute 2 value",
								"#text":       "child element 1 value",
							},
							map[string]interface{}{
								"-attribute1": "attribute 1 value",
								"-attribute2": "attribute 2 value",
								"#text":       "child element 2 value",
							},
						},
					},
				},
			),
		},
	}

	for _, c := range cases {
		relativeFilename := fmt.Sprintf("%s/%s", "_testdata/util", c.filename)
		actual := loadXML(relativeFilename)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`%s
Call:
loadXML(%s) != {{expected}}

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
				c.expected,
				actual)
		}
	}
}

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
				addLine("export multiline_foo=\" \\").
				addLine("foo \\").
				addLine("bar \\").
				addLine("baz\"").
				addLine("").
				toSlice(),
		},
	}

	for _, c := range cases {
		relativeFilename := fmt.Sprintf("%s/%s", "_testdata/util", c.filename)
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
