# Console-Viz CLI Usage Examples

## Basic Usage

### Single Widget

```bash
# Table view (default) - relative path
console-viz data.csv

# Bar chart - absolute path (Mac/Linux)
console-viz /Users/username/Documents/data.csv --widget=barchart

# Line plot - absolute path (Windows)
console-viz C:\Users\username\Documents\data.csv --widget=plot

# List view - any path on your computer
console-viz ~/Downloads/sales.csv --widget=list

# Using --file flag
console-viz --file=/path/to/data.csv --widget=barchart
```

**Note:** You can use **any file path** - relative, absolute, or anywhere on your computer!

---

## Layout Ratios

### Two Widgets with Ratios

```bash
# 80% bar chart, 20% plot
console-viz data.csv --widget="barchart,plot" --layout="80:20"

# 60% table, 40% chart
console-viz data.csv --widget="table,barchart" --layout="60:40"

# Named ratios
console-viz data.csv --widget="barchart,plot" --layout="barchart:80,plot:20"
```

### Three Widgets

```bash
# Equal distribution (33% each)
console-viz data.csv --widget="table,barchart,plot" --layout="33:33:34"

# Custom ratios
console-viz data.csv --widget="table,barchart,plot" --layout="40:35:25"
```

---

## Data Selection

### Column Selection

```bash
# Specific columns by index
console-viz data.csv --columns="1-3"

# Specific columns by name
console-viz data.csv --columns="name,value,date"

# Mixed selection
console-viz data.csv --columns="1,3,5,name"
```

### Row Filtering

```bash
# First 100 rows
console-viz data.csv --rows="1-100"

# Rows 10-50
console-viz data.csv --rows="10-50"

# Last 50 rows
console-viz data.csv --rows="-50"

# Skip header row
console-viz data.csv --skip-rows=1

# Limit to 50 rows
console-viz data.csv --limit=50
```

### Combined Filtering

```bash
# First 100 rows, columns 1-3, skip header
console-viz data.csv --rows="1-100" --columns="1-3" --skip-rows=1 --widget=barchart
```

---

## File Formats

### CSV Files

```bash
console-viz sales.csv --widget=barchart
console-viz data.csv --widget=table --columns="month,sales"
```

### JSON Files

```bash
console-viz data.json --widget=table
console-viz metrics.json --widget=plot
```

### Auto-Detection

The tool automatically detects file format by extension:
- `.csv` → CSV format
- `.json` → JSON format
- `.txt` → Text format

Force format:
```bash
console-viz data.txt --format=csv
```

---

## Themes

```bash
# Dark mode
console-viz data.csv --theme=dark

# Light mode
console-viz data.csv --theme=light

# Default theme
console-viz data.csv --theme=default
```

---

## Titles

```bash
# Custom widget title
console-viz data.csv --widget=barchart --title="Sales by Month"

# Multiple widgets with titles (via config file)
console-viz --config=my-config.yaml
```

---

## Complete Examples

### Sales Dashboard

```bash
# Sales data: 70% bar chart, 30% table
console-viz sales.csv \
  --widget="barchart,table" \
  --layout="70:30" \
  --columns="month,sales,profit" \
  --skip-rows=1 \
  --theme=dark \
  --title="Sales Dashboard"
```

### Time Series Analysis

```bash
# Time series: 60% plot, 40% sparkline
console-viz timeseries.csv \
  --widget="plot,sparkline" \
  --layout="60:40" \
  --rows="1-1000" \
  --columns="timestamp,value"
```

### Multi-Widget Dashboard

```bash
# Three widgets: table, bar chart, plot
console-viz data.csv \
  --widget="table,barchart,plot" \
  --layout="40:30:30" \
  --rows="1-50" \
  --theme=dark
```

---

## Advanced: Config Files

Create `dashboard.yaml`:

```yaml
data:
  source: "sales.csv"
  skip_rows: 1
  columns: ["month", "sales", "profit"]
  rows: "1-12"

layout:
  widgets:
    - type: "barchart"
      ratio: 0.7
      title: "Sales by Month"
    - type: "table"
      ratio: 0.3
      title: "Data Table"

theme: "dark"
```

Run:
```bash
console-viz --config=dashboard.yaml
```

---

## Tips

1. **Start Simple**: Begin with `--widget=table` to see your data
2. **Use Ratios**: `--layout="80:20"` is easier than pixel calculations
3. **Filter First**: Use `--rows` and `--columns` to focus on relevant data
4. **Theme**: Use `--theme=dark` for better visibility in low light
5. **Multiple Widgets**: Combine different views for comprehensive analysis

---

## Keyboard Shortcuts

- **ESC** - Exit
- **Arrow Keys** - Navigate (in list/table widgets)
- **q** - Quit (alternative)

---

## Troubleshooting

**No data displayed?**
- Check file path is correct
- Verify file format (CSV/JSON)
- Check column/row selections are valid

**Widgets overlap?**
- Adjust `--layout` ratios (should sum to 100%)
- Use fewer widgets if terminal is small

**Wrong data shown?**
- Check `--skip-rows` for headers
- Verify `--columns` selection
- Check `--rows` range
