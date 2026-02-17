# CLI Data Visualization Design

## Research Summary: Standard Approaches

Based on research of tools like **VisiData**, **csvkit**, **Gurita**, **WTF**, and **vcli**, here are standard patterns:

### 1. **Command-Line Flags** (Most Common)
- Widget selection via flags: `--widget=barchart`, `--widget=table`
- Layout ratios: `--layout="80:20"` or `--layout="barchart:80,plot:20"`
- Data selection: `--columns=1-3`, `--rows=1-100`

### 2. **Grid/Ratio Layouts** (WTF Pattern)
- Grid-based positioning: `[10, 10, 10, 10]` (columns)
- Ratio-based sizing: `0.8` (80%), `0.2` (20%)
- Widget positioning: `top`, `left`, `width`, `height`

### 3. **Chainable Commands** (Gurita Pattern)
- Pipe operations: `load data.csv | barchart --columns=1,2 | plot --columns=3`
- Multiple visualizations: `--viz="barchart:80,plot:20"`

### 4. **Config Files** (Complex Layouts)
- YAML/JSON configs for complex dashboards
- Reusable layouts

---

## Proposed Design: `console-viz` CLI

### Command Structure

```bash
# Basic usage
console-viz <data-file> [options]

# Examples
console-viz data.csv --widget=barchart
console-viz data.json --widget=table --columns=name,value
console-viz data.csv --layout="barchart:80,plot:20"
console-viz data.csv --widget=barchart --rows=1-50
```

---

## CLI Flags Design

### Widget Selection

```bash
--widget=<type>          # Widget type: table, barchart, plot, list, gauge, etc.
--widgets=<list>        # Multiple widgets: "barchart,plot,table"
```

**Widget Types:**
- `table` - Tabular display
- `barchart` - Bar chart
- `plot` - Line/scatter plot
- `sparkline` - Mini chart
- `piechart` - Pie chart
- `list` - Scrollable list
- `gauge` - Progress gauge
- `tree` - Hierarchical tree

### Layout Ratios

```bash
--layout=<ratio>        # Single widget: "80" or "0.8" (80% of screen)
--layout=<spec>         # Multiple: "barchart:80,plot:20"
--layout=<grid>         # Grid: "2x2" (2 columns, 2 rows)
```

**Layout Formats:**
1. **Percentage**: `--layout="80:20"` (80% widget1, 20% widget2)
2. **Named Ratios**: `--layout="barchart:80,plot:20"`
3. **Grid**: `--layout="2x2"` (equal grid)
4. **Custom Grid**: `--layout="[0.6,0.4],[0.5,0.5]"` (rows x columns)

### Data Selection

```bash
--columns=<spec>        # Column selection: "1-3", "name,value", "1,3,5"
--rows=<spec>           # Row selection: "1-100", "10-50", "all"
--skip-rows=<n>         # Skip first N rows (headers)
--limit=<n>             # Limit rows displayed
```

### Data Source

```bash
--file=<path>           # File path (CSV, JSON, TXT)
--db=<connection>       # Database connection string
--format=<type>         # Force format: csv, json, txt
```

### Advanced Options

```bash
--x-axis=<column>       # X-axis column for charts
--y-axis=<column>       # Y-axis column for charts
--labels=<column>       # Labels column
--title=<text>          # Widget title
--theme=<mode>          # dark, light, default
--config=<file>         # Config file path
```

---

## Example Commands

### Simple Single Widget

```bash
# Bar chart from CSV
console-viz sales.csv --widget=barchart

# Table with specific columns
console-viz data.csv --widget=table --columns=name,value,date

# Plot with row limit
console-viz timeseries.csv --widget=plot --rows=1-100
```

### Ratio-Based Layouts

```bash
# 80% bar chart, 20% plot
console-viz data.csv --layout="barchart:80,plot:20"

# Three widgets with ratios
console-viz data.csv --layout="table:40,barchart:35,plot:25"

# Grid layout
console-viz data.csv --layout="2x2" --widgets="table,barchart,plot,sparkline"
```

