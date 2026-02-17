package widgets

import (
	"console-viz/draw"
	"console-viz/styling"
	"console-viz/utils"
	"fmt"
	"image"

	rw "github.com/mattn/go-runewidth"
)

// StackedBarChart displays a bar chart with stacked segments
type StackedBarChart struct {
	draw.Base
	BarColors    []styling.Color              // Colors for bars (cycled)
	LabelStyles  []styling.Style              // Styles for labels (cycled)
	NumStyles    []styling.Style              // Styles for numbers (cycled)
	NumFormatter func(float64) string         // Function to format numbers
	Data         [][]float64                  // Data values (each []float64 is one bar with multiple segments)
	Labels       []string                     // Labels for each bar
	BarWidth     int                          // Width of each bar
	BarGap       int                          // Gap between bars
	MaxVal       float64                      // Maximum value (0 = auto-calculate)
}

// NewStackedBarChart creates a new StackedBarChart widget with default settings
func NewStackedBarChart() *StackedBarChart {
	theme := styling.GetTheme()
	return &StackedBarChart{
		Base:        *draw.NewBase(),
		BarColors:   theme.StackedBarChart.Bars,
		NumStyles:   theme.StackedBarChart.Nums,
		LabelStyles: theme.StackedBarChart.Labels,
		NumFormatter: func(n float64) string {
			return fmt.Sprint(n)
		},
		BarGap:   1,
		BarWidth: 3,
	}
}

// Draw renders the stacked bar chart widget
func (sbc *StackedBarChart) Draw(buf *draw.Buffer) {
	sbc.Base.Draw(buf)

	// Calculate max value
	maxVal := sbc.MaxVal
	if maxVal == 0 {
		for _, bar := range sbc.Data {
			sum := utils.SumFloat64Slice(bar)
			maxVal = utils.MaxFloat64(maxVal, sum)
		}
	}

	if maxVal == 0 {
		return // No data to display
	}

	barX := sbc.Inner.Min.X

	// Draw each bar
	for i, bar := range sbc.Data {
		stackedY := 0

		// Draw each segment of the stacked bar
		for j, data := range bar {
			if data > 0 {
				// Calculate segment height
				height := int((data / maxVal) * float64(sbc.Inner.Dy()-1))

				// Draw segment
				for x := barX; x < utils.MinInt(barX+sbc.BarWidth, sbc.Inner.Max.X); x++ {
					for y := (sbc.Inner.Max.Y - 2) - stackedY; y > (sbc.Inner.Max.Y-2)-stackedY-height; y-- {
						barColor := utils.SelectColor(sbc.BarColors, j)
						c := draw.NewCell(' ', styling.NewStyle(styling.ColorClear, barColor))
						buf.SetCell(c, image.Pt(x, y))
					}
				}

				// Draw number on segment
				numberX := barX + (sbc.BarWidth / 2) - 1
				numStr := sbc.NumFormatter(data)
				numStyle := utils.SelectStyle(sbc.NumStyles, j+1)
				barColor := utils.SelectColor(sbc.BarColors, j)
				style := styling.NewStyle(numStyle.Fg, barColor, numStyle.Modifier)
				buf.SetString(numStr, style, image.Pt(numberX, (sbc.Inner.Max.Y-2)-stackedY))

				stackedY += height
			}
		}

		// Draw label
		if i < len(sbc.Labels) {
			labelX := barX + utils.MaxInt(
				(sbc.BarWidth/2)-(rw.StringWidth(sbc.Labels[i])/2),
				0,
			)
			label := utils.TrimString(sbc.Labels[i], sbc.BarWidth)
			labelStyle := utils.SelectStyle(sbc.LabelStyles, i)
			buf.SetString(label, labelStyle, image.Pt(labelX, sbc.Inner.Max.Y-1))
		}

		barX += (sbc.BarWidth + sbc.BarGap)
	}
}
