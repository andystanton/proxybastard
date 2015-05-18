package uilayout

import "github.com/nsf/termbox-go"

// ListItem is as ListItem does.
type ListItem struct {
	label   string
	enabled bool
}

// ListBox is as ListBox does.
type ListBox struct {
	x            int
	y            int
	width        int
	widthStr     string
	height       int
	heightStr    string
	listItems    map[int]*ListItem
	selectedItem int
}

// CreateListBox creates a list box.
func CreateListBox(x, y int, widthStr string, heightStr string) *ListBox {
	lb := ListBox{
		x: x,
		y: y,
		listItems: map[int]*ListItem{
			0: {"env", true},
			1: {"git", false},
			2: {"npm", false},
		},
		selectedItem: 0,
		widthStr:     widthStr,
		heightStr:    heightStr,
	}
	return &lb
}

// ProcessEvent does what it says on the tin.
func (lb *ListBox) ProcessEvent(ev termbox.Event) {
	switch ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyArrowUp:
			lb.selectedItem--
			if lb.selectedItem < 0 {
				lb.selectedItem = 0
			}
		case termbox.KeyArrowDown:
			lb.selectedItem++
			if lb.selectedItem >= len(lb.listItems) {
				lb.selectedItem = len(lb.listItems) - 1
			}
		}
	}
}

// Update a ListBox on Resize
func (lb *ListBox) Update(width, height int, selected bool) {
	lb.width = parseDimensionString(lb.widthStr, lb.width, lb.height, width, height)
	lb.height = parseDimensionString(lb.heightStr, lb.width, lb.height, width, height)
}

// Draw a ListBox
func (lb *ListBox) Draw() {
	const coldef = termbox.ColorDefault

	drawBox(lb.x, lb.y, lb.width, lb.height)

	textFGColour := termbox.ColorDefault
	textBGColour := termbox.ColorDefault
	selectedTextFGColour := termbox.AttrBold
	selectedTextBGColour := termbox.ColorBlue

	for index, item := range lb.listItems {
		// ✔ ✘
		fgColour := textFGColour
		bgColour := textBGColour
		if index == lb.selectedItem {
			fgColour = selectedTextFGColour
			bgColour = selectedTextBGColour
		}
		printText(lb.x+2, lb.y+1+index, fgColour, bgColour, item.label)
	}

}
