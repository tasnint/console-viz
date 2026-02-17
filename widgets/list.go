package widgets

import (
	"console-viz/draw"
	"console-viz/styling"
	"console-viz/utils"
	"image"

	rw "github.com/mattn/go-runewidth"
)

// List displays a scrollable list of items with selection highlighting
type List struct {
	draw.Base
	Rows             []string        // List items (can contain style markup)
	WrapText         bool            // Whether to wrap text
	TextStyle        styling.Style   // Default text style
	SelectedRow      int             // Currently selected row index
	topRow           int             // Top visible row (for scrolling)
	SelectedRowStyle styling.Style   // Style for selected row
}

// NewList creates a new List widget with default settings
func NewList() *List {
	theme := styling.GetTheme()
	return &List{
		Base:            *draw.NewBase(),
		TextStyle:       theme.List.Text,
		SelectedRowStyle: theme.List.Text,
		WrapText:        true,
	}
}

// Draw renders the list widget
func (l *List) Draw(buf *draw.Buffer) {
	l.Base.Draw(buf)

	point := l.Inner.Min

	// Adjust view to show selected row
	if l.SelectedRow >= l.Inner.Dy()+l.topRow {
		l.topRow = l.SelectedRow - l.Inner.Dy() + 1
	} else if l.SelectedRow < l.topRow {
		l.topRow = l.SelectedRow
	}

	// Draw visible rows
	for row := l.topRow; row < len(l.Rows) && point.Y < l.Inner.Max.Y; row++ {
		// Parse styles from row text
		cells := draw.ParseStyles(l.Rows[row], l.TextStyle)

		// Wrap if enabled
		if l.WrapText {
			cells = utils.WrapCells(cells, uint(l.Inner.Dx()))
		}

		// Draw cells
		for j := 0; j < len(cells) && point.Y < l.Inner.Max.Y; j++ {
			style := cells[j].Style
			if row == l.SelectedRow {
				style = l.SelectedRowStyle
			}

			if cells[j].Rune == '\n' {
				point = image.Pt(l.Inner.Min.X, point.Y+1)
			} else {
				// Check if we need ellipsis
				if point.X+1 == l.Inner.Max.X+1 && len(cells) > l.Inner.Dx() {
					buf.SetCell(draw.NewCell(styling.ELLIPSES, style), point.Add(image.Pt(-1, 0)))
					break
				} else {
					buf.SetCell(draw.NewCell(cells[j].Rune, style), point)
					point = point.Add(image.Pt(rw.RuneWidth(cells[j].Rune), 0))
				}
			}
		}
		point = image.Pt(l.Inner.Min.X, point.Y+1)
	}

	// Draw scroll indicators
	if l.topRow > 0 {
		buf.SetCell(
			draw.NewCell(styling.UP_ARROW, styling.NewStyle(styling.ColorWhite)),
			image.Pt(l.Inner.Max.X-1, l.Inner.Min.Y),
		)
	}

	if len(l.Rows) > int(l.topRow)+l.Inner.Dy() {
		buf.SetCell(
			draw.NewCell(styling.DOWN_ARROW, styling.NewStyle(styling.ColorWhite)),
			image.Pt(l.Inner.Max.X-1, l.Inner.Max.Y-1),
		)
	}
}

// ScrollAmount scrolls by the given amount (negative = up, positive = down)
func (l *List) ScrollAmount(amount int) {
	if len(l.Rows)-int(l.SelectedRow) <= amount {
		l.SelectedRow = len(l.Rows) - 1
	} else if int(l.SelectedRow)+amount < 0 {
		l.SelectedRow = 0
	} else {
		l.SelectedRow += amount
	}
}

// ScrollUp scrolls up one row
func (l *List) ScrollUp() {
	l.ScrollAmount(-1)
}

// ScrollDown scrolls down one row
func (l *List) ScrollDown() {
	l.ScrollAmount(1)
}

// ScrollPageUp scrolls up one page
func (l *List) ScrollPageUp() {
	if l.SelectedRow > l.topRow {
		l.SelectedRow = l.topRow
	} else {
		l.ScrollAmount(-l.Inner.Dy())
	}
}

// ScrollPageDown scrolls down one page
func (l *List) ScrollPageDown() {
	l.ScrollAmount(l.Inner.Dy())
}

// ScrollTop scrolls to the top
func (l *List) ScrollTop() {
	l.SelectedRow = 0
}

// ScrollBottom scrolls to the bottom
func (l *List) ScrollBottom() {
	l.SelectedRow = len(l.Rows) - 1
}
