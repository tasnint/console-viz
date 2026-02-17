package draw

import (
	"console-viz/styling"
	tb "github.com/nsf/termbox-go"
)

// Init initializes termbox-go and sets up input/output modes
// This must be called before using any drawing functions
// After initialization, the library must be finalized with Close()
func Init() error {
	if err := tb.Init(); err != nil {
		return err
	}
	tb.SetInputMode(tb.InputEsc | tb.InputMouse) // Enable mouse and ESC key detection
	tb.SetOutputMode(tb.Output256)               // Enable 256 color mode
	return nil
}

// Close closes termbox-go and restores terminal state
// Should be called when done with the application
func Close() {
	tb.Close()
}

// TerminalDimensions returns the current terminal width and height
// Syncs termbox state first to ensure accurate dimensions
func TerminalDimensions() (int, int) {
	tb.Sync()
	width, height := tb.Size()
	return width, height
}

// Clear clears the terminal with the default background color from theme
func Clear() {
	theme := styling.GetTheme()
	bgColor := tb.ColorDefault
	if theme.Default.Bg != styling.ColorClear {
		bgColor = tb.Attribute(theme.Default.Bg + 1)
	}
	tb.Clear(tb.ColorDefault, bgColor)
}
