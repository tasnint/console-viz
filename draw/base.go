// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package draw

import (
	"console-viz/styling"
	"image" // Go's standard library for rectangles and points
	"sync"  // for Mutex, thread safety when multiple goroutines access the same widget
)

// rename all instances of block to base
// Base is the base struct inherited by most widgets.

// Base is the base struct inherited by most widgets.
// Base manages size, position, border, and title.
// It implements all 3 of the methods needed for the `Drawable` interface.
// Custom widgets will override the Draw method.
// Base is the foundation for all widgets and containers.
type Base struct {
	Border       bool                 // should we draw a border around base?
	BorderType   styling.BorderType   // single, double, rounded, etc.
	BorderStyle  styling.Style        // default border style
	BorderStyles styling.BorderStyles // per-side border styles (optional)

	BorderLeft, BorderRight, BorderTop, BorderBottom bool // controls which sides of the border to draw

	PaddingLeft, PaddingRight, PaddingTop, PaddingBottom int // space between border and content (in cells)
	MarginLeft, MarginRight, MarginTop, MarginBottom     int // margin outside the border

	BackgroundStyle styling.Style // background fill style

	Focused    bool          // is this base focused?
	FocusStyle styling.Style // style to use when focused

	image.Rectangle                 // embedded structure, means base inherits all fields from rectangle
	Inner           image.Rectangle // area inside the border and padding

	Title      string        // text (title) displayed at the top of the border
	TitleStyle styling.Style // style for the title text

	// Mouse event hooks
	OnClick func(x, y int)
	OnHover func(x, y int)

	sync.Mutex // embedded mutex to allow base.Lock() and base.Unlock() for when multiple goroutines access the same base
}

// factory fuction that creates a Base with default settings
// NewBase creates a Base with sensible defaults and theme fallbacks.
func NewBase() *Base {
	theme := styling.GetTheme()
	defaultBorderStyle := theme.Block.Border
	if defaultBorderStyle == (styling.Style{}) {
		defaultBorderStyle = styling.Style{Fg: styling.ColorWhite, Bg: styling.ColorBlack}
	}
	defaultTitleStyle := theme.Block.Title
	if defaultTitleStyle == (styling.Style{}) {
		defaultTitleStyle = styling.Style{Fg: styling.ColorWhite, Bg: styling.ColorBlack, Modifier: styling.ModifierBold}
	}
	// Note: Background and Focus styles are not in BlockTheme, using defaults
	defaultBackgroundStyle := styling.Style{Bg: styling.ColorBlack}
	defaultFocusStyle := styling.Style{Fg: styling.ColorYellow, Modifier: styling.ModifierBold}
	return &Base{
		Border:       true,
		BorderType:   styling.BorderSingle,
		BorderStyle:  defaultBorderStyle,
		BorderStyles: styling.BorderStyles{},
		BorderLeft:   true,
		BorderRight:  true,
		BorderTop:    true,
		BorderBottom: true,

		PaddingLeft:   1,
		PaddingRight:  1,
		PaddingTop:    1,
		PaddingBottom: 1,
		MarginLeft:    0,
		MarginRight:   0,
		MarginTop:     0,
		MarginBottom:  0,

		BackgroundStyle: defaultBackgroundStyle,
		Focused:         false,
		FocusStyle:      defaultFocusStyle,
		TitleStyle:      defaultTitleStyle,
	}
}

