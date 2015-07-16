package util

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/clbanning/mxj"
)

func TestGenerateRandomFilename(t *testing.T) {
	cases := []struct {
		filepath string
		filename string
	}{
		{
			"/path/to/some/",
			"file",
		},
		{
			"",
			"file",
		},
	}

	uuidFormat := "[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}"

	for _, c := range cases {
		safeFilename := generateRandomFilename(c.filepath + c.filename)
		safeFilenameFormat := strings.Replace(c.filepath, "/", "\\/", -1) + uuidFormat
		safeFilenameRegexp := regexp.MustCompile(safeFilenameFormat)
		if !safeFilenameRegexp.MatchString(safeFilename) {
			t.Errorf("randomFilename %s did not match the expected format %s", safeFilename, safeFilenameFormat)
		}
	}
}

func TestLoadXML(t *testing.T) {
	cases := []struct {
		name     string
		filename string
		expected mxj.Map
	}{
		{
			"LoadXML",
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
		relativeFilename := fmt.Sprintf("%s/%s", "_testdata", c.filename)
		actual := LoadXML(relativeFilename)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				`%s
Call:
LoadXML(%s) != {{expected}}

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
			"LoadFileIntoSlice",
			"example_shell_file",
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
		relativeFilename := fmt.Sprintf("%s/%s", "_testdata", c.filename)
		actual, _ := LoadFileIntoSlice(relativeFilename)
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
