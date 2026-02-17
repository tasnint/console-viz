package styling

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//
// -----------------------------------------------------
// Shared Color Palettes
// -----------------------------------------------------
//

var StandardColors = []Color{
	ColorRed,
	ColorGreen,
	ColorYellow,
	ColorBlue,
	ColorMagenta,
	ColorCyan,
	ColorWhite,
}

var StandardStyles = []Style{
	NewStyle(ColorRed),
	NewStyle(ColorGreen),
	NewStyle(ColorYellow),
	NewStyle(ColorBlue),
	NewStyle(ColorMagenta),
	NewStyle(ColorCyan),
	NewStyle(ColorWhite),
}

//
// -----------------------------------------------------
// Root Theme Structure (Widget Compatible)
// -----------------------------------------------------
//

type RootTheme struct {
	Default Style

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

	Accessibility AccessibilityVariants
}

//
// -----------------------------------------------------
// Accessibility Variants (no recursion issues)
// -----------------------------------------------------
//

type AccessibilityVariants struct {
	HighContrast *RootTheme
	ColorBlind   *RootTheme
}

//
// -----------------------------------------------------
// Theme Registry (Multiple Theme Support)
// -----------------------------------------------------
//

var themeRegistry = map[string]*RootTheme{}

var activeTheme *RootTheme
var currentThemeMode string // "dark", "light", or "default"

// ThemeMode represents the theme mode
type ThemeMode string

const (
	ThemeModeDark   ThemeMode = "dark"
	ThemeModeLight  ThemeMode = "light"
	ThemeModeDefault ThemeMode = "default"
)

// ThemeChangeCallback is called when theme changes
type ThemeChangeCallback func(mode ThemeMode)

var themeChangeCallbacks []ThemeChangeCallback

// RegisterThemeChangeCallback registers a callback to be called when theme changes
func RegisterThemeChangeCallback(callback ThemeChangeCallback) {
	themeChangeCallbacks = append(themeChangeCallbacks, callback)
}

// notifyThemeChange notifies all registered callbacks of theme change
func notifyThemeChange(mode ThemeMode) {
	for _, callback := range themeChangeCallbacks {
		callback(mode)
	}
}

// RegisterTheme registers a theme by name
func RegisterTheme(name string, theme *RootTheme) {
	themeRegistry[name] = theme
}

// SwitchTheme switches active theme at runtime
func SwitchTheme(name string) error {
	t, ok := themeRegistry[name]
	if !ok {
		return fmt.Errorf("theme '%s' not found", name)
	}
	activeTheme = t
	
	// Update mode based on theme name
	if name == "dark" {
		currentThemeMode = string(ThemeModeDark)
		notifyThemeChange(ThemeModeDark)
	} else if name == "light" {
		currentThemeMode = string(ThemeModeLight)
		notifyThemeChange(ThemeModeLight)
	} else {
		currentThemeMode = string(ThemeModeDefault)
		notifyThemeChange(ThemeModeDefault)
	}
	
	return nil
}

// GetTheme returns currently active theme
func GetTheme() *RootTheme {
	if activeTheme == nil {
		return themeRegistry["default"]
	}
	return activeTheme
}

// GetCurrentThemeMode returns the current theme mode
func GetCurrentThemeMode() ThemeMode {
	if currentThemeMode == "" {
		return ThemeModeDefault
	}
	return ThemeMode(currentThemeMode)
}

// ToggleDarkMode switches to dark mode theme
func ToggleDarkMode() error {
	return SwitchTheme("dark")
}

// ToggleLightMode switches to light mode theme
func ToggleLightMode() error {
	return SwitchTheme("light")
}

// ToggleThemeMode toggles between dark and light mode
func ToggleThemeMode() error {
	current := GetCurrentThemeMode()
	if current == ThemeModeDark {
		return ToggleLightMode()
	}
	return ToggleDarkMode()
}

// ListThemes returns a list of all registered theme names
func ListThemes() []string {
	themes := make([]string, 0, len(themeRegistry))
	for name := range themeRegistry {
		themes = append(themes, name)
	}
	return themes
}

