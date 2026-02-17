package draw

import (
	"console-viz/styling"
	"strings"
)

// ============================================================================
// Style Parser - Embedded Style Markup Syntax
// ============================================================================
//
// ParseStyles parses strings with embedded style markup and converts them
// to styled cells. This enables rich text formatting within widgets.
//
// Syntax: [text](fg:<color>,bg:<color>,mod:<modifier>)
// Example: [Hello World](fg:red,bg:blue,mod:bold)
//
// Features:
// - Supports nested brackets
// - Graceful error handling (falls back to plain text on syntax errors)
// - All style fields are optional
// - Order-independent style specification

// Parser tokens for style markup syntax
const (
	tokenFg       = "fg"  // foreground color keyword
	tokenBg       = "bg"  // background color keyword
	tokenModifier = "mod" // modifier keyword

	tokenItemSeparator  = "," // separates style items (fg:red,bg:blue)
	tokenValueSeparator = ":" // separates key from value (fg:red)

	tokenBeginStyledText = '[' // starts styled text section
	tokenEndStyledText   = ']' // ends styled text section

	tokenBeginStyle = '(' // starts style definition
	tokenEndStyle   = ')' // ends style definition
)

// ParserState represents the current state of the style parser
type ParserState uint

const (
	// ParserStateDefault - parsing normal text, not inside any markup
	ParserStateDefault ParserState = iota
	// ParserStateStyleItems - parsing style definition inside (...)
	ParserStateStyleItems
	// ParserStateStyledText - parsing text content inside [...]
	ParserStateStyledText
)

// ColorMap maps string color names to Color constants
// Can be extended by users to add custom color names
var ColorMap = map[string]styling.Color{
	"red":     styling.ColorRed,
	"blue":    styling.ColorBlue,
	"black":   styling.ColorBlack,
	"cyan":    styling.ColorCyan,
	"yellow":  styling.ColorYellow,
	"white":   styling.ColorWhite,
	"clear":   styling.ColorClear,
	"green":   styling.ColorGreen,
	"magenta": styling.ColorMagenta,
}

// ModifierMap maps string modifier names to Modifier constants
var ModifierMap = map[string]styling.Modifier{
	"bold":      styling.ModifierBold,
	"underline": styling.ModifierUnderline,
	"reverse":   styling.ModifierReverse,
}

// RegisterColor adds a custom color name to the parser
// Useful for extending the parser with custom color names
// Example: RegisterColor("orange", styling.ColorYellow)
func RegisterColor(name string, color styling.Color) {
	ColorMap[strings.ToLower(name)] = color
}

// RegisterModifier adds a custom modifier name to the parser
// Useful for extending the parser with custom modifier names
func RegisterModifier(name string, modifier styling.Modifier) {
	ModifierMap[strings.ToLower(name)] = modifier
}

// parseStyleString parses a style string like "fg:red,bg:blue,mod:bold"
// and returns a Style with the specified attributes
// Uses defaultStyle as a base and only overrides specified attributes
func parseStyleString(styleStr string, defaultStyle styling.Style) styling.Style {
	style := defaultStyle // Start with default style

	// Split by comma to get individual style items
	items := strings.Split(styleStr, tokenItemSeparator)

	for _, item := range items {
		item = strings.TrimSpace(item) // Remove whitespace
		parts := strings.Split(item, tokenValueSeparator)

		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			switch key {
			case tokenFg:
				// Set foreground color
				if color, ok := ColorMap[strings.ToLower(value)]; ok {
					style.Fg = color
				}
			case tokenBg:
				// Set background color
				if color, ok := ColorMap[strings.ToLower(value)]; ok {
					style.Bg = color
				}
			case tokenModifier:
				// Set modifier (can be combined with bitwise OR)
				if modifier, ok := ModifierMap[strings.ToLower(value)]; ok {
					style.Modifier = modifier
				}
			}
		}
	}

	return style
}

