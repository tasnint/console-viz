package widgets

import (
	"console-viz/draw"
	"console-viz/styling"
	"fmt"
	"image"
)

// Gauge displays a progress bar with percentage
type Gauge struct {
	draw.Base
	Percent    int             // Progress percentage (0-100)
	BarColor   styling.Color   // Color of the progress bar
	Label      string          // Custom label (if empty, shows percentage)
	LabelStyle styling.Style   // Style for the label text
}

// NewGauge creates a new Gauge widget with default settings
func NewGauge() *Gauge {
	theme := styling.GetTheme()
	return &Gauge{
		Base:       *draw.NewBase(),
		BarColor:   theme.Gauge.Bar,
		LabelStyle: theme.Gauge.Label,
	}
}

// Draw renders the gauge widget
func (g *Gauge) Draw(buf *draw.Buffer) {
	g.Base.Draw(buf)

	// Determine label text
	label := g.Label
	if label == "" {
		label = fmt.Sprintf("%d%%", g.Percent)
	}

	// Calculate bar width
	barWidth := int((float64(g.Percent) / 100.0) * float64(g.Inner.Dx()))

	// Draw progress bar
	buf.Fill(
		draw.NewCell(' ', styling.NewStyle(styling.ColorClear, g.BarColor)),
		image.Rect(g.Inner.Min.X, g.Inner.Min.Y, g.Inner.Min.X+barWidth, g.Inner.Max.Y),
	)

	// Draw label centered
	labelX := g.Inner.Min.X + (g.Inner.Dx() / 2) - (len(label) / 2)
	labelY := g.Inner.Min.Y + ((g.Inner.Dy() - 1) / 2)

	if labelY < g.Inner.Max.Y {
		for i, char := range label {
			style := g.LabelStyle
			// If label is over the bar, use reverse style
			if labelX+i+1 <= g.Inner.Min.X+barWidth {
				style = styling.NewStyle(g.BarColor, styling.ColorClear, styling.ModifierReverse)
			}
			buf.SetCell(draw.NewCell(char, style), image.Pt(labelX+i, labelY))
		}
	}
}
