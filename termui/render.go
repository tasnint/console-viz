// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"image"
	"sync"

	tb "github.com/nsf/termbox-go"
)

// drawable interface that all widgets implement
type Drawable interface {
	GetRect() image.Rectangle   // return widget's rectangle
	SetRect(int, int, int, int) // set widget's rectangle
	Draw(*Buffer)               // draw widget into buffer
	sync.Locker                 // embeds locking methods for the widget
}

// the Render function takes in any number of Drawable items and loops to render them to terminal
func Render(items ...Drawable) { // takes in any number of Drawable (widget) arguments
	for _, item := range items {
		buf := NewBuffer(item.GetRect())       // create new buffer for the widget's rectangle
		item.Lock()                            // lock the widget for thread safety
		item.Draw(buf)                         // draw the widget into the buffer (in memory)
		item.Unlock()                          // unlock the widget once drawing completes
		for point, cell := range buf.CellMap { // loop through the buffer's cells and render them to the terminal
			if point.In(buf.Rectangle) { //check if the point is within the buffer's rectangle, to avoid rendering outside of the widget's area
				tb.SetCell( // draws character at (x,y) with specified style
					point.X, point.Y, // coordinates of the cell
					cell.Rune, // character to be displayed
					tb.Attribute(cell.Style.Fg+1)|tb.Attribute(cell.Style.Modifier), tb.Attribute(cell.Style.Bg+1), // style attributes by calling tb (termbox-go) SetCell function
				)
			}
		}
	}
	tb.Flush() // flush is a function from termbox-go that updates the terminal, must be called after setting cells to see the changes in the terminal
}
