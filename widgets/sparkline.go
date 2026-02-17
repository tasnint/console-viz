package widgets

import (
	"console-viz/draw"
	"console-viz/styling"
	"console-viz/utils"
	"image"
)

// Sparkline displays a mini line chart (single data series)
type Sparkline struct {
	Data       []float64      // Data points (should be non-negative)
	Title      string          // Title text
	TitleStyle styling.Style   // Style for title
	LineColor  styling.Color   // Color of the sparkline
	MaxVal     float64         // Maximum value (0 = auto-calculate)
	MaxHeight  int             // Maximum height of bars
}

// SparklineGroup groups multiple sparklines together
type SparklineGroup struct {
	draw.Base
	Sparklines []*Sparkline // List of sparklines to display
}

// NewSparkline creates a new Sparkline (must be added to a SparklineGroup)
func NewSparkline() *Sparkline {
	theme := styling.GetTheme()
	return &Sparkline{
		TitleStyle: theme.Sparkline.Title,
		LineColor:  theme.Sparkline.Line,
	}
}

// NewSparklineGroup creates a new SparklineGroup widget
func NewSparklineGroup(sparklines ...*Sparkline) *SparklineGroup {
	return &SparklineGroup{
		Base:       *draw.NewBase(),
		Sparklines: sparklines,
	}
}

// Draw renders the sparkline group widget
func (sg *SparklineGroup) Draw(buf *draw.Buffer) {
	sg.Base.Draw(buf)

	if len(sg.Sparklines) == 0 {
		return
	}

	sparklineHeight := sg.Inner.Dy() / len(sg.Sparklines)

	for i, sl := range sg.Sparklines {
		heightOffset := sparklineHeight * (i + 1)
		barHeight := sparklineHeight
		if i == len(sg.Sparklines)-1 {
			heightOffset = sg.Inner.Dy()
			barHeight = sg.Inner.Dy() - (sparklineHeight * i)
		}
		if sl.Title != "" {
			barHeight--
		}

		// Calculate max value
		maxVal := sl.MaxVal
		if maxVal == 0 {
			var err error
			maxVal, err = utils.GetMaxFloat64FromSlice(sl.Data)
			if err != nil || maxVal == 0 {
				continue
			}
		}

		// Draw sparkline bars
		for j := 0; j < len(sl.Data) && j < sg.Inner.Dx(); j++ {
			data := sl.Data[j]
			height := int((data / maxVal) * float64(barHeight))
			if height > sl.MaxHeight && sl.MaxHeight > 0 {
				height = sl.MaxHeight
			}

			sparkChar := styling.BARS[len(styling.BARS)-1] // Full block
			for k := 0; k < height; k++ {
				buf.SetCell(
					draw.NewCell(sparkChar, styling.NewStyle(sl.LineColor)),
					image.Pt(j+sg.Inner.Min.X, sg.Inner.Min.Y-1+heightOffset-k),
				)
			}
			if height == 0 {
				sparkChar = styling.BARS[1] // Small bar
				buf.SetCell(
					draw.NewCell(sparkChar, styling.NewStyle(sl.LineColor)),
					image.Pt(j+sg.Inner.Min.X, sg.Inner.Min.Y-1+heightOffset),
				)
			}
		}

		// Draw title
		if sl.Title != "" {
			title := utils.TrimString(sl.Title, sg.Inner.Dx())
			buf.SetString(
				title,
				sl.TitleStyle,
				image.Pt(sg.Inner.Min.X, sg.Inner.Min.Y-1+heightOffset-barHeight),
			)
		}
	}
}
