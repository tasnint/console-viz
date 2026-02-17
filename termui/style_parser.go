// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"strings"
)

// defining the markup syntax for embedded styles
// this styling is applied to text in widgets
// e.g. [Hello World](fg:red, bg:blue, mod:bold) would render "Hello World" in bold red text on a blue background
const (
	tokenFg       = "fg"  // foreground color keyword (font of text)
	tokenBg       = "bg"  // background color keyword (color of the background behind the text)
	tokenModifier = "mod" // modifier keyword

	tokenItemSeparator  = "," // separates style items
	tokenValueSeparator = ":" // separates key from value

	tokenBeginStyledText = '[' // starts styled text
	tokenEndStyledText   = ']' // ends styled text

	tokenBeginStyle = '(' // starts style definition
	tokenEndStyle   = ')' // ends style definition
)

type parserState uint // creating a custom type parserState based on unsigned integer
// this will be used to track the state of the style parser as it processes the input string

const (
	parserStateDefault parserState = iota // default state, not currently parsing styled text or style items
	parserStateStyleItems
	parserStateStyledText
)

// StyleParserColorMap can be modified to add custom color parsing to text
var StyleParserColorMap = map[string]Color{
	"red":     ColorRed,
	"blue":    ColorBlue,
	"black":   ColorBlack,
	"cyan":    ColorCyan,
	"yellow":  ColorYellow,
	"white":   ColorWhite,
	"clear":   ColorClear,
	"green":   ColorGreen,
	"magenta": ColorMagenta,
}

var modifierMap = map[string]Modifier{
	"bold":      ModifierBold,
	"underline": ModifierUnderline,
	"reverse":   ModifierReverse,
}

// readStyle translates an []rune like `fg:red,mod:bold,bg:white` to a style
func readStyle(runes []rune, defaultStyle Style) Style {
	style := defaultStyle // fallback style
	split := strings.Split(string(runes), tokenItemSeparator)
	for _, item := range split {
		pair := strings.Split(item, tokenValueSeparator) // splits the item into key and value, e.g. "fg:red" -> ["fg", "red"]
		if len(pair) == 2 {                              // styling is correctly done if the key has a value, for example "fg:red" is correct but "fg:" or "fg" is not
			switch pair[0] {
			case tokenFg: // if the key is "fg", we set the foreground color of the style
				style.Fg = StyleParserColorMap[pair[1]]
			case tokenBg: // if the key is "bg", we set the background color of the style
				style.Bg = StyleParserColorMap[pair[1]]
			case tokenModifier: // if the key is "mod", we set the modifier of the style
				style.Modifier = modifierMap[pair[1]]
			}
		}
	}
	return style
}

// ParseStyles parses a string for embedded Styles and returns []Cell with the correct styling.
// Uses defaultStyle for any text without an embedded style.
// Syntax is of the form [text](fg:<color>,mod:<attribute>,bg:<color>).
// Ordering does not matter. All fields are optional.
func ParseStyles(s string, defaultStyle Style) []Cell {
	cells := []Cell{}           // output: array of styled cells
	runes := []rune(s)          // convert string to runes
	state := parserStateDefault // initialize parser state to default
	styledText := []rune{}      // collects text inside [..]
	styleItems := []rune{}      // collects style inside (..)
	squareCount := 0            // tracks nested brackets

	// function that clears everything and returns to starting state
	// after successsfully converted to styled cells, we reset
	reset := func() {
		styledText = []rune{}      // clears the styled text buffer
		styleItems = []rune{}      // clears the style items buffer
		state = parserStateDefault // resets the parser state to default
		squareCount = 0
	}

	// if style parsing fails at any point, we want to roll back to treating the text as unstyled and reset the parser
	rollback := func() { // convert to plain cells and then reset
		cells = append(cells, RunesToStyledCells(styledText, defaultStyle)...) // convert styled text to og input example: "Hello World" -> ['H', 'e', 'l', 'l', 'o', ' ', 'W', 'o', 'r', 'l', 'd']
		cells = append(cells, RunesToStyledCells(styleItems, defaultStyle)...) // convert style items to og input
		reset()
	}

	// chop first and last runes
	chop := func(s []rune) []rune {
		return s[1 : len(s)-1] // removes first and last characters
	}

	for i, _rune := range runes { // loops through every character in the rune
		switch state { // checks current state
		case parserStateDefault: // if current state is default:
			if _rune == tokenBeginStyledText { // is character '['?]
				state = parserStateStyledText          // yes: switch to state of parsing styled text
				squareCount = 1                        // square bracket count is now 1
				styledText = append(styledText, _rune) // adds '[' to styled text buffer
			} else {
				// normal character - add as plain cell
				cells = append(cells, Cell{_rune, defaultStyle})
			}

		// STATE: STYLEDTEXT - inside [...], collecting text to be styled
		case parserStateStyledText:
			switch {

			// text portion complete (saw ']'), now expecting '('
			case squareCount == 0:
				switch _rune {
				case tokenBeginStyle: // saw '(' - start collecting style
					state = parserStateStyleItems          // change state to parsing style items
					styleItems = append(styleItems, _rune) // add '(' to style items buffer
				default: // didn't see '(' - invalid syntax!
					rollback() // treat as plain text
					switch _rune {
					case tokenBeginStyledText: // but if '[', start new styled section
						state = parserStateStyledText          // switch to state of parsing styled text
						squareCount = 1                        // square bracket count is now 1
						styleItems = append(styleItems, _rune) // add '[' to style items buffer
					default:
						cells = append(cells, Cell{_rune, defaultStyle}) // otherwise, just add the character as a plain cell
					}
				}

			// string ended while still parsing - incomplete syntax
			case len(runes) == i+1:
				rollback() // might be a bug? rollback should be called after character is appended to cell or styletext buffer
				styledText = append(styledText, _rune)

			// nested '[' - increment bracket counter
			case _rune == tokenBeginStyledText:
				squareCount++
				styledText = append(styledText, _rune)

			// closing ']' - decrement counter (when 0, text portion done)
			case _rune == tokenEndStyledText:
				squareCount--
				styledText = append(styledText, _rune)

			// normal character inside brackets - collect it
			default:
				styledText = append(styledText, _rune)
			}

		// STATE: STYLEITEMS - inside (...), collecting style definition
		case parserStateStyleItems:
			styleItems = append(styleItems, _rune) // always collect the character

			if _rune == tokenEndStyle { // saw ')' - SUCCESS! style complete
				// parse style string: "(fg:red)" → "fg:red" → Style{Fg: ColorRed}
				style := readStyle(chop(styleItems), defaultStyle)
				// create styled cells: "[Hello]" → "Hello" → []Cell with style
				// ... unpacks the slice so each Cell is appended individually
				cells = append(cells, RunesToStyledCells(chop(styledText), style)...) // ... is Go's spread operator, it unpacks the slice so each Cell is appended individually, e.g. "Hello" -> ['H', 'e', 'l', 'l', 'o'] with the specified style
				reset()                                                               // clear buffers, back to default state
			} else if len(runes) == i+1 { // string ended without ')' - incomplete
				rollback() //
			}
		}
	}

	return cells // return all parsed cells
}
