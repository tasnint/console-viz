package widgets

import (
	"console-viz/draw"
	"console-viz/styling"
	"console-viz/utils"
	"image"
)

// Table displays tabular data with columns and rows
type Table struct {
	draw.Base
	Rows          [][]string              // Table rows (each row is []string of cells)
	ColumnWidths  []int                   // Width of each column (0 = auto-calculate)
	TextStyle     styling.Style           // Default text style
	RowSeparator  bool                    // Whether to show separators between rows
	TextAlignment draw.Alignment          // Text alignment (Left, Center, Right)
	RowStyles     map[int]styling.Style    // Per-row styles
	FillRow       bool                    // Whether to fill entire row width

	// ColumnResizer is called on each Draw for custom column sizing
	ColumnResizer func()
}

// NewTable creates a new Table widget with default settings
func NewTable() *Table {
	theme := styling.GetTheme()
	return &Table{
		Base:          *draw.NewBase(),
		TextStyle:     theme.Table.Text,
		RowSeparator:  true,
		RowStyles:     make(map[int]styling.Style),
		ColumnResizer: func() {},
	}
}

// Draw renders the table widget
func (t *Table) Draw(buf *draw.Buffer) {
	t.Base.Draw(buf)

	t.ColumnResizer()

	// Calculate column widths if not set
	columnWidths := t.ColumnWidths
	if len(columnWidths) == 0 && len(t.Rows) > 0 {
		columnCount := len(t.Rows[0])
		columnWidth := t.Inner.Dx() / columnCount
		for i := 0; i < columnCount; i++ {
			columnWidths = append(columnWidths, columnWidth)
		}
	}

	y := t.Inner.Min.Y

	// Draw rows
	for i := 0; i < len(t.Rows) && y < t.Inner.Max.Y; i++ {
		row := t.Rows[i]
		colX := t.Inner.Min.X

		// Get row style
		rowStyle := t.TextStyle
		if style, ok := t.RowStyles[i]; ok {
			rowStyle = style
		}

		// Fill row if enabled
		if t.FillRow {
			blankCell := draw.NewCell(' ', rowStyle)
			buf.Fill(blankCell, image.Rect(t.Inner.Min.X, y, t.Inner.Max.X, y+1))
		}

		// Draw row cells
		for j := 0; j < len(row) && j < len(columnWidths); j++ {
			// Parse styles from cell text
			cells := draw.ParseStyles(row[j], rowStyle)

			// Draw cell based on alignment
			if len(cells) > columnWidths[j] || t.TextAlignment == draw.AlignLeft {
				cellArray := utils.BuildCellWithXArray(cells)
				for _, cx := range cellArray {
					k, cell := cx.X, cx.Cell
					if k == columnWidths[j] || colX+k == t.Inner.Max.X {
						cell.Rune = styling.ELLIPSES
						buf.SetCell(cell, image.Pt(colX+k-1, y))
						break
					} else {
						buf.SetCell(cell, image.Pt(colX+k, y))
					}
				}
			} else if t.TextAlignment == draw.AlignCenter {
				xOffset := (columnWidths[j] - len(cells)) / 2
				stringX := xOffset + colX
				cellArray := utils.BuildCellWithXArray(cells)
				for _, cx := range cellArray {
					k, cell := cx.X, cx.Cell
					buf.SetCell(cell, image.Pt(stringX+k, y))
				}
			} else if t.TextAlignment == draw.AlignRight {
				stringX := utils.MinInt(colX+columnWidths[j], t.Inner.Max.X) - len(cells)
				cellArray := utils.BuildCellWithXArray(cells)
				for _, cx := range cellArray {
					k, cell := cx.X, cx.Cell
					buf.SetCell(cell, image.Pt(stringX+k, y))
				}
			}
			colX += columnWidths[j] + 1
		}

		// Draw vertical separators
		separatorStyle := t.BorderStyle
		separatorX := t.Inner.Min.X
		verticalCell := draw.NewCell(styling.VERTICAL_LINE, separatorStyle)

		for i, width := range columnWidths {
			if t.FillRow && i < len(columnWidths)-1 {
				verticalCell.Style.Bg = rowStyle.Bg
			} else {
				verticalCell.Style.Bg = t.BorderStyle.Bg
			}
			separatorX += width
			buf.SetCell(verticalCell, image.Pt(separatorX, y))
			separatorX++
		}

		y++

		// Draw horizontal separator
		if t.RowSeparator && y < t.Inner.Max.Y && i != len(t.Rows)-1 {
			horizontalCell := draw.NewCell(styling.HORIZONTAL_LINE, separatorStyle)
			buf.Fill(horizontalCell, image.Rect(t.Inner.Min.X, y, t.Inner.Max.X, y+1))
			y++
		}
	}
}
