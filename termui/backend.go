// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	// termbox-go is a low level library: that puts terminal into raw mode,
	// polls keyboard/mouse events and manages terminal colors
	tb "github.com/nsf/termbox-go"
)

// Init initializes termbox-go and is required to render anything.
// After initialization, the library must be finalized with `Close`.
func Init() error {
	if err := tb.Init(); err != nil {
		return err
	}
	tb.SetInputMode(tb.InputEsc | tb.InputMouse) // set up that allows detection of mouse clicks and esc key
	tb.SetOutputMode(tb.Output256)               // makes 256 colors available for the terminal
	return nil                                   // nil means no error
}

// Close closes termbox-go.
func Close() {
	tb.Close() // termbox-go close and cleanup
}

func TerminalDimensions() (int, int) {
	tb.Sync()                  // sync termbox's internal state with the actual terminal, ensures we dont get stale dimensions
	width, height := tb.Size() // retrieve terminal dimensions
	return width, height
}

func Clear() {
	tb.Clear(tb.ColorDefault, tb.Attribute(Theme.Default.Bg+1))
}
