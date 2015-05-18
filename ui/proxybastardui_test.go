package ui

import "testing"

func TestGetNext(t *testing.T) {
	cases := []struct {
		elements []string
		current  string
		want     string
	}{
		{
			[]string{"a", "b", "c"},
			"a",
			"b",
		},
		{
			[]string{"a", "b", "c"},
			"b",
			"c",
		},
		{
			[]string{"a", "b", "c"},
			"c",
			"a",
		},
		{
			[]string{"a", "b", "c"},
			"",
			"a",
		},
	}
	for _, c := range cases {
		next := getNext(c.elements, c.current)
		if next != c.want {
			t.Errorf("getNext(%q, %s) != %s", c.elements, c.current, c.want)
		}
	}
}
