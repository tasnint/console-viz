package widgets

import (
	"console-viz/draw"
	"console-viz/styling"
	"console-viz/utils"
	"fmt"
	"image"

	rw "github.com/mattn/go-runewidth"
)

// HorizontalBarChart displays a horizontal bar chart (bars extend right)
// Perfect for category distributions like Gender Distribution
type HorizontalBarChart struct {
	draw.Base
	BarColors    []styling.Color              // Colors for bars (cycled)
	LabelStyles  []styling.Style              // Styles for labels (cycled)
	ValueStyles  []styling.Style              // Styles for values (cycled)
	ValueFormatter func(float64) string       // Function to format values
	Data         []float64                    // Data values
	Labels       []string                     // Labels for each bar (displayed on left)
	MaxVal       float64                      // Maximum value (0 = auto-calculate)
	LabelWidth   int                          // Width reserved for labels (0 = auto-calculate)
	ValueWidth   int                          // Width reserved for values (0 = auto-calculate)
	BarChar      rune                         // Character to use for bar ('█', '▉', etc.)
	Gap          int                          // Gap between label, bar, and value
}

// NewHorizontalBarChart creates a new HorizontalBarChart widget
func NewHorizontalBarChart() *HorizontalBarChart {
	theme := styling.GetTheme()
	return &HorizontalBarChart{
		Base:        *draw.NewBase(),
		BarColors:   theme.BarChart.Bars,
		ValueStyles: theme.BarChart.Nums,
		LabelStyles: theme.BarChart.Labels,
		ValueFormatter: func(n float64) string {
			return fmt.Sprintf("%.0f", n)
		},
		BarChar:    styling.SHADED_BLOCKS[4], // Full block character '█'
		Gap:        2,   // Gap between sections
		LabelWidth: 0,   // Auto-calculate
		ValueWidth: 0,   // Auto-calculate
	}
}

// Draw renders the horizontal bar chart widget
func (hbc *HorizontalBarChart) Draw(buf *draw.Buffer) {
	hbc.Base.Draw(buf)

	if len(hbc.Data) == 0 {
		return
	}

	// Calculate max value
	maxVal := hbc.MaxVal
	if maxVal == 0 {
		var err error
		maxVal, err = utils.GetMaxFloat64FromSlice(hbc.Data)
		if err != nil || maxVal == 0 {
			return
		}
	}

	// Calculate label width (max label length + padding)
	labelWidth := hbc.LabelWidth
	if labelWidth == 0 {
		maxLabelLen := 0
		for _, label := range hbc.Labels {
			len := rw.StringWidth(label)
			if len > maxLabelLen {
				maxLabelLen = len
			}
		}
		labelWidth = maxLabelLen + 2 // Add padding
	}

	// Calculate value width (max value string length + padding)
	valueWidth := hbc.ValueWidth
	if valueWidth == 0 {
		maxValueLen := 0
		for _, data := range hbc.Data {
			valueStr := hbc.ValueFormatter(data)
			len := rw.StringWidth(valueStr)
			if len > maxValueLen {
				maxValueLen = len
			}
		}
		valueWidth = maxValueLen + 2 // Add padding
	}

	// Calculate available width for bars
	// Total width - labels - values - gaps
	availableBarWidth := hbc.Inner.Dx() - labelWidth - valueWidth - (hbc.Gap * 2)
	if availableBarWidth < 1 {
		availableBarWidth = 1
	}

	// Draw each row
	for i, data := range hbc.Data {
		y := hbc.Inner.Min.Y + i
		if y >= hbc.Inner.Max.Y {
			break // Out of bounds
		}

		// 1. Draw label on LEFT
		labelX := hbc.Inner.Min.X
		if i < len(hbc.Labels) {
			labelStyle := utils.SelectStyle(hbc.LabelStyles, i)
			buf.SetString(hbc.Labels[i], labelStyle, image.Pt(labelX, y))
		}

		// 2. Calculate bar width (proportional to value)
		barWidth := int((data / maxVal) * float64(availableBarWidth))
		if barWidth < 0 {
			barWidth = 0
		}

		// 3. Draw bar extending RIGHT from label
		barStartX := hbc.Inner.Min.X + labelWidth + hbc.Gap
		barColor := utils.SelectColor(hbc.BarColors, i)
		barCell := draw.NewCell(hbc.BarChar, styling.NewStyle(styling.ColorClear, barColor))

		for x := barStartX; x < barStartX+barWidth && x < hbc.Inner.Max.X-valueWidth-hbc.Gap; x++ {
			buf.SetCell(barCell, image.Pt(x, y))
		}

		// 4. Draw value on RIGHT
		valueX := hbc.Inner.Min.X + labelWidth + hbc.Gap + availableBarWidth + hbc.Gap
		if valueX < hbc.Inner.Max.X {
			valueStr := hbc.ValueFormatter(data)
			valueStyle := utils.SelectStyle(hbc.ValueStyles, i)
			buf.SetString(valueStr, valueStyle, image.Pt(valueX, y))
		}
	}
}
