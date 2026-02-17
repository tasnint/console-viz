# Quick Start: Commands to Run Console-Viz

## 🚀 Build the CLI Tool

### Option 1: Build Executable
```bash
# Build the CLI tool
go build -o console-viz.exe ./cmd/console-viz

# Or on Mac/Linux:
go build -o console-viz ./cmd/console-viz
```

### Option 2: Run Directly (No Build)
```bash
# Run directly with go run
go run ./cmd/console-viz/main.go <file-path> [options]
```

---

## 📊 Run with Sample Data

### Using Sample CSV File

```bash
# Basic table view
go run ./cmd/console-viz/main.go data/sample_sales.csv

# Bar chart
go run ./cmd/console-viz/main.go data/sample_sales.csv --widget=barchart

# Line plot
go run ./cmd/console-viz/main.go data/sample_sales.csv --widget=plot

# List view
go run ./cmd/console-viz/main.go data/sample_sales.csv --widget=list
```

**⚠️ IMPORTANT:** Always use `go run ./cmd/console-viz/main.go <data-file>` - don't use `go run <data-file>` directly!

### Using Sample JSON File

```bash
# Table view (JSON)
go run ./cmd/console-viz/main.go data/sample_data.json --widget=table
```

---

## 🎨 With Layout Ratios

```bash
# 80% bar chart, 20% plot
go run ./cmd/console-viz/main.go data/sample_sales.csv \
  --widget="barchart,plot" \
  --layout="80:20"

# 60% table, 40% chart
go run ./cmd/console-viz/main.go data/sample_sales.csv \
  --widget="table,barchart" \
  --layout="60:40"
```

---

## 📁 With Your Own Data Files

### CSV Files

```bash
# From project root
go run ./cmd/console-viz/main.go /path/to/your/data.csv --widget=barchart

# Windows example
go run ./cmd/console-viz/main.go "C:\Users\YourName\Documents\sales.csv" --widget=plot

# Mac/Linux example
go run ./cmd/console-viz/main.go ~/Downloads/data.csv --widget=table

# Your specific file (Windows)
go run ./cmd/console-viz/main.go "C:\Users\tanis\Downloads\MOCK_DATA.json" --widget=sparkline
```

**⚠️ Note:** Use quotes around paths with spaces or special characters!

### JSON Files

```bash
go run ./cmd/console-viz/main.go /path/to/data.json --widget=table
```

---

## 🎯 Complete Examples

### Example 1: Sales Dashboard

```bash
go run ./cmd/console-viz/main.go data/sample_sales.csv \
  --widget="barchart,table" \
  --layout="70:30" \
  --columns="month,sales,profit" \
  --skip-rows=1 \
  --theme=dark \
  --title="Sales Dashboard"
```

### Example 2: Filtered Data

```bash
# First 6 rows only
go run ./cmd/console-viz/main.go data/sample_sales.csv \
  --widget=barchart \
  --rows="1-6" \
  --skip-rows=1
```

### Example 3: Multiple Widgets

```bash
go run ./cmd/console-viz/main.go data/sample_sales.csv \
  --widget="table,barchart,plot" \
  --layout="40:30:30" \
  --theme=dark
```

---

## 🖥️ Run Data Visualizer Example

### Default (uses ./data folder)

```bash
go run examples/data_visualizer.go
```

### Custom Folder Path

```bash
# Relative path
go run examples/data_visualizer.go ./my-data-folder

# Absolute path (Mac/Linux)
go run examples/data_visualizer.go /Users/username/Documents/my-data

# Absolute path (Windows)
go run examples/data_visualizer.go C:\Users\username\Documents\my-data
```

---

## 🎨 Theme Options

```bash
# Dark mode
go run ./cmd/console-viz/main.go data/sample_sales.csv --widget=barchart --theme=dark

# Light mode
go run ./cmd/console-viz/main.go data/sample_sales.csv --widget=barchart --theme=light
```

---

## 📋 All Available Options

```bash
go run ./cmd/console-viz/main.go <file-path> \
  --widget=<type>          # table, barchart, plot, list
  --layout=<ratios>        # "80:20" or "barchart:80,plot:20"
  --columns=<spec>         # "1-3" or "name,value"
  --rows=<spec>           # "1-100" or "-50" (last 50)
  --skip-rows=<n>         # Skip first N rows
  --limit=<n>             # Limit number of rows
  --theme=<mode>          # dark, light, default
  --title=<text>          # Widget title
  --format=<type>         # csv, json, txt (auto-detected)
```

---

## 🔧 Troubleshooting

### If "file not found" error:

```bash
# Check if file exists
ls data/sample_sales.csv        # Mac/Linux
dir data\sample_sales.csv       # Windows

# Use absolute path instead
go run ./cmd/console-viz/main.go $(pwd)/data/sample_sales.csv --widget=barchart
```

### If build fails:

```bash
# Make sure you're in project root
cd "c:\Users\tanis\Downloads\GitHub Repos\console-viz"

# Try building first
go build ./cmd/console-viz
```

---

## ⚡ Quick Test Commands

```bash
# 1. Test with sample CSV
go run ./cmd/console-viz/main.go data/sample_sales.csv --widget=barchart

# 2. Test with layout
go run ./cmd/console-viz/main.go data/sample_sales.csv --widget="barchart,plot" --layout="80:20"

# 3. Test data visualizer
go run examples/data_visualizer.go

# 4. Test with your own file
go run ./cmd/console-viz/main.go /path/to/your/file.csv --widget=table
```

---

## 🎯 Most Common Commands

```bash
# Quick start - see your CSV as a table
go run ./cmd/console-viz/main.go data/sample_sales.csv

# Quick start - see your CSV as a bar chart
go run ./cmd/console-viz/main.go data/sample_sales.csv --widget=barchart

# Quick start - see your CSV as a plot
go run ./cmd/console-viz/main.go data/sample_sales.csv --widget=plot

# Quick start - interactive file browser
go run examples/data_visualizer.go
```

---

**Press ESC to exit when viewing visualizations!**
