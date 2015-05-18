package uilayout

import (
	"unicode/utf8"

	"github.com/nsf/termbox-go"
)

const editBoxWidth = 40

const preferredHorizontalThreshold = 5
const tabstopLength = 8

// EditBox is as EditBox does.
type EditBox struct {
	text          []byte
	x             int
	xStr          string
	y             int
	yStr          string
	width         int
	widthStr      string
	height        int
	heightStr     string
	lineVOffset   int
	cursorBOffset int // cursor offset in bytes
	cursorVOffset int // visual cursor offset in termbox cells
	cursorCOffset int // cursor offset in unicode code points
}

// CreateEditBox foo bar.
func CreateEditBox(xStr, yStr string) *EditBox {
	eb := EditBox{
		width:  editBoxWidth,
		height: 1,
		xStr:   xStr,
		yStr:   yStr,
	}
	return &eb
}

// ProcessEvent does what it says on the tin.
func (eb *EditBox) ProcessEvent(ev termbox.Event) {
	switch ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyArrowLeft, termbox.KeyCtrlB:
			eb.MoveCursorOneRuneBackward()
		case termbox.KeyArrowRight, termbox.KeyCtrlF:
			eb.MoveCursorOneRuneForward()
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			eb.DeleteRuneBackward()
		case termbox.KeyDelete, termbox.KeyCtrlD:
			eb.DeleteRuneForward()
		case termbox.KeySpace:
			eb.InsertRune(' ')
		case termbox.KeyCtrlK:
			eb.DeleteTheRestOfTheLine()
		case termbox.KeyHome, termbox.KeyCtrlA:
			eb.MoveCursorToBeginningOfTheLine()
		case termbox.KeyEnd, termbox.KeyCtrlE:
			eb.MoveCursorToEndOfTheLine()
		default:
			if ev.Ch != 0 {
				eb.InsertRune(ev.Ch)
			}
		}
	}
}

// Draw the EditBox in the given location, 'h' is not used at the moment
func (eb *EditBox) Draw() {
	eb.AdjustVOffset(eb.width)
	drawBox(eb.x-2, eb.y-1, editBoxWidth+3, 3)

	const coldef = termbox.ColorDefault
	fillArea(eb.x, eb.y, eb.width, eb.height, termbox.Cell{Ch: ' '})

	t := eb.text
	lx := 0
	tabstop := 0
	for {
		rx := lx - eb.lineVOffset
		if len(t) == 0 {
			break
		}

		if lx == tabstop {
			tabstop += tabstopLength
		}

		if rx >= eb.width {
			termbox.SetCell(eb.x+eb.width-1, eb.y, '→',
				coldef, coldef)
			break
		}

		r, size := utf8.DecodeRune(t)
		if r == '\t' {
			for ; lx < tabstop; lx++ {
				rx = lx - eb.lineVOffset
				if rx >= eb.width {
					goto next
				}

				if rx >= 0 {
					termbox.SetCell(eb.x+rx, eb.y, ' ', coldef, coldef)
				}
			}
		} else {
			if rx >= 0 {
				termbox.SetCell(eb.x+rx, eb.y, r, coldef, coldef)
			}
			lx++
		}
	next:
		t = t[size:]
	}

	if eb.lineVOffset != 0 {
		termbox.SetCell(eb.x, eb.y, '←', coldef, coldef)
	}
}

// Update an EditBox on Resize
func (eb *EditBox) Update(width, height int, selected bool) {
	eb.x = parseDimensionString(eb.xStr, eb.width, eb.height, width, height)
	eb.y = parseDimensionString(eb.yStr, eb.width, eb.height, width, height)

	if selected {
		termbox.SetCursor(eb.x+eb.CursorX(), eb.y)
	}
}

func runeAdvanceLen(r rune, pos int) int {
	if r == '\t' {
		return tabstopLength - pos%tabstopLength
	}
	return 1
}

func vOffsetCOffset(text []byte, boffset int) (voffset, coffset int) {
	text = text[:boffset]
	for len(text) > 0 {
		r, size := utf8.DecodeRune(text)
		text = text[size:]
		coffset++
		voffset += runeAdvanceLen(r, voffset)
	}
	return
}

func byteSliceGrow(s []byte, desiredCap int) []byte {
	if cap(s) < desiredCap {
		ns := make([]byte, len(s), desiredCap)
		copy(ns, s)
		return ns
	}
	return s
}

