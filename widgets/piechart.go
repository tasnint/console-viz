package widgets

import (
	"console-viz/draw"
	"console-viz/styling"
	"console-viz/utils"
	"image"
	"math"
)

const (
	piechartOffsetUp = -.5 * math.Pi // Start angle (north)
	resolutionFactor = .0001         // Circle resolution
	fullCircle       = 2.0 * math.Pi // Full circle
	xStretch         = 2.0           // Horizontal stretch factor
)

// PieChartLabel is a callback function for formatting pie chart labels
type PieChartLabel func(dataIndex int, currentValue float64) string

// PieChart displays a pie chart visualization
type PieChart struct {
	draw.Base
	Data           []float64     // Data values for each slice
	Colors         []styling.Color // Colors for slices (cycled)
	LabelFormatter PieChartLabel // Callback for label formatting
	AngleOffset    float64       // Starting angle offset
}

// NewPieChart creates a new PieChart widget with default settings
func NewPieChart() *PieChart {
	theme := styling.GetTheme()
	return &PieChart{
		Base:       *draw.NewBase(),
		Colors:     theme.PieChart.Slices,
		AngleOffset: piechartOffsetUp,
	}
}

// Draw renders the pie chart widget
func (pc *PieChart) Draw(buf *draw.Buffer) {
	pc.Base.Draw(buf)

	if len(pc.Data) == 0 {
		return
	}

	center := pc.Inner.Min.Add(pc.Inner.Size().Div(2))
	radius := utils.MinFloat64(float64(pc.Inner.Dx()/2)/xStretch, float64(pc.Inner.Dy()/2))

	// Calculate slice sizes
	sum := utils.SumFloat64Slice(pc.Data)
	if sum == 0 {
		return
	}

	sliceSizes := make([]float64, len(pc.Data))
	for i, v := range pc.Data {
		sliceSizes[i] = v / sum * fullCircle
	}

	borderCircle := &circle{center, radius}
	middleCircle := circle{Point: center, radius: radius / 2.0}

	// Draw sectors
	phi := pc.AngleOffset
	for i, size := range sliceSizes {
		for j := 0.0; j < size; j += resolutionFactor {
			borderPoint := borderCircle.at(phi + j)
			line := line{P1: center, P2: borderPoint}
			line.draw(draw.NewCell(styling.SHADED_BLOCKS[1], styling.NewStyle(utils.SelectColor(pc.Colors, i))), buf)
		}
		phi += size
	}

	// Draw labels
	if pc.LabelFormatter != nil {
		phi = pc.AngleOffset
		for i, size := range sliceSizes {
			labelPoint := middleCircle.at(phi + size/2.0)
			if len(pc.Data) == 1 {
				labelPoint = center
			}
			label := pc.LabelFormatter(i, pc.Data[i])
			buf.SetString(
				label,
				styling.NewStyle(utils.SelectColor(pc.Colors, i)),
				image.Pt(labelPoint.X, labelPoint.Y),
			)
			phi += size
		}
	}
}

// circle represents a circle for pie chart calculations
type circle struct {
	image.Point
	radius float64
}

// at returns the point at a given angle phi
func (c *circle) at(phi float64) image.Point {
	x := c.X + int(utils.RoundFloat64(xStretch*c.radius*math.Cos(phi)))
	y := c.Y + int(utils.RoundFloat64(c.radius*math.Sin(phi)))
	return image.Point{X: x, Y: y}
}

// line represents a line between two points
type line struct {
	P1, P2 image.Point
}

// draw draws the line on the buffer
func (l *line) draw(cell draw.Cell, buf *draw.Buffer) {
	isLeftOf := func(p1, p2 image.Point) bool {
		return p1.X <= p2.X
	}
	isTopOf := func(p1, p2 image.Point) bool {
		return p1.Y <= p2.Y
	}

	p1, p2 := l.P1, l.P2
	buf.SetCell(draw.NewCell('*', cell.Style), l.P2)

	width, height := l.size()
	if width > height { // Paint left to right
		if !isLeftOf(p1, p2) {
			p1, p2 = p2, p1
		}
		flip := 1.0
		if !isTopOf(p1, p2) {
			flip = -1.0
		}
		for x := p1.X; x <= p2.X; x++ {
			ratio := float64(height) / float64(width)
			factor := float64(x - p1.X)
			y := ratio * factor * flip
			buf.SetCell(cell, image.Pt(x, int(utils.RoundFloat64(y))+p1.Y))
		}
	} else { // Paint top to bottom
		if !isTopOf(p1, p2) {
			p1, p2 = p2, p1
		}
		flip := 1.0
		if !isLeftOf(p1, p2) {
			flip = -1.0
		}
		for y := p1.Y; y <= p2.Y; y++ {
			ratio := float64(width) / float64(height)
			factor := float64(y - p1.Y)
			x := ratio * factor * flip
			buf.SetCell(cell, image.Pt(int(utils.RoundFloat64(x))+p1.X, y))
		}
	}
}

// size returns the width and height of the line
func (l *line) size() (w, h int) {
	return utils.AbsInt(l.P2.X - l.P1.X), utils.AbsInt(l.P2.Y - l.P1.Y)
}
