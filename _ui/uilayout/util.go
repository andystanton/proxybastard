package uilayout

import (
	"fmt"
	"regexp"

	"github.com/marcmak/calc/calc"
	"github.com/nsf/termbox-go"
)

func drawBox(xOffset, yOffset, width, height int) {
	const coldef = termbox.ColorDefault
	termbox.SetCell(xOffset, yOffset, '┌', coldef, coldef)
	fillArea(xOffset+1, yOffset, width-2, 1, termbox.Cell{Ch: '─'})
	termbox.SetCell(xOffset+width-1, yOffset, '┐', coldef, coldef)
	fillArea(xOffset+width-1, yOffset+1, 1, height-2, termbox.Cell{Ch: '│'})
	termbox.SetCell(xOffset+width-1, yOffset+height-1, '┘', coldef, coldef)
	fillArea(xOffset+1, yOffset+height-1, width-2, 1, termbox.Cell{Ch: '─'})
	termbox.SetCell(xOffset, yOffset+height-1, '└', coldef, coldef)
	fillArea(xOffset, yOffset+1, 1, height-2, termbox.Cell{Ch: '│'})
}

func printText(x, y int, fg, bg termbox.Attribute, msg string) {
	dx := x
	for _, c := range msg {
		if c == '\n' {
			y++
			x = dx
		} else {
			termbox.SetCell(dx, y, c, fg, bg)
		}
		dx++
	}
}

func fillArea(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func parseDimensionString(dimensionString string,
	elementDimensionX int,
	elementDimensionY int,
	windowDimensionX int,
	windowDimensionY int) int {

	windowWidthR := regexp.MustCompile("{{window-width}}")
	elementWidthR := regexp.MustCompile("{{element-width}}")
	windowHeightR := regexp.MustCompile("{{window-height}}")
	elementHeightR := regexp.MustCompile("{{element-height}}")

	target := dimensionString

	if len(windowWidthR.FindStringSubmatch(target)) > 0 {
		target = fmt.Sprintf(windowWidthR.ReplaceAllString(target, "%d"), windowDimensionX)
	}

	if len(elementWidthR.FindStringSubmatch(target)) > 0 {
		target = fmt.Sprintf(elementWidthR.ReplaceAllString(target, "%d"), elementDimensionX)
	}

	if len(windowHeightR.FindStringSubmatch(target)) > 0 {
		target = fmt.Sprintf(windowHeightR.ReplaceAllString(target, "%d"), windowDimensionY)
	}

	if len(elementHeightR.FindStringSubmatch(target)) > 0 {
		target = fmt.Sprintf(elementHeightR.ReplaceAllString(target, "%d"), elementDimensionY)
	}

	return int(calc.Solve(target))
}
