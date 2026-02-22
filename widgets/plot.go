package widgets

import (
	"console-viz/draw"
	"console-viz/styling"
	"console-viz/utils"
	"fmt"
	"image"
)

const (
	xAxisLabelsHeight = 1
	yAxisLabelsWidth  = 4
	xAxisLabelsGap    = 2
	yAxisLabelsGap    = 1
)

// PlotType defines the type of plot
type PlotType uint

const (
	LineChart PlotType = iota
	ScatterPlot
)

// PlotMarker defines the marker type for data points
type PlotMarker uint

const (
	MarkerBraille PlotMarker = iota
	MarkerDot
)

// DrawDirection defines drawing direction (not fully implemented)
type DrawDirection uint

const (
	DrawLeft DrawDirection = iota
	DrawRight
)

// Plot displays line charts or scatter plots
// Supports multiple data series with different colors
type Plot struct {
	draw.Base
	Data           [][]float64      // Data series (each []float64 is one series)
	DataLabels     []string          // Labels for each data series
	MaxVal         float64           // Maximum value (0 = auto-calculate)
	LineColors     []styling.Color   // Colors for lines (cycled)
	AxesColor      styling.Color      // Color for axes
	ShowAxes       bool               // Whether to show axes
	Marker         PlotMarker         // Marker type (Braille or Dot)
	DotMarkerRune  rune               // Rune to use for dot markers
	PlotType       PlotType           // LineChart or ScatterPlot
	HorizontalScale int               // Horizontal scaling factor
	DrawDirection   DrawDirection     // Drawing direction
}

// NewPlot creates a new Plot widget with default settings
func NewPlot() *Plot {
	theme := styling.GetTheme()
	return &Plot{
		Base:           *draw.NewBase(),
		LineColors:     theme.Plot.Lines,
		AxesColor:      theme.Plot.Axes,
		Marker:         MarkerDot, // Use dot by default (braille requires Canvas)
		DotMarkerRune:  styling.DOT,
		Data:           [][]float64{},
		HorizontalScale: 1,
		DrawDirection:   DrawRight,
		ShowAxes:        true,
		PlotType:        LineChart,
	}
}

// plotAxes draws the axes and labels
func (p *Plot) plotAxes(buf *draw.Buffer, maxVal float64) {
	// Draw origin cell
	buf.SetCell(
		draw.NewCell(styling.BOTTOM_LEFT, styling.NewStyle(styling.ColorWhite)),
		image.Pt(p.Inner.Min.X+yAxisLabelsWidth, p.Inner.Max.Y-xAxisLabelsHeight-1),
	)

	// Draw x axis line
	for i := yAxisLabelsWidth + 1; i < p.Inner.Dx(); i++ {
		buf.SetCell(
			draw.NewCell(styling.HORIZONTAL_LINE, styling.NewStyle(p.AxesColor)),
			image.Pt(i+p.Inner.Min.X, p.Inner.Max.Y-xAxisLabelsHeight-1),
		)
	}

	// Draw y axis line
	for i := 0; i < p.Inner.Dy()-xAxisLabelsHeight-1; i++ {
		buf.SetCell(
			draw.NewCell(styling.VERTICAL_LINE, styling.NewStyle(p.AxesColor)),
			image.Pt(p.Inner.Min.X+yAxisLabelsWidth, i+p.Inner.Min.Y),
		)
	}

	// Draw x axis labels
	buf.SetString("0", styling.NewStyle(p.AxesColor), image.Pt(p.Inner.Min.X+yAxisLabelsWidth, p.Inner.Max.Y-1))

	for x := p.Inner.Min.X + yAxisLabelsWidth + (xAxisLabelsGap)*p.HorizontalScale + 1; x < p.Inner.Max.X-1; {
		label := fmt.Sprintf("%d", (x-(p.Inner.Min.X+yAxisLabelsWidth)-1)/(p.HorizontalScale)+1)
		buf.SetString(label, styling.NewStyle(p.AxesColor), image.Pt(x, p.Inner.Max.Y-1))
		x += (len(label) + xAxisLabelsGap) * p.HorizontalScale
	}

	// Draw y axis labels
	verticalScale := maxVal / float64(p.Inner.Dy()-xAxisLabelsHeight-1)
	for i := 0; i*(yAxisLabelsGap+1) < p.Inner.Dy()-1; i++ {
		label := fmt.Sprintf("%.2f", float64(i)*verticalScale*(yAxisLabelsGap+1))
		buf.SetString(label, styling.NewStyle(p.AxesColor), image.Pt(p.Inner.Min.X, p.Inner.Max.Y-(i*(yAxisLabelsGap+1))-2))
	}
}