// drawBorder draws the border of the base using the selected BorderType and per-side styles.
func (self *Base) drawBorder(buf *Buffer) {
	// Choose border symbols based on BorderType
	var vert, horiz, tl, tr, bl, br rune
	switch self.BorderType {
	case styling.BorderDouble:
		vert, horiz, tl, tr, bl, br = styling.DOUBLE_VERTICAL, styling.DOUBLE_HORIZONTAL, styling.DOUBLE_TOP_LEFT, styling.DOUBLE_TOP_RIGHT, styling.DOUBLE_BOTTOM_LEFT, styling.DOUBLE_BOTTOM_RIGHT
	case styling.BorderRounded:
		vert, horiz, tl, tr, bl, br = styling.ROUNDED_VERTICAL, styling.ROUNDED_HORIZONTAL, styling.ROUNDED_TOP_LEFT, styling.ROUNDED_TOP_RIGHT, styling.ROUNDED_BOTTOM_LEFT, styling.ROUNDED_BOTTOM_RIGHT
	default:
		vert, horiz, tl, tr, bl, br = styling.VERTICAL_LINE, styling.HORIZONTAL_LINE, styling.TOP_LEFT, styling.TOP_RIGHT, styling.BOTTOM_LEFT, styling.BOTTOM_RIGHT
	}
	// Use per-side styles if set, else fallback to BorderStyle
	leftStyle := self.BorderStyle
	if self.BorderStyles.Left != (styling.Style{}) {
		leftStyle = self.BorderStyles.Left
	}
	rightStyle := self.BorderStyle
	if self.BorderStyles.Right != (styling.Style{}) {
		rightStyle = self.BorderStyles.Right
	}
	topStyle := self.BorderStyle
	if self.BorderStyles.Top != (styling.Style{}) {
		topStyle = self.BorderStyles.Top
	}
	bottomStyle := self.BorderStyle
	if self.BorderStyles.Bottom != (styling.Style{}) {
		bottomStyle = self.BorderStyles.Bottom
	}
	// Draw lines
	if self.BorderTop {
		buf.Fill(Cell{horiz, topStyle}, image.Rect(self.Min.X, self.Min.Y, self.Max.X, self.Min.Y+1))
	}
	if self.BorderBottom {
		buf.Fill(Cell{horiz, bottomStyle}, image.Rect(self.Min.X, self.Max.Y-1, self.Max.X, self.Max.Y))
	}
	if self.BorderLeft {
		buf.Fill(Cell{vert, leftStyle}, image.Rect(self.Min.X, self.Min.Y, self.Min.X+1, self.Max.Y))
	}
	if self.BorderRight {
		buf.Fill(Cell{vert, rightStyle}, image.Rect(self.Max.X-1, self.Min.Y, self.Max.X, self.Max.Y))
	}
	// Draw corners
	if self.BorderTop && self.BorderLeft {
		buf.SetCell(Cell{tl, topStyle}, self.Min)
	}
	if self.BorderTop && self.BorderRight {
		buf.SetCell(Cell{tr, topStyle}, image.Pt(self.Max.X-1, self.Min.Y))
	}
	if self.BorderBottom && self.BorderLeft {
		buf.SetCell(Cell{bl, bottomStyle}, image.Pt(self.Min.X, self.Max.Y-1))
	}
	if self.BorderBottom && self.BorderRight {
		buf.SetCell(Cell{br, bottomStyle}, self.Max.Sub(image.Pt(1, 1)))
	}
}

// Draw implements the Drawable interface. It draws background, border, title, and focus indicator.
func (self *Base) Draw(buf *Buffer) {
	// Fill background
	if self.BackgroundStyle != (styling.Style{}) {
		buf.Fill(Cell{' ', self.BackgroundStyle}, self.Rectangle)
	}
	// Draw border
	if self.Border {
		self.drawBorder(buf)
	}
	// Draw title
	buf.SetString(
		self.Title,
		self.TitleStyle,
		image.Pt(self.Min.X+2, self.Min.Y),
	)
	// Draw focus indicator (e.g., highlight border)
	if self.Focused && self.FocusStyle != (styling.Style{}) {
		// Example: overlay border with focus style (could be improved)
		if self.Border {
			buf.Fill(Cell{' ', self.FocusStyle}, image.Rect(self.Min.X, self.Min.Y, self.Max.X, self.Min.Y+1))
		}
	}
}

// SetRect sets the outer rectangle and calculates the inner content area based on padding and margin.
func (self *Base) SetRect(x1, y1, x2, y2 int) {
	self.Rectangle = image.Rect(x1+self.MarginLeft, y1+self.MarginTop, x2-self.MarginRight, y2-self.MarginBottom)
	self.Inner = image.Rect(
		self.Min.X+1+self.PaddingLeft,
		self.Min.Y+1+self.PaddingTop,
		self.Max.X-1-self.PaddingRight,
		self.Max.Y-1-self.PaddingBottom,
	)
}

// GetRect returns the coordinates of the base's outer rectangle (min and max points)
func (self *Base) GetRect() image.Rectangle {
	return self.Rectangle
}

// SetPadding sets all four paddings at once.
func (self *Base) SetPadding(left, right, top, bottom int) {
	self.PaddingLeft = left
	self.PaddingRight = right
	self.PaddingTop = top
	self.PaddingBottom = bottom
}

// SetMargin sets all four margins at once.
func (self *Base) SetMargin(left, right, top, bottom int) {
	self.MarginLeft = left
	self.MarginRight = right
	self.MarginTop = top
	self.MarginBottom = bottom
}
