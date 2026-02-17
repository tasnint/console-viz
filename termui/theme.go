// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

// reusable color palette for widgets that need multiple colors
// widgets cycle through these colors for different data series
var StandardColors = []Color{
	ColorRed,
	ColorGreen,
	ColorYellow,
	ColorBlue,
	ColorMagenta,
	ColorCyan,
	ColorWhite,
}

// same as standard colors but in Style form, for widgets that need multiple styles
var StandardStyles = []Style{
	NewStyle(ColorRed),
	NewStyle(ColorGreen),
	NewStyle(ColorYellow),
	NewStyle(ColorBlue),
	NewStyle(ColorMagenta),
	NewStyle(ColorCyan),
	NewStyle(ColorWhite),
}

// master theme container holding themes for every widget type
// each widget has its own subtheme struct with widget specific styling
type RootTheme struct {
	Default Style // fallback default style

	Block BlockTheme

	BarChart        BarChartTheme
	Gauge           GaugeTheme
	Plot            PlotTheme
	List            ListTheme
	Tree            TreeTheme
	Paragraph       ParagraphTheme
	PieChart        PieChartTheme
	Sparkline       SparklineTheme
	StackedBarChart StackedBarChartTheme
	Tab             TabTheme
	Table           TableTheme
}

type BlockTheme struct {
	Title  Style // style for the title text at the top of the widget
	Border Style // style for border lines of the block
}

type BarChartTheme struct {
	Bars   []Color // colors for each bar in the chart, if there are more bars than colors, it will cycle through the colors
	Nums   []Style // styles for the numbers displayed on top of the bars, if there are more bars than styles, it will cycle through the styles
	Labels []Style // styles for the labels displayed below the bars, if there are more bars than styles, it will cycle through the styles
}

type GaugeTheme struct {
	Bar   Color // color for the filled portion of the gauge
	Label Style // style for the percentage label displayed on the gauge
}

type PlotTheme struct {
	Lines []Color // colors for each line in the plot, if there are more lines than colors, it will cycle through the colors
	Axes  Color   // color for the X and Y axes of the plot
}

type ListTheme struct {
	Text Style // style for the text of each list item
}

type TreeTheme struct {
	Text      Style // style for the text of each tree node
	Collapsed rune  // style for symbol indicating a collapsed node, default is "+"
	Expanded  rune  // style for symbol indicating an expanded node, default is "-"
}

type ParagraphTheme struct {
	Text Style // style for the text of the paragraph
}

type PieChartTheme struct {
	Slices []Color // colors for each slice of the pie chart, if there are more slices than colors, it will cycle through the colors
}

type SparklineTheme struct {
	Title Style // style for the title of the sparkline
	Line  Color // color for the line of the sparkline
}

type StackedBarChartTheme struct {
	Bars   []Color // colors for each bar segment in the stacked bar chart
	Nums   []Style // styles for the numbers displayed on top of each bar segment
	Labels []Style // styles for the labels displayed below each bar segment
}

type TabTheme struct {
	Active   Style // style for the active tab
	Inactive Style // style for the inactive tab
}

type TableTheme struct {
	Text Style // style for the text in each cell of the table
}

// Theme holds the default Styles and Colors for all widgets.
// You can set default widget Styles by modifying the Theme before creating the widgets.
var Theme = RootTheme{
	Default: NewStyle(ColorWhite),

	Block: BlockTheme{
		Title:  NewStyle(ColorWhite),
		Border: NewStyle(ColorWhite),
	},

	BarChart: BarChartTheme{
		Bars:   StandardColors,
		Nums:   StandardStyles,
		Labels: StandardStyles,
	},

	Paragraph: ParagraphTheme{
		Text: NewStyle(ColorWhite),
	},

	PieChart: PieChartTheme{
		Slices: StandardColors,
	},

	List: ListTheme{
		Text: NewStyle(ColorWhite),
	},

	Tree: TreeTheme{
		Text:      NewStyle(ColorWhite),
		Collapsed: COLLAPSED,
		Expanded:  EXPANDED,
	},

	StackedBarChart: StackedBarChartTheme{
		Bars:   StandardColors,
		Nums:   StandardStyles,
		Labels: StandardStyles,
	},

	Gauge: GaugeTheme{
		Bar:   ColorWhite,
		Label: NewStyle(ColorWhite),
	},

	Sparkline: SparklineTheme{
		Title: NewStyle(ColorWhite),
		Line:  ColorWhite,
	},

	Plot: PlotTheme{
		Lines: StandardColors,
		Axes:  ColorWhite,
	},

	Table: TableTheme{
		Text: NewStyle(ColorWhite),
	},

	Tab: TabTheme{
		Active:   NewStyle(ColorRed),
		Inactive: NewStyle(ColorWhite),
	},
}
