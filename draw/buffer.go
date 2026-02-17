package draw

import (
	"console-viz/styling"
	"image"

	rw "github.com/mattn/go-runewidth"
)

// Cell represents a viewable terminal cell
type Cell struct {
	Rune  rune           // the character to be displayed in the cell
	Style styling.Style // the style of the cell, including foreground color, background color and modifiers like bold, underline etc.
}

// CellClear is the eraser which clears area by filling with blank cells
var CellClear = Cell{
	Rune:  ' ',
	Style: styling.StyleClear,
}

// NewCell takes 1 to 2 arguments
// 1st argument = rune
// 2nd argument = optional style
func NewCell(rune rune, args ...interface{}) Cell {
	style := styling.StyleClear
	if len(args) == 1 {
		style = args[0].(styling.Style)
	}
	return Cell{
		Rune:  rune,
		Style: style,
	}
}

// Buffer represents a section of a terminal and is a renderable rectangle of cells.
// A buffer is a collection of cells
type Buffer struct {
	image.Rectangle                      // rectangle is from Go's image package, it defines a rectangular area in terms of its minimum and maximum points (Min and Max)
	CellMap         map[image.Point]Cell // the coordinates of the cell in the buffer, and the cell itself
}

// creates a new buffer, takes boundaries as input
// example usage: NewBuffer(image.Rect(0, 0, 10, 10)) creates a buffer that covers the area from (0,0) to (10,10)
func NewBuffer(r image.Rectangle) *Buffer {
	buf := &Buffer{
		Rectangle: r,                          // sets buffer's position and size
		CellMap:   make(map[image.Point]Cell), // creates map to store the cells
	}
	buf.Fill(CellClear, r) // clears out specified area
	return buf             // returns pointer to the new buffer
}

// retrieves the cell at a specific coordinate in the buffer, takes a point as input and returns the cell at that point
func (self *Buffer) GetCell(p image.Point) Cell {
	return self.CellMap[p]
}

// sets the cell at a specific coordinate in the buffer, takes a cell and a point as input and sets the cell at that point
func (self *Buffer) SetCell(c Cell, p image.Point) {
	self.CellMap[p] = c // stores the cell in the buffer's map at the specified point
}

// places the same cell at every position inside a specified rectangle, kind of like a paint bucket tool
func (self *Buffer) Fill(c Cell, rect image.Rectangle) {
	// outer loop goes through each column in the rectangle (left to right)
	for x := rect.Min.X; x < rect.Max.X; x++ {
		// inner loop goes through each row in the rectangle (top to bottom)
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			self.SetCell(c, image.Pt(x, y))
		}
	}
}

// writes a string horizontally starting at the given point
// takes a string, a style and a point as input
func (self *Buffer) SetString(s string, style styling.Style, p image.Point) {
	runes := []rune(s) // convert string into slice of runes, e.g. Hello becomes ['H', 'e', 'l', 'l', 'o']
	x := 0
	for _, char := range runes { // loops through each rune (character) in the string
		self.SetCell(Cell{char, style}, image.Pt(p.X+x, p.Y)) // sets the cell at the current position to the character
		x += rw.RuneWidth(char)                               // moves the x position by the width of the character
	}
}
