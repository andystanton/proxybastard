package uilayout

// Frame is as Frame does.
type Frame struct {
	x         int
	xStr      string
	y         int
	yStr      string
	width     int
	widthStr  string
	height    int
	heightStr string
}

// CreateFrame creates a frame.
func CreateFrame(widthStr, heightStr, xStr, yStr string) *Frame {
	frame := Frame{
		xStr:      xStr,
		yStr:      yStr,
		widthStr:  widthStr,
		heightStr: heightStr,
	}
	return &frame
}

// Update a frame.
func (f *Frame) Update(width, height int, selected bool) {
	f.width = parseDimensionString(f.widthStr, f.width, f.height, width, height)
	f.height = parseDimensionString(f.heightStr, f.width, f.height, width, height)
	f.x = parseDimensionString(f.xStr, f.width, f.height, width, height)
	f.y = parseDimensionString(f.yStr, f.width, f.height, width, height)
}

// Draw a frame.
func (f *Frame) Draw() {
	drawBox(f.x, f.y, f.width, f.height)
}