// HasTheme checks if a theme with the given name exists
func HasTheme(name string) bool {
	_, exists := themeRegistry[name]
	return exists
}

//
// -----------------------------------------------------
// CLI and Environment Variable Support
// -----------------------------------------------------
//

// InitThemeFromCLI initializes theme from command-line arguments, environment variables, or config file
// Priority: CLI args > Environment variable > Config file > Default
// This should be called early in application startup
func InitThemeFromCLI(cliTheme string) error {
	// Priority 1: CLI argument (highest priority)
	if cliTheme != "" {
		if err := SwitchTheme(cliTheme); err != nil {
			return fmt.Errorf("invalid CLI theme '%s': %w", cliTheme, err)
		}
		// Save preference
		saveThemePreference(cliTheme)
		return nil
	}

	// Priority 2: Environment variable
	if envTheme := os.Getenv("CONSOLE_VIZ_THEME"); envTheme != "" {
		envTheme = strings.ToLower(strings.TrimSpace(envTheme))
		if err := SwitchTheme(envTheme); err != nil {
			return fmt.Errorf("invalid environment theme '%s': %w", envTheme, err)
		}
		// Save preference
		saveThemePreference(envTheme)
		return nil
	}

	// Priority 3: Config file
	if savedTheme := loadThemePreference(); savedTheme != "" {
		if err := SwitchTheme(savedTheme); err == nil {
			return nil
		}
		// If saved theme is invalid, continue to default
	}

	// Priority 4: Default (already set in init())
	return nil
}

// getConfigFilePath returns the path to the theme config file
func getConfigFilePath() string {
	// Try to use user's config directory
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			configDir = filepath.Join(homeDir, ".config")
		}
	}
	
	if configDir != "" {
		appConfigDir := filepath.Join(configDir, "console-viz")
		os.MkdirAll(appConfigDir, 0755) // Create directory if it doesn't exist
		return filepath.Join(appConfigDir, "theme.json")
	}
	
	// Fallback to current directory
	return "theme.json"
}

// saveThemePreference saves the theme preference to a config file
func saveThemePreference(themeName string) {
	configPath := getConfigFilePath()
	config := map[string]string{
		"theme": themeName,
	}
	
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return // Silently fail if we can't save
	}
	
	os.WriteFile(configPath, data, 0644)
}

// loadThemePreference loads the theme preference from config file
func loadThemePreference() string {
	configPath := getConfigFilePath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}
	
	var config map[string]string
	if err := json.Unmarshal(data, &config); err != nil {
		return ""
	}
	
	theme, ok := config["theme"]
	if !ok {
		return ""
	}
	
	return strings.ToLower(strings.TrimSpace(theme))
}

// SetThemeFromString sets theme from a string (useful for CLI/API)
// Valid values: "dark", "light", "default"
func SetThemeFromString(themeStr string) error {
	themeStr = strings.ToLower(strings.TrimSpace(themeStr))
	if !HasTheme(themeStr) {
		return fmt.Errorf("theme '%s' not found. Available themes: %v", themeStr, ListThemes())
	}
	
	if err := SwitchTheme(themeStr); err != nil {
		return err
	}
	
	// Save preference
	saveThemePreference(themeStr)
	return nil
}

//
// -----------------------------------------------------
// Theme Inheritance
// -----------------------------------------------------
//