func byteSliceRemove(text []byte, from, to int) []byte {
	size := to - from
	copy(text[from:], text[to:])
	text = text[:len(text)-size]
	return text
}

func byteSliceInsert(text []byte, offset int, what []byte) []byte {
	n := len(text) + len(what)
	text = byteSliceGrow(text, n)
	text = text[:n]
	copy(text[offset+len(what):], text[offset:])
	copy(text[offset:], what)
	return text
}

// AdjustVOffset adjusts line visual offset to a proper value depending on width
func (eb *EditBox) AdjustVOffset(width int) {
	ht := preferredHorizontalThreshold
	maxHThreshold := (width - 1) / 2
	if ht > maxHThreshold {
		ht = maxHThreshold
	}

	threshold := width - 1
	if eb.lineVOffset != 0 {
		threshold = width - ht
	}
	if eb.cursorVOffset-eb.lineVOffset >= threshold {
		eb.lineVOffset = eb.cursorVOffset + (ht - width + 1)
	}

	if eb.lineVOffset != 0 && eb.cursorVOffset-eb.lineVOffset < ht {
		eb.lineVOffset = eb.cursorVOffset - ht
		if eb.lineVOffset < 0 {
			eb.lineVOffset = 0
		}
	}
}

// MoveCursorTo does what it says on the tin.
func (eb *EditBox) MoveCursorTo(boffset int) {
	eb.cursorBOffset = boffset
	eb.cursorVOffset, eb.cursorCOffset = vOffsetCOffset(eb.text, boffset)
}

// RuneUnderCursor does what it says on the tin.
func (eb *EditBox) RuneUnderCursor() (rune, int) {
	return utf8.DecodeRune(eb.text[eb.cursorBOffset:])
}

// RuneBeforeCursor does what it says on the tin.
func (eb *EditBox) RuneBeforeCursor() (rune, int) {
	return utf8.DecodeLastRune(eb.text[:eb.cursorBOffset])
}

// MoveCursorOneRuneBackward does what it says on the tin.
func (eb *EditBox) MoveCursorOneRuneBackward() {
	if eb.cursorBOffset == 0 {
		return
	}
	_, size := eb.RuneBeforeCursor()
	eb.MoveCursorTo(eb.cursorBOffset - size)
}

// MoveCursorOneRuneForward does what it says on the tin.
func (eb *EditBox) MoveCursorOneRuneForward() {
	if eb.cursorBOffset == len(eb.text) {
		return
	}
	_, size := eb.RuneUnderCursor()
	eb.MoveCursorTo(eb.cursorBOffset + size)
}

// MoveCursorToBeginningOfTheLine does what it says on the tin.
func (eb *EditBox) MoveCursorToBeginningOfTheLine() {
	eb.MoveCursorTo(0)
}

// MoveCursorToEndOfTheLine does what it says on the tin.
func (eb *EditBox) MoveCursorToEndOfTheLine() {
	eb.MoveCursorTo(len(eb.text))
}

// DeleteRuneBackward does what it says on the tin.
func (eb *EditBox) DeleteRuneBackward() {
	if eb.cursorBOffset == 0 {
		return
	}

	eb.MoveCursorOneRuneBackward()
	_, size := eb.RuneUnderCursor()
	eb.text = byteSliceRemove(eb.text, eb.cursorBOffset, eb.cursorBOffset+size)
}

// DeleteRuneForward does what it says on the tin.
func (eb *EditBox) DeleteRuneForward() {
	if eb.cursorBOffset == len(eb.text) {
		return
	}
	_, size := eb.RuneUnderCursor()
	eb.text = byteSliceRemove(eb.text, eb.cursorBOffset, eb.cursorBOffset+size)
}

// DeleteTheRestOfTheLine does what it says on the tin.
func (eb *EditBox) DeleteTheRestOfTheLine() {
	eb.text = eb.text[:eb.cursorBOffset]
}

// InsertRune does what it says on the tin.
func (eb *EditBox) InsertRune(r rune) {
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], r)
	eb.text = byteSliceInsert(eb.text, eb.cursorBOffset, buf[:n])
	eb.MoveCursorOneRuneForward()
}

// CursorX does what it says on the tin.
// Please, keep in mind that cursor depends on the value of lineVOffset, which
// is being set on Draw() call, so.. call this method after Draw() one.
func (eb *EditBox) CursorX() int {
	return eb.cursorVOffset - eb.lineVOffset
}
