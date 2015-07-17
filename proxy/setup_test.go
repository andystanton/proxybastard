package proxy

import "testing"

func TestGetHighestFrequency(t *testing.T) {
	cases := []struct {
		frequencyMap map[string]int
		expected     string
	}{
		{
			map[string]int{
				"a": 10,
				"b": 8,
				"c": -2,
			},
			"a",
		},
		{
			map[string]int{
				"a": 1,
				"b": 7,
				"c": 6,
			},
			"b",
		},
		{
			map[string]int{
				"a": 3,
				"b": 3,
				"c": 3,
			},
			"a",
		},
	}

	for _, c := range cases {
		if result := getHighestFrequency(c.frequencyMap); result != c.expected {
			t.Errorf("getHighestFrequency(%v) == %s did not equal %s", c.frequencyMap, result, c.expected)
		}
	}
}