// renderDot renders the plot using dot markers
func (p *Plot) renderDot(buf *draw.Buffer, drawArea image.Rectangle, maxVal float64) {
	switch p.PlotType {
	case ScatterPlot:
		for i, line := range p.Data {
			for j, val := range line {
				height := int((val / maxVal) * float64(drawArea.Dy()-1))
				point := image.Pt(drawArea.Min.X+(j*p.HorizontalScale), drawArea.Max.Y-1-height)
				if point.In(drawArea) {
					color := utils.SelectColor(p.LineColors, i)
					buf.SetCell(
						draw.NewCell(p.DotMarkerRune, styling.NewStyle(color)),
						point,
					)
				}
			}
		}
	case LineChart:
		for i, line := range p.Data {
			for j := 0; j < len(line) && j*p.HorizontalScale < drawArea.Dx(); j++ {
				val := line[j]
				height := int((val / maxVal) * float64(drawArea.Dy()-1))
				color := utils.SelectColor(p.LineColors, i)
				point := image.Pt(drawArea.Min.X+(j*p.HorizontalScale), drawArea.Max.Y-1-height)
				if point.In(drawArea) {
					buf.SetCell(
						draw.NewCell(p.DotMarkerRune, styling.NewStyle(color)),
						point,
					)
				}
				// Draw line to next point
				if j < len(line)-1 {
					nextVal := line[j+1]
					nextHeight := int((nextVal / maxVal) * float64(drawArea.Dy()-1))
					nextPoint := image.Pt(drawArea.Min.X+((j+1)*p.HorizontalScale), drawArea.Max.Y-1-nextHeight)
					p.drawLine(buf, point, nextPoint, color)
				}
			}
		}
	}
}

// drawLine draws a line between two points using the same dot rune as the data points,
// so the graph looks like connected dots (thin line), not thick bars.
func (p *Plot) drawLine(buf *draw.Buffer, p1, p2 image.Point, color styling.Color) {
	dx := utils.AbsInt(p2.X - p1.X)
	dy := utils.AbsInt(p2.Y - p1.Y)
	dotRune := p.DotMarkerRune // same as data points so line looks like connected dots
	if dx > dy {
		// Mostly horizontal: iterate x, compute y
		startX := utils.MinInt(p1.X, p2.X)
		endX := utils.MaxInt(p1.X, p2.X)
		denom := p2.X - p1.X
		if denom == 0 {
			return
		}
		for x := startX; x <= endX; x++ {
			y := p1.Y + (p2.Y-p1.Y)*(x-p1.X)/denom
			pt := image.Pt(x, y)
			if pt.In(buf.Rectangle) {
				buf.SetCell(draw.NewCell(dotRune, styling.NewStyle(color)), pt)
			}
		}
	} else {
		// Mostly vertical: iterate y, compute x
		startY := utils.MinInt(p1.Y, p2.Y)
		endY := utils.MaxInt(p1.Y, p2.Y)
		denom := p2.Y - p1.Y
		if denom == 0 {
			return
		}
		for y := startY; y <= endY; y++ {
			x := p1.X + (p2.X-p1.X)*(y-p1.Y)/denom
			pt := image.Pt(x, y)
			if pt.In(buf.Rectangle) {
				buf.SetCell(draw.NewCell(dotRune, styling.NewStyle(color)), pt)
			}
		}
	}
}

// Draw renders the plot widget
func (p *Plot) Draw(buf *draw.Buffer) {
	p.Base.Draw(buf)

	// Calculate max value
	maxVal := p.MaxVal
	if maxVal == 0 {
		var err error
		maxVal, err = utils.GetMaxFloat64From2dSlice(p.Data)
		if err != nil || maxVal == 0 {
			return
		}
	}

	// Draw axes if enabled
	if p.ShowAxes {
		p.plotAxes(buf, maxVal)
	}

	// Calculate draw area
	drawArea := p.Inner
	if p.ShowAxes {
		drawArea = image.Rect(
			p.Inner.Min.X+yAxisLabelsWidth+1,
			p.Inner.Min.Y,
			p.Inner.Max.X,
			p.Inner.Max.Y-xAxisLabelsHeight-1,
		)
	}

	// Render based on marker type
	switch p.Marker {
	case MarkerBraille:
		// Braille rendering requires Canvas - fallback to dot
		p.renderDot(buf, drawArea, maxVal)
	case MarkerDot:
		p.renderDot(buf, drawArea, maxVal)
	}
}
