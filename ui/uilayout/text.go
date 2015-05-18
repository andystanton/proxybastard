package uilayout

import (
	"log"

	"github.com/nsf/termbox-go"
)

// Text is as Text does.
type Text struct {
	x    int
	xStr string
	y    int
	yStr string
	text string
}

// CreateText foo bar.
func CreateText(xStr, yStr, text string) *Text {
	txt := Text{
		xStr: xStr,
		yStr: yStr,
		text: text,
	}
	return &txt
}

// Update a Text on Resize
func (txt *Text) Update(width, height int, selected bool) {
	txt.x = parseDimensionString(txt.xStr, len(txt.text), 1, width, height)
	txt.y = parseDimensionString(txt.yStr, len(txt.text), 1, width, height)
	log.Printf("x: %d, y: %d\n", txt.x, txt.y)
}

// Draw a Text
func (txt *Text) Draw() {
	const coldef = termbox.ColorDefault

	fgColour := termbox.ColorDefault
	bgColour := termbox.ColorDefault
	log.Printf("x: %d, y: %d\n", txt.x, txt.y)
	printText(txt.x, txt.y, fgColour, bgColour, txt.text)
}
