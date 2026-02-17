package termui

// Color is an integer from -1 to 255
// -1 = ColorClear
// 0-255 = Xterm colors
type Color int

// ColorClear clears the Fg or Bg color of a Style
const ColorClear Color = -1 // denotes a clear/transparent colour

// Basic terminal colors
const (
	ColorBlack   Color = 0
	ColorRed     Color = 1
	ColorGreen   Color = 2
	ColorYellow  Color = 3
	ColorBlue    Color = 4
	ColorMagenta Color = 5
	ColorCyan    Color = 6
	ColorWhite   Color = 7
)

type Modifier uint

const (
	// ModifierClear clears any modifiers
	// modifiers are used for things like bold, underline, reverse etc. they are represented as bit flags, so we can combine them using bitwise OR
	ModifierClear     Modifier = 0
	ModifierBold      Modifier = 1 << 9
	ModifierUnderline Modifier = 1 << 10
	ModifierReverse   Modifier = 1 << 11
)

// Style represents the style of one terminal cell
type Style struct {
	Fg       Color
	Bg       Color
	Modifier Modifier
}

// StyleClear represents a default Style, with no colors or modifiers
var StyleClear = Style{
	Fg:       ColorClear,
	Bg:       ColorClear,
	Modifier: ModifierClear,
}

// NewStyle takes 1 to 3 arguments
// 1st argument = Fg <font color>
// 2nd argument = optional Bg <background color>
// 3rd argument = optional Modifier <format modifier: bold, underline, reverse etc.>
func NewStyle(fg Color, args ...interface{}) Style {
	bg := ColorClear
	modifier := ModifierClear
	if len(args) >= 1 {
		bg = args[0].(Color) // if there is at least one argument, we treat it as the background color
	}
	if len(args) == 2 {
		modifier = args[1].(Modifier) // if there are two arguments, we treat the second one as the modifier
	}
	return Style{
		fg,
		bg,
		modifier,
	}
}
