package ui

import (
	"log"

	"github.com/andystanton/proxybastard/ui/uilayout"
	"github.com/nsf/termbox-go"
)

func hardFlush() {
	termbox.HideCursor()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()
}

func drawNode(node *uilayout.Node, w, h int, selectedElement string) {
	log.Printf("Drawing %s with %d children\n", node.ID(), len(node.Children()))
	if node.ID() != "root" {
		node.Element().Update(w, h, selectedElement == node.ID())
		node.Element().Draw()
	}
	for _, child := range node.Children() {
		drawNode(child, w, h, selectedElement)
	}
}

func redrawAll(selectedElement string, root *uilayout.Node) {
	w, h := termbox.Size()

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	drawNode(root, w, h, selectedElement)

	termbox.Flush()
}

func getNext(elements []string, current string) string {
	selectedIndex := -1
	for index, value := range elements {
		if current == value {
			selectedIndex = index
		}
	}
	next := ""
	switch selectedIndex {
	case -1:
		next = elements[0]
	case len(elements) - 1:
		next = elements[0]
	default:
		next = elements[selectedIndex+1]
	}
	return next
}

// UIMain is the entrypoint for running the UI.
func UIMain() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	rootNode := uilayout.CreateRootNode()

	uilayout.CreateNode("outerFrame", "root",
		uilayout.CreateFrame("{{window-width}}", "{{window-height}}", "0", "0"))

	// title

	uilayout.CreateNode("title-box", "root",
		uilayout.CreateFrame("{{window-width}}-4", "3", "2", "1"))

	uilayout.CreateNode("title-text", "title-box",
		uilayout.CreateText("({{window-width}}/2)-6", "2", "PROXY BASTARD"))

	// list

	uilayout.CreateNode("list-box", "root",
		uilayout.CreateListBox(2, 4, "30", "{{window-height}}-5"))

	// main content

	uilayout.CreateNode("content-box", "root",
		uilayout.CreateFrame("{{window-width}}-35", "{{window-height}}-5", "33", "4"))

	uilayout.CreateNode("shell-file-text", "content-box",
		uilayout.CreateText("35", "6", "Shell file:"))

	uilayout.CreateNode("shell-edit-box", "content-box",
		uilayout.CreateEditBox("{{window-width}}-{{element-width}}-5", "6"))

	uilayout.CreateNode("shell-file-description-text", "content-box",
		uilayout.CreateText("{{window-width}}-{{element-width}}-4", "9", "(Location of env vars e.g. ~/.bash_profile ~/.zshrc)"))

	uilayout.CreateNode("other-thing-text", "content-box",
		uilayout.CreateText("35", "13", "Something else:"))

	uilayout.CreateNode("other-edit-box", "content-box",
		uilayout.CreateEditBox("{{window-width}}-{{element-width}}-5", "13"))

	uilayout.CreateNode("other-thing-description-text", "content-box",
		uilayout.CreateText("{{window-width}}-{{element-width}}-4", "16", "(Some other description containing detailed information)"))

	uilayout.CreateNode("escape-text", "content-box",
		uilayout.CreateText("({{window-width}}/2)-9", "{{window-height}}-3", "Press ESC to quit"))

	// selectable bits

	Selectables := []string{"list-box", "shell-edit-box", "other-edit-box"}
	selectedElement := "list-box"

	redrawAll(selectedElement, rootNode)
mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			case termbox.KeyTab:
				selectedElement = getNext(Selectables, selectedElement)
				hardFlush()
			}
			selectable := uilayout.FindElementNodeByID(selectedElement).Element().(uilayout.Selectable)
			selectable.ProcessEvent(ev)
		case termbox.EventError:
			log.Println(ev.Err)
			panic(ev.Err)
		}
		redrawAll(selectedElement, rootNode)
	}
}
