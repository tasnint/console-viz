// +build ignore

package main

import (
	"console-viz/draw"
	"console-viz/styling"
	"fmt"
	"log"

	termbox "github.com/nsf/termbox-go"
)

// Example demonstrating CLI theme switching
func main() {
	// Parse theme from command line
	// Usage: go run theme_demo.go --theme=dark
	//        go run theme_demo.go --theme=light
	cliTheme := styling.ParseThemeFlag()
	
	// Initialize theme from CLI, environment variable, or config file
	if err := styling.InitThemeFromCLI(cliTheme); err != nil {
		log.Printf("Warning: Theme initialization failed: %v", err)
		log.Println("Available themes:", styling.ListThemes())
	}

	// Initialize termbox
	if err := termbox.Init(); err != nil {
		log.Fatalf("Failed to initialize terminal: %v", err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	termbox.SetOutputMode(termbox.Output256)

	// Initialize renderer
	draw.InitRenderer()

	// Create a simple widget to demonstrate theme
	base := draw.NewBase()
	base.Title = fmt.Sprintf("Theme Demo - Current: %s", styling.GetCurrentThemeMode())
	base.SetRect(10, 5, 50, 15)

	// Render
	draw.Render(base)

	// Instructions
	termbox.SetCell(10, 16, 'P', termbox.ColorYellow, termbox.ColorBlack)
	termbox.SetCell(11, 16, 'r', termbox.ColorYellow, termbox.ColorBlack)
	termbox.SetCell(12, 16, 'e', termbox.ColorYellow, termbox.ColorBlack)
	termbox.SetCell(13, 16, 's', termbox.ColorYellow, termbox.ColorBlack)
	termbox.SetCell(14, 16, 's', termbox.ColorYellow, termbox.ColorBlack)
	termbox.SetCell(15, 16, ' ', termbox.ColorYellow, termbox.ColorBlack)
	termbox.SetCell(16, 16, 'T', termbox.ColorYellow, termbox.ColorBlack)
	termbox.SetCell(17, 16, ' ', termbox.ColorYellow, termbox.ColorBlack)
	termbox.SetCell(18, 16, 't', termbox.ColorYellow, termbox.ColorBlack)
	termbox.SetCell(19, 16, 'o', termbox.ColorYellow, termbox.ColorBlack)
	termbox.SetCell(20, 16, ' ', termbox.ColorYellow, termbox.ColorBlack)
	termbox.SetCell(21, 16, 't', termbox.ColorYellow, termbox.ColorBlack)
	termbox.SetCell(22, 16, 'o', termbox.ColorYellow, termbox.ColorBlack)
	termbox.SetCell(23, 16, 'g', termbox.ColorYellow, termbox.ColorBlack)
	termbox.SetCell(24, 16, 'g', termbox.ColorYellow, termbox.ColorBlack)
	termbox.SetCell(25, 16, 'l', termbox.ColorYellow, termbox.ColorBlack)
	termbox.SetCell(26, 16, 'e', termbox.ColorYellow, termbox.ColorBlack)
	termbox.Flush()

	// Event loop
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			if ev.Key == termbox.KeyEsc {
				break
			}
			if ev.Ch == 't' || ev.Ch == 'T' {
				// Toggle theme
				if err := styling.ToggleThemeMode(); err != nil {
					log.Printf("Failed to toggle theme: %v", err)
				} else {
					base.Title = fmt.Sprintf("Theme Demo - Current: %s", styling.GetCurrentThemeMode())
					draw.Render(base)
				}
			}
		}
	}
}
