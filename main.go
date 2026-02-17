package main

import (
	"console-viz/draw"
	"console-viz/styling"
	"log"
)

func main() {
	// Initialize theme from CLI, environment variable, or config file
	cliTheme := styling.ParseThemeFlag()
	if err := styling.InitThemeFromCLI(cliTheme); err != nil {
		log.Printf("Warning: Theme initialization failed: %v", err)
	}

	// Initialize termbox via draw package
	if err := draw.Init(); err != nil {
		log.Fatalf("Failed to initialize terminal: %v", err)
	}
	defer draw.Close()

	// Initialize renderer
	draw.InitRenderer()

	// Clear terminal
	draw.Clear()

	// Example: print terminal dimensions
	w, h := draw.TerminalDimensions()
	log.Printf("Terminal size: %dx%d", w, h)
	log.Printf("Current theme: %s", styling.GetCurrentThemeMode())

	// Wait for ESC to exit using draw.PollEvents()
	for e := range draw.PollEvents() {
		if e.Type == draw.KeyboardEvent && e.ID == "<Escape>" {
			break
		}
	}
}
