package styling

import (
	"flag"
	"fmt"
	"os"
)

// ParseThemeFlag parses the --theme flag from command-line arguments
// Returns the theme name if provided, empty string otherwise
// Example: app --theme=dark or app --theme light
func ParseThemeFlag() string {
	var theme string
	flag.StringVar(&theme, "theme", "", "Theme to use: dark, light, or default")
	flag.Parse()
	return theme
}

// ParseThemeFlagWithArgs parses theme flag from custom argument slice
// Useful when you want to parse flags without affecting os.Args
func ParseThemeFlagWithArgs(args []string) string {
	fs := flag.NewFlagSet("theme", flag.ContinueOnError)
	var theme string
	fs.StringVar(&theme, "theme", "", "Theme to use: dark, light, or default")
	fs.Parse(args)
	return theme
}

// PrintThemeHelp prints help information about theme options
func PrintThemeHelp() {
	fmt.Fprintf(os.Stderr, "Theme Options:\n")
	fmt.Fprintf(os.Stderr, "  --theme=dark|light|default    Set theme mode\n")
	fmt.Fprintf(os.Stderr, "  CONSOLE_VIZ_THEME=dark|light  Set via environment variable\n")
	fmt.Fprintf(os.Stderr, "\nAvailable themes: %v\n", ListThemes())
}
