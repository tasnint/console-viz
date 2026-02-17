// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"image" // Go/s standard library for rectangles and points
	"sync"  // for Mutex, threas safety when multiple goroutines access the same widget
)

// Block is the base struct inherited by most widgets.
// Block manages size, position, border, and title.
// It implements all 3 of the methods needed for the `Drawable` interface.
// Custom widgets will override the Draw method.
type Block struct {
	Border      bool  // should we draw a border around block?
	BorderStyle Style // color/modifier for the border

	BorderLeft, BorderRight, BorderTop, BorderBottom bool // controls which sides of the border to draw

	PaddingLeft, PaddingRight, PaddingTop, PaddingBottom int // space between border and content (in cells)

	image.Rectangle                 // embedded structure, means block inherits all fields from rectangle
	Inner           image.Rectangle // different rectangle that defines the inner (content) area of the block, which is the area inside the border and padding, where content will be drawn

	Title      string // text (title) displayed at the top of the border
	TitleStyle Style  // style for the title text

	sync.Mutex // embedded mutex to allow block.Lock() and block.Unlock() for when multiple goroutines access the same block
}

// factory fuction that creates a Block with default settings
func NewBlock() *Block { // * returns a pointer to a new Block struct
	return &Block{ // & creates a pointer to this new Block
		Border:       true,
		BorderStyle:  Theme.Block.Border,
		BorderLeft:   true,
		BorderRight:  true,
		BorderTop:    true,
		BorderBottom: true,

		TitleStyle: Theme.Block.Title,
	}
}

// private method to draw the border of the block, takes the Buffer as input to draw into
func (self *Block) drawBorder(buf *Buffer) {
	verticalCell := Cell{VERTICAL_LINE, self.BorderStyle}     // cell for vertical line
	horizontalCell := Cell{HORIZONTAL_LINE, self.BorderStyle} // cell for horizontal line

	// draw lines
	if self.BorderTop {
		buf.Fill(horizontalCell, image.Rect(self.Min.X, self.Min.Y, self.Max.X, self.Min.Y+1))
	}
	if self.BorderBottom {
		buf.Fill(horizontalCell, image.Rect(self.Min.X, self.Max.Y-1, self.Max.X, self.Max.Y))
	}
	if self.BorderLeft {
		buf.Fill(verticalCell, image.Rect(self.Min.X, self.Min.Y, self.Min.X+1, self.Max.Y))
	}
	if self.BorderRight {
		buf.Fill(verticalCell, image.Rect(self.Max.X-1, self.Min.Y, self.Max.X, self.Max.Y))
	}

	// draw corners
	if self.BorderTop && self.BorderLeft {
		buf.SetCell(Cell{TOP_LEFT, self.BorderStyle}, self.Min)
	}
	if self.BorderTop && self.BorderRight {
		buf.SetCell(Cell{TOP_RIGHT, self.BorderStyle}, image.Pt(self.Max.X-1, self.Min.Y))
	}
	if self.BorderBottom && self.BorderLeft {
		buf.SetCell(Cell{BOTTOM_LEFT, self.BorderStyle}, image.Pt(self.Min.X, self.Max.Y-1))
	}
	if self.BorderBottom && self.BorderRight {
		buf.SetCell(Cell{BOTTOM_RIGHT, self.BorderStyle}, self.Max.Sub(image.Pt(1, 1)))
	}
}

// Draw implements the Drawable interface.
func (self *Block) Draw(buf *Buffer) {
	if self.Border {
		self.drawBorder(buf)
	}
	buf.SetString(
		self.Title,
		self.TitleStyle,
		image.Pt(self.Min.X+2, self.Min.Y),
	)
}

// defines/calculates the area where the content will be drawn
// content written inside block with 1 character of padding on all sides
func (self *Block) SetRect(x1, y1, x2, y2 int) {
	self.Rectangle = image.Rect(x1, y1, x2, y2)
	self.Inner = image.Rect(
		self.Min.X+1+self.PaddingLeft,
		self.Min.Y+1+self.PaddingTop,
		self.Max.X-1-self.PaddingRight,
		self.Max.Y-1-self.PaddingBottom,
	)
}

// returns the coordinates of the block's outer rectangle (min and max points)
func (self *Block) GetRect() image.Rectangle {
	return self.Rectangle
}
