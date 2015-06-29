package uilayout

import (
	"testing"
)

type vec2 struct {
	x int
	y int
}

func TestParseDimensionString(t *testing.T) {
	cases := []struct {
		dimensionString   string
		elementDimensions vec2
		windowDimensions  vec2
		expected          int
	}{
		{
			"{{window-width}}-{{element-width}}-3",
			vec2{10, 1},
			vec2{72, 48},
			59,
		},
		{
			"{{window-height}}-{{element-height}}-3",
			vec2{10, 1},
			vec2{72, 48},
			44,
		},
	}
	for _, c := range cases {
		actual := parseDimensionString(c.dimensionString,
			c.elementDimensions.x, c.elementDimensions.y,
			c.windowDimensions.x, c.windowDimensions.y)
		if actual != c.expected {
			t.Errorf("parseDimensionString(%s, %d, %d, %d, %d) != %d",
				c.dimensionString,
				c.elementDimensions.x, c.elementDimensions.y,
				c.windowDimensions.x, c.elementDimensions.y,
				c.expected)
		}
	}
}
