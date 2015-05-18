package uilayout

import (
	"log"

	"github.com/nsf/termbox-go"
)

// Element interface
type Element interface {
	Update(x, y int, selectedElement bool)
	Draw()
}

// Selectable interface
type Selectable interface {
	ProcessEvent(e termbox.Event)
}

// Node is a node in a tree of Elements.
type Node struct {
	id       string
	element  Element
	parent   *Node
	children []*Node
}

var nodeDictionary = make(map[string]*Node)

// FindElementNodeByID returns an ElementNode by its id.
func FindElementNodeByID(id string) *Node {
	return nodeDictionary[id]
}

// CreateRootNode creates the root Node.
func CreateRootNode() *Node {
	node := Node{
		id: "root",
	}
	nodeDictionary[node.id] = &node
	return &node
}

// CreateNode creates a Node.
func CreateNode(id, parent string, element Element) *Node {
	if id == "root" {
		panic("'root' is a reserved id")
	}
	node := Node{
		id:      id,
		element: element,
		parent:  nodeDictionary[parent],
	}
	nodeDictionary[id] = &node
	log.Printf("Creating node %s with parent %s\n", id, nodeDictionary[parent].ID())
	nodeDictionary[parent].children = append(nodeDictionary[parent].children, &node)
	return &node
}

// ID returns a node's id.
func (node *Node) ID() string {
	return node.id
}

// Children returns a node's children.
func (node *Node) Children() []*Node {
	return node.children
}

// Element returns the node's element.
func (node *Node) Element() Element {
	return node.element
}