### Data Filtering

```bash
# First 50 rows, columns 1-3
console-viz data.csv --widget=table --rows=1-50 --columns=1-3

# Skip header, specific columns
console-viz data.csv --widget=barchart --skip-rows=1 --columns=name,value

# Last 100 rows
console-viz data.csv --widget=plot --rows=-100
```

### Multiple Data Sources

```bash
# Multiple files
console-viz file1.csv file2.json --widget=table

# Database
console-viz --db="sqlite:data.db" --query="SELECT * FROM sales" --widget=barchart
```

---

## Config File Format (YAML)

For complex layouts, use config files:

```yaml
# viz-config.yaml
data:
  source: "sales.csv"
  format: "csv"
  skip_rows: 1
  columns: ["month", "sales", "profit"]
  rows: "1-12"

layout:
  type: "ratio"
  widgets:
    - type: "barchart"
      ratio: 0.6
      title: "Sales by Month"
      x_axis: "month"
      y_axis: "sales"
    - type: "plot"
      ratio: 0.4
      title: "Profit Trend"
      x_axis: "month"
      y_axis: "profit"

theme: "dark"
```

Usage:
```bash
console-viz --config=viz-config.yaml
```

---

## Implementation Plan

### Phase 1: Basic CLI (MVP)
1. ✅ File detection (CSV, JSON, TXT)
2. ✅ Single widget selection (`--widget`)
3. ✅ Basic data loading
4. ✅ Simple rendering

### Phase 2: Layout System
1. ✅ Ratio-based layouts (`--layout="80:20"`)
2. ✅ Named widget ratios (`--layout="barchart:80,plot:20"`)
3. ✅ Grid layouts (`--layout="2x2"`)

### Phase 3: Data Selection
1. ✅ Column selection (`--columns`)
2. ✅ Row filtering (`--rows`, `--limit`)
3. ✅ Data transformation

### Phase 4: Advanced Features
1. ✅ Config files (YAML/JSON)
2. ✅ Database support
3. ✅ Multiple data sources
4. ✅ Interactive mode

---

## Standard Patterns Used

### 1. **Flag-Based Selection** (csvkit pattern)
```bash
--widget=barchart        # Single selection
--widgets="a,b,c"       # Multiple selection
```

### 2. **Ratio Layouts** (WTF pattern)
```bash
--layout="80:20"        # Percentage ratios
--layout="0.8:0.2"      # Decimal ratios
```

### 3. **Data Selection** (VisiData pattern)
```bash
--columns=1-3          # Range
--columns=name,value    # Named
--rows=1-100           # Range
```

### 4. **Config Files** (Standard)
```bash
--config=file.yaml     # Complex layouts
```

---

## Comparison with Standard Tools

| Feature | VisiData | csvkit | WTF | **console-viz** |
|---------|----------|--------|-----|-----------------|
| Interactive | ✅ | ❌ | ✅ | ✅ |
| CLI Flags | ❌ | ✅ | ✅ | ✅ |
| Ratio Layouts | ❌ | ❌ | ✅ | ✅ |
| Widget Selection | ❌ | ❌ | ✅ | ✅ |
| Config Files | ✅ | ❌ | ✅ | ✅ |
| Multiple Formats | ✅ | CSV only | ❌ | ✅ |

**console-viz combines:**
- ✅ CLI flexibility (csvkit)
- ✅ Ratio layouts (WTF)
- ✅ Multiple formats (VisiData)
- ✅ Interactive terminal UI (all)

---

## Next Steps

1. **Implement CLI parser** (`cmd/console-viz/main.go`)
2. **Add layout ratio parser** (parse `"80:20"` → ratios)
3. **Implement data loaders** (CSV, JSON, DB)
4. **Create config file parser** (YAML/JSON)
5. **Build interactive mode** (widget selection menu)

This design follows industry standards while adding unique features like ratio-based layouts and multi-format support!
