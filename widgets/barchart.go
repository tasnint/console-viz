package widgets

import (
	"console-viz/draw"
	"console-viz/styling"
	"console-viz/utils"
	"fmt"
	"image"

	rw "github.com/mattn/go-runewidth"
)

// BarChart displays a vertical bar chart
type BarChart struct {
	draw.Base
	BarColors    []styling.Color              // Colors for bars (cycled)
	LabelStyles  []styling.Style              // Styles for labels (cycled)
	NumStyles    []styling.Style              // Styles for numbers (cycled)
	NumFormatter func(float64) string         // Function to format numbers
	Data         []float64                    // Data values
	Labels       []string                     // Labels for each bar
	BarWidth     int                          // Width of each bar
	BarGap       int                          // Gap between bars
	MaxVal       float64                      // Maximum value (0 = auto-calculate)
}

// NewBarChart creates a new BarChart widget with default settings
func NewBarChart() *BarChart {
	theme := styling.GetTheme()
	return &BarChart{
		Base:        *draw.NewBase(),
		BarColors:   theme.BarChart.Bars,
		NumStyles:   theme.BarChart.Nums,
		LabelStyles: theme.BarChart.Labels,
		NumFormatter: func(n float64) string {
			return fmt.Sprint(n)
		},
		BarGap:   1,
		BarWidth: 3,
	}
}

// Draw renders the bar chart widget
func (bc *BarChart) Draw(buf *draw.Buffer) {
	bc.Base.Draw(buf)

	// Calculate max value
	maxVal := bc.MaxVal
	if maxVal == 0 {
		var err error
		maxVal, err = utils.GetMaxFloat64FromSlice(bc.Data)
		if err != nil || maxVal == 0 {
			return // No data to display
		}
	}

	barX := bc.Inner.Min.X

	// Draw each bar
	for i, data := range bc.Data {
		if data > 0 {
			// Calculate bar height
			height := int((data / maxVal) * float64(bc.Inner.Dy()-1))

			// Draw bar
			for x := barX; x < utils.MinInt(barX+bc.BarWidth, bc.Inner.Max.X); x++ {
				for y := bc.Inner.Max.Y - 2; y > (bc.Inner.Max.Y-2)-height; y-- {
					barColor := utils.SelectColor(bc.BarColors, i)
					c := draw.NewCell(' ', styling.NewStyle(styling.ColorClear, barColor))
					buf.SetCell(c, image.Pt(x, y))
				}
			}
		}

		// Draw label
		if i < len(bc.Labels) {
			labelX := barX + (bc.BarWidth / 2) - (rw.StringWidth(bc.Labels[i]) / 2)
			labelStyle := utils.SelectStyle(bc.LabelStyles, i)
			buf.SetString(bc.Labels[i], labelStyle, image.Pt(labelX, bc.Inner.Max.Y-1))
		}

		// Draw number
		numberX := barX + (bc.BarWidth / 2)
		if numberX <= bc.Inner.Max.X {
			numStr := bc.NumFormatter(data)
			numStyle := utils.SelectStyle(bc.NumStyles, i+1)
			barColor := utils.SelectColor(bc.BarColors, i)
			style := styling.NewStyle(numStyle.Fg, barColor, numStyle.Modifier)
			buf.SetString(numStr, style, image.Pt(numberX, bc.Inner.Max.Y-2))
		}

		barX += (bc.BarWidth + bc.BarGap)
	}
}
