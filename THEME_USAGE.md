# Theme Switching Guide

## Overview

The console-viz library now supports dark mode and light mode themes that can be switched via command-line, environment variables, or config files - **no code changes required!**

## Color Schemes

### Dark Mode
- **Background**: Black (`ColorBlack`)
- **Text**: White (`ColorWhite`)
- **Accents**: Cyan for gauges, Red for active tabs
- **Best for**: Low-light environments, reducing eye strain

### Light Mode
- **Background**: White (`ColorWhite`)
- **Text**: Black (`ColorBlack`)
- **Accents**: Blue for gauges, Red for active tabs
- **Best for**: Bright environments, high contrast

## Usage Methods

### 1. Command-Line Flag (Highest Priority)

```bash
# Set dark mode
./your-app --theme=dark

# Set light mode
./your-app --theme=light

# Set default theme
./your-app --theme=default
```

### 2. Environment Variable

```bash
# Linux/Mac
export CONSOLE_VIZ_THEME=dark
./your-app

# Windows (PowerShell)
$env:CONSOLE_VIZ_THEME="dark"
./your-app

# Windows (CMD)
set CONSOLE_VIZ_THEME=dark
./your-app
```

### 3. Config File (Persistent)

The theme preference is automatically saved to:
- **Linux/Mac**: `~/.config/console-viz/theme.json`
- **Windows**: `%APPDATA%\console-viz\theme.json`
- **Fallback**: `./theme.json` (current directory)

The config file format:
```json
{
  "theme": "dark"
}
```

Once set via CLI or environment variable, the preference is saved automatically.

### 4. Programmatic (Code)

```go
import "console-viz/styling"

// Switch to dark mode
styling.ToggleDarkMode()

// Switch to light mode
styling.ToggleLightMode()

// Toggle between dark/light
styling.ToggleThemeMode()

// Set specific theme
styling.SetThemeFromString("dark")

// Check current mode
mode := styling.GetCurrentThemeMode()
```

## Priority Order

Theme selection follows this priority (highest to lowest):
1. **CLI flag** (`--theme=dark`)
2. **Environment variable** (`CONSOLE_VIZ_THEME=dark`)
3. **Config file** (`~/.config/console-viz/theme.json`)
4. **Default** (dark mode)

## Example Application

See `examples/theme_demo.go` for a complete example:

```bash
# Run with dark theme
go run examples/theme_demo.go --theme=dark

# Run with light theme
go run examples/theme_demo.go --theme=light

# Press 'T' in the app to toggle themes interactively
```

## Integration in Your Code

Add this to your `main()` function:

```go
package main

import (
    "console-viz/styling"
    // ... other imports
)

func main() {
    // Parse and initialize theme from CLI/env/config
    cliTheme := styling.ParseThemeFlag()
    if err := styling.InitThemeFromCLI(cliTheme); err != nil {
        log.Printf("Warning: Theme initialization failed: %v", err)
    }
    
    // Your application code...
}
```

## Available Themes

List all available themes:
```go
themes := styling.ListThemes()
// Returns: ["dark", "light", "default"]
```

Check if a theme exists:
```go
if styling.HasTheme("dark") {
    // Theme exists
}
```

## Theme Change Callbacks

Register callbacks to be notified when theme changes:

```go
styling.RegisterThemeChangeCallback(func(mode styling.ThemeMode) {
    fmt.Printf("Theme changed to: %s\n", mode)
    // Trigger re-render, update UI, etc.
})
```

## Troubleshooting

**Theme not applying?**
- Check priority order (CLI > Env > Config > Default)
- Verify theme name is correct: `dark`, `light`, or `default`
- Check config file permissions

**Config file location?**
- Run your app once with `--theme=dark` to create the config file
- Check `~/.config/console-viz/theme.json` (Linux/Mac)
- Check `%APPDATA%\console-viz\theme.json` (Windows)

**Invalid theme error?**
- Use `styling.ListThemes()` to see available themes
- Theme names are case-insensitive but must match exactly
