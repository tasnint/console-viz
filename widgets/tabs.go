package widgets

import (
	"console-viz/draw"
	"console-viz/styling"
	"console-viz/utils"
	"image"
)

// TabPane displays a tab bar with selectable tabs
// Used for switching between different views/panels
type TabPane struct {
	draw.Base
	TabNames         []string        // Names of tabs
	ActiveTabIndex   int             // Currently active tab index
	ActiveTabStyle   styling.Style   // Style for active tab
	InactiveTabStyle styling.Style   // Style for inactive tabs
}

// NewTabPane creates a new TabPane widget with the given tab names
func NewTabPane(names ...string) *TabPane {
	theme := styling.GetTheme()
	return &TabPane{
		Base:            *draw.NewBase(),
		TabNames:       names,
		ActiveTabStyle:   theme.Tab.Active,
		InactiveTabStyle: theme.Tab.Inactive,
	}
}

// FocusLeft moves focus to the left tab
func (tp *TabPane) FocusLeft() {
	if tp.ActiveTabIndex > 0 {
		tp.ActiveTabIndex--
	}
}

// FocusRight moves focus to the right tab
func (tp *TabPane) FocusRight() {
	if tp.ActiveTabIndex < len(tp.TabNames)-1 {
		tp.ActiveTabIndex++
	}
}

// Draw renders the tab pane widget
func (tp *TabPane) Draw(buf *draw.Buffer) {
	tp.Base.Draw(buf)

	x := tp.Inner.Min.X
	for i, name := range tp.TabNames {
		// Choose style based on active state
		style := tp.InactiveTabStyle
		if i == tp.ActiveTabIndex {
			style = tp.ActiveTabStyle
		}

		// Draw tab name
		trimmedName := utils.TrimString(name, tp.Inner.Max.X-x)
		buf.SetString(
			trimmedName,
			style,
			image.Pt(x, tp.Inner.Min.Y),
		)

		x += 1 + len(name)

		// Draw separator between tabs
		if i < len(tp.TabNames)-1 && x < tp.Inner.Max.X {
			buf.SetCell(
				draw.NewCell(styling.VERTICAL_LINE, styling.NewStyle(styling.ColorWhite)),
				image.Pt(x, tp.Inner.Min.Y),
			)
		}

		x += 2
	}
}
