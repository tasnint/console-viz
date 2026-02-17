package widgets

import (
	"console-viz/draw"
	"console-viz/styling"
	"console-viz/utils"
	"image"
)

// Paragraph displays text with optional wrapping and style parsing
// Supports embedded style markup: [text](fg:red,bg:blue,mod:bold)
type Paragraph struct {
	draw.Base
	Text      string        // Text content (can contain style markup)
	TextStyle styling.Style // Default style for unmarked text
	WrapText  bool          // Whether to wrap text to fit width
}

// NewParagraph creates a new Paragraph widget with default settings
func NewParagraph() *Paragraph {
	theme := styling.GetTheme()
	return &Paragraph{
		Base:      *draw.NewBase(),
		TextStyle: theme.Paragraph.Text,
		WrapText:  true,
	}
}

// Draw renders the paragraph widget
func (p *Paragraph) Draw(buf *draw.Buffer) {
	p.Base.Draw(buf)

	// Parse styles from text
	cells := draw.ParseStyles(p.Text, p.TextStyle)

	// Wrap cells if enabled
	if p.WrapText {
		cells = utils.WrapCells(cells, uint(p.Inner.Dx()))
	}

	// Split cells by newlines
	rows := utils.SplitCells(cells, '\n')

	// Render each row
	for y, row := range rows {
		if y+p.Inner.Min.Y >= p.Inner.Max.Y {
			break
		}

		// Trim row to fit width
		row = utils.TrimCells(row, p.Inner.Dx())

		// Build cell array with X positions
		cellArray := utils.BuildCellWithXArray(row)
		for _, cx := range cellArray {
			x, cell := cx.X, cx.Cell
			pos := image.Pt(x, y).Add(p.Inner.Min)
			if pos.In(p.Inner) {
				buf.SetCell(cell, pos)
			}
		}
	}
}