// ParseStyles parses a string with embedded style markup and returns styled cells
//
// Syntax: [text](fg:<color>,bg:<color>,mod:<modifier>)
//
// Examples:
//   - "Hello [World](fg:red)" → "Hello " (default) + "World" (red)
//   - "[Bold Text](mod:bold)" → "Bold Text" (bold)
//   - "Normal [Red](fg:red) and [Blue](fg:blue)" → Mixed styles
//
// Features:
//   - Supports nested brackets: [[inner]](fg:red)
//   - Graceful error handling: invalid syntax falls back to plain text
//   - All style fields are optional
//   - Order-independent: (bg:blue,fg:red) same as (fg:red,bg:blue)
//
// Parameters:
//   - s: Input string with optional style markup
//   - defaultStyle: Style to use for unmarked text
//
// Returns:
//   - Slice of styled cells ready for rendering
func ParseStyles(s string, defaultStyle styling.Style) []Cell {
	if s == "" {
		return []Cell{}
	}

	cells := make([]Cell, 0)
	runes := []rune(s)
	state := ParserStateDefault
	styledText := make([]rune, 0) // Text inside [...]
	styleItems := make([]rune, 0)  // Style definition inside (...)
	bracketDepth := 0              // Track nested brackets

	// Reset parser state to default
	reset := func() {
		styledText = make([]rune, 0)
		styleItems = make([]rune, 0)
		state = ParserStateDefault
		bracketDepth = 0
	}

	// Helper function to convert runes to cells
	runesToCells := func(runes []rune, style styling.Style) []Cell {
		result := make([]Cell, len(runes))
		for i, r := range runes {
			result[i] = Cell{
				Rune:  r,
				Style: style,
			}
		}
		return result
	}

	// Rollback: treat collected text as plain text and reset
	rollback := func() {
		// Convert styled text to plain cells
		if len(styledText) > 0 {
			cells = append(cells, runesToCells(styledText, defaultStyle)...)
		}
		// Convert style items to plain cells (in case of syntax error)
		if len(styleItems) > 0 {
			cells = append(cells, runesToCells(styleItems, defaultStyle)...)
		}
		reset()
	}

	// Remove first and last rune (used to strip brackets/parentheses)
	chop := func(r []rune) []rune {
		if len(r) < 2 {
			return r
		}
		return r[1 : len(r)-1]
	}

	// Main parsing loop
	for i, r := range runes {
		switch state {
		case ParserStateDefault:
			// Looking for start of styled text '['
			if r == tokenBeginStyledText {
				state = ParserStateStyledText
				bracketDepth = 1
				styledText = append(styledText, r)
			} else {
				// Plain text - add as default style
				cells = append(cells, Cell{
					Rune:  r,
					Style: defaultStyle,
				})
			}

		case ParserStateStyledText:
			// Collecting text inside [...]
			switch {
			case bracketDepth == 0:
				// Text portion complete, now expecting '('
				if r == tokenBeginStyle {
					// Start collecting style definition
					state = ParserStateStyleItems
					styleItems = append(styleItems, r)
				} else {
					// Invalid syntax - expected '(' but got something else
					rollback()
					// Handle the current character
					if r == tokenBeginStyledText {
						// Start new styled section
						state = ParserStateStyledText
						bracketDepth = 1
						styledText = append(styledText, r)
					} else {
						// Add as plain text
						cells = append(cells, Cell{
							Rune:  r,
							Style: defaultStyle,
						})
					}
				}

			case i == len(runes)-1:
				// String ended while parsing - incomplete syntax
				rollback()

			case r == tokenBeginStyledText:
				// Nested '[' - increment depth
				bracketDepth++
				styledText = append(styledText, r)

			case r == tokenEndStyledText:
				// Closing ']' - decrement depth
				bracketDepth--
				styledText = append(styledText, r)

			default:
				// Normal character inside brackets
				styledText = append(styledText, r)
			}

		case ParserStateStyleItems:
			// Collecting style definition inside (...)
			styleItems = append(styleItems, r)

			if r == tokenEndStyle {
				// Style definition complete!
				// Parse the style string and create styled cells
				styleStr := string(chop(styleItems))
				style := parseStyleString(styleStr, defaultStyle)

				// Create styled cells from the text
				textRunes := chop(styledText)
				styledCells := runesToCells(textRunes, style)
				cells = append(cells, styledCells...)

				reset()
			} else if i == len(runes)-1 {
				// String ended without closing ')' - incomplete syntax
				rollback()
			}
		}
	}

	// Handle any remaining unprocessed text
	if state != ParserStateDefault {
		rollback()
	}

	return cells
}

// ParseStylesSimple is a convenience function that parses styles with a simple default
// Uses StyleClear as the default style
func ParseStylesSimple(s string) []Cell {
	return ParseStyles(s, styling.StyleClear)
}

// HasStyleMarkup checks if a string contains style markup syntax
// Useful for conditional parsing or validation
func HasStyleMarkup(s string) bool {
	return strings.Contains(s, "[") && strings.Contains(s, "]") &&
		strings.Contains(s, "(") && strings.Contains(s, ")")
}

// StripStyleMarkup removes all style markup from a string, returning plain text
// Useful for extracting text content without styles
func StripStyleMarkup(s string) string {
	runes := []rune(s)
	result := make([]rune, 0)
	inMarkup := false
	bracketDepth := 0

	for i, r := range runes {
		switch {
		case r == tokenBeginStyledText:
			inMarkup = true
			bracketDepth = 1
		case r == tokenEndStyledText && inMarkup:
			bracketDepth--
			if bracketDepth == 0 {
				// Skip until we find the closing ')'
				for j := i + 1; j < len(runes); j++ {
					if runes[j] == tokenEndStyle {
						inMarkup = false
						break
					}
				}
			}
		case r == tokenBeginStyledText && inMarkup:
			bracketDepth++
		case !inMarkup:
			result = append(result, r)
		}
	}

	return string(result)
}