// InheritFrom fills missing values from base theme
func (t *RootTheme) InheritFrom(base *RootTheme) {

	if t.Default == (Style{}) {
		t.Default = base.Default
	}

	// Block
	if t.Block.Title == (Style{}) {
		t.Block.Title = base.Block.Title
	}
	if t.Block.Border == (Style{}) {
		t.Block.Border = base.Block.Border
	}

	// Gauge
	if t.Gauge.Bar == 0 {
		t.Gauge.Bar = base.Gauge.Bar
	}
	if t.Gauge.Label == (Style{}) {
		t.Gauge.Label = base.Gauge.Label
	}

	// Plot
	if len(t.Plot.Lines) == 0 {
		t.Plot.Lines = base.Plot.Lines
	}
	if t.Plot.Axes == 0 {
		t.Plot.Axes = base.Plot.Axes
	}

	// List
	if t.List.Text == (Style{}) {
		t.List.Text = base.List.Text
	}

	// Paragraph
	if t.Paragraph.Text == (Style{}) {
		t.Paragraph.Text = base.Paragraph.Text
	}

	// Table
	if t.Table.Text == (Style{}) {
		t.Table.Text = base.Table.Text
	}

	// Tab
	if t.Tab.Active == (Style{}) {
		t.Tab.Active = base.Tab.Active
	}
	if t.Tab.Inactive == (Style{}) {
		t.Tab.Inactive = base.Tab.Inactive
	}

	// BarChart
	if len(t.BarChart.Bars) == 0 {
		t.BarChart.Bars = base.BarChart.Bars
	}
	if len(t.BarChart.Nums) == 0 {
		t.BarChart.Nums = base.BarChart.Nums
	}
	if len(t.BarChart.Labels) == 0 {
		t.BarChart.Labels = base.BarChart.Labels
	}

	// Sparkline
	if t.Sparkline.Title == (Style{}) {
		t.Sparkline.Title = base.Sparkline.Title
	}
	if t.Sparkline.Line == 0 {
		t.Sparkline.Line = base.Sparkline.Line
	}

	// PieChart
	if len(t.PieChart.Slices) == 0 {
		t.PieChart.Slices = base.PieChart.Slices
	}

	// StackedBarChart
	if len(t.StackedBarChart.Bars) == 0 {
		t.StackedBarChart.Bars = base.StackedBarChart.Bars
	}
	if len(t.StackedBarChart.Nums) == 0 {
		t.StackedBarChart.Nums = base.StackedBarChart.Nums
	}
	if len(t.StackedBarChart.Labels) == 0 {
		t.StackedBarChart.Labels = base.StackedBarChart.Labels
	}
}

//
// -----------------------------------------------------
// Validation
// -----------------------------------------------------
//

// Validate fills missing values using default theme
func (t *RootTheme) Validate() {
	defaultTheme, ok := themeRegistry["default"]
	if ok && t != defaultTheme {
		t.InheritFrom(defaultTheme)
	}
}

//
// -----------------------------------------------------
// JSON Loader
// -----------------------------------------------------
//

func LoadThemeFromFile(name, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var theme RootTheme
	if err := json.NewDecoder(file).Decode(&theme); err != nil {
		return err
	}

	theme.Validate()
	RegisterTheme(name, &theme)

	return nil
}

//
// -----------------------------------------------------
// Widget SubThemes (unchanged from original)
// -----------------------------------------------------
//

type BlockTheme struct {
	Title  Style
	Border Style
}

type BarChartTheme struct {
	Bars   []Color
	Nums   []Style
	Labels []Style
}

type GaugeTheme struct {
	Bar   Color
	Label Style
}

type PlotTheme struct {
	Lines []Color
	Axes  Color
}

type ListTheme struct {
	Text Style
}

type TreeTheme struct {
	Text      Style
	Collapsed rune
	Expanded  rune
}

type ParagraphTheme struct {
	Text Style
}

type PieChartTheme struct {
	Slices []Color
}

type SparklineTheme struct {
	Title Style
	Line  Color
}

type StackedBarChartTheme struct {
	Bars   []Color
	Nums   []Style
	Labels []Style
}

type TabTheme struct {
	Active   Style
	Inactive Style
}

type TableTheme struct {
	Text Style
}

//
// -----------------------------------------------------
// Theme Definitions (Dark, Light, Default)
// -----------------------------------------------------
//

// createDarkTheme creates a dark mode theme (dark background, light text)
func createDarkTheme() *RootTheme {
	return &RootTheme{
		Default: NewStyle(ColorWhite, ColorBlack),

		Block: BlockTheme{
			Title:  NewStyle(ColorWhite, ColorBlack, ModifierBold),
			Border: NewStyle(ColorWhite, ColorBlack),
		},

		BarChart: BarChartTheme{
			Bars:   StandardColors,
			Nums:   StandardStyles,
			Labels: StandardStyles,
		},

		Gauge: GaugeTheme{
			Bar:   ColorCyan,
			Label: NewStyle(ColorWhite, ColorBlack),
		},

		Plot: PlotTheme{
			Lines: StandardColors,
			Axes:  ColorWhite,
		},

		List: ListTheme{
			Text: NewStyle(ColorWhite, ColorBlack),
		},

		Tree: TreeTheme{
			Text:      NewStyle(ColorWhite, ColorBlack),
			Collapsed: '+',
			Expanded:  '−',
		},

		Paragraph: ParagraphTheme{
			Text: NewStyle(ColorWhite, ColorBlack),
		},

		PieChart: PieChartTheme{
			Slices: StandardColors,
		},

		Sparkline: SparklineTheme{
			Title: NewStyle(ColorWhite, ColorBlack),
			Line:  ColorCyan,
		},

		StackedBarChart: StackedBarChartTheme{
			Bars:   StandardColors,
			Nums:   StandardStyles,
			Labels: StandardStyles,
		},

		Table: TableTheme{
			Text: NewStyle(ColorWhite, ColorBlack),
		},

		Tab: TabTheme{
			Active:   NewStyle(ColorRed, ColorBlack, ModifierBold),
			Inactive: NewStyle(ColorWhite, ColorBlack),
		},
	}
}

// createLightTheme creates a light mode theme (light background, dark text)
func createLightTheme() *RootTheme {
	return &RootTheme{
		Default: NewStyle(ColorBlack, ColorWhite),

		Block: BlockTheme{
			Title:  NewStyle(ColorBlack, ColorWhite, ModifierBold),
			Border: NewStyle(ColorBlack, ColorWhite),
		},

		BarChart: BarChartTheme{
			Bars:   StandardColors,
			Nums:   []Style{
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
			},
			Labels: []Style{
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
			},
		},

		Gauge: GaugeTheme{
			Bar:   ColorBlue,
			Label: NewStyle(ColorBlack, ColorWhite),
		},

		Plot: PlotTheme{
			Lines: StandardColors,
			Axes:  ColorBlack,
		},

		List: ListTheme{
			Text: NewStyle(ColorBlack, ColorWhite),
		},

		Tree: TreeTheme{
			Text:      NewStyle(ColorBlack, ColorWhite),
			Collapsed: '+',
			Expanded:  '−',
		},

		Paragraph: ParagraphTheme{
			Text: NewStyle(ColorBlack, ColorWhite),
		},

		PieChart: PieChartTheme{
			Slices: StandardColors,
		},

		Sparkline: SparklineTheme{
			Title: NewStyle(ColorBlack, ColorWhite),
			Line:  ColorBlue,
		},

		StackedBarChart: StackedBarChartTheme{
			Bars:   StandardColors,
			Nums:   []Style{
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
			},
			Labels: []Style{
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
				NewStyle(ColorBlack, ColorWhite),
			},
		},

		Table: TableTheme{
			Text: NewStyle(ColorBlack, ColorWhite),
		},

		Tab: TabTheme{
			Active:   NewStyle(ColorRed, ColorWhite, ModifierBold),
			Inactive: NewStyle(ColorBlack, ColorWhite),
		},
	}
}

//
// -----------------------------------------------------
// Default Theme Registration
// -----------------------------------------------------
//

func init() {
	// Create and register dark theme
	darkTheme := createDarkTheme()
	RegisterTheme("dark", darkTheme)

	// Create and register light theme
	lightTheme := createLightTheme()
	RegisterTheme("light", lightTheme)

	// Create default theme (dark mode style)
	defaultTheme := createDarkTheme()
	RegisterTheme("default", defaultTheme)

	// Set default theme as active
	activeTheme = defaultTheme
	currentThemeMode = string(ThemeModeDark) // Default is dark mode
}
