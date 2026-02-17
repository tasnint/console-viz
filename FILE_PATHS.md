# File Path Handling in Console-Viz

## ✅ **You Can Use ANY Path on Your Computer!**

Both the CLI tool and examples support **absolute paths**, **relative paths**, and **paths anywhere on your system**.

---

## CLI Tool (`cmd/console-viz/main.go`)

### Supported Path Formats

```bash
# Relative path (from current directory)
console-viz data.csv
console-viz ./data/sales.csv
console-viz ../data/sales.csv

# Absolute path (full path)
console-viz /Users/username/Documents/data.csv          # Mac/Linux
console-viz C:\Users\username\Documents\data.csv        # Windows
console-viz C:/Users/username/Documents/data.csv         # Windows (forward slashes work too)

# Using --file flag
console-viz --file=/path/to/data.csv --widget=barchart
console-viz --file="C:\Users\username\data.csv" --widget=table
```

### How It Works

The CLI uses Go's `os.Open()` which accepts:
- ✅ **Relative paths** - `data.csv`, `./data/file.csv`, `../data/file.csv`
- ✅ **Absolute paths** - `/full/path/to/file.csv` (Unix) or `C:\full\path\to\file.csv` (Windows)
- ✅ **Home directory** - `~/Documents/data.csv` (on Unix systems)
- ✅ **Current directory** - `./file.csv` or just `file.csv`

**Example:**
```bash
# From project root
console-viz data/sample_sales.csv

# From anywhere on your computer
console-viz /Users/john/Documents/sales_data.csv --widget=barchart

# Windows example
console-viz C:\Users\John\Documents\sales_data.csv --widget=plot
```

---

## Data Visualizer Example (`examples/data_visualizer.go`)

### Default Behavior

By default, it looks in the `data/` folder relative to where you run it:

```bash
# Default: looks in ./data/ folder
go run examples/data_visualizer.go
```

### Custom Folder Path

You can specify **any folder path** as an argument:

```bash
# Relative folder path
go run examples/data_visualizer.go ./my-data

# Absolute folder path
go run examples/data_visualizer.go /Users/username/Documents/my-data      # Mac/Linux
go run examples/data_visualizer.go C:\Users\username\Documents\my-data  # Windows

# Parent directory
go run examples/data_visualizer.go ../data-folder
```

**Example:**
```bash
# Visualize files from your Downloads folder
go run examples/data_visualizer.go ~/Downloads

# Visualize files from Windows Documents
go run examples/data_visualizer.go C:\Users\YourName\Documents\DataFiles
```

---

## Path Examples by Operating System

### Mac/Linux

```bash
# Home directory
console-viz ~/Documents/data.csv
console-viz ~/Downloads/sales.csv

# Absolute path
console-viz /Users/john/Documents/data.csv
console-viz /home/john/data/sales.csv

# Relative path
console-viz ./data.csv
console-viz ../data/sales.csv
```

### Windows

```bash
# Absolute path (backslashes)
console-viz C:\Users\John\Documents\data.csv

# Absolute path (forward slashes - also works!)
console-viz C:/Users/John/Documents/data.csv

# Relative path
console-viz .\data.csv
console-viz ..\data\sales.csv

# With quotes (if path has spaces)
console-viz "C:\Users\John\My Documents\data.csv"
```

---

## Current Working Directory

The path is resolved relative to where you run the command:

```bash
# If you're in: /Users/john/projects/console-viz
console-viz data.csv                    # Looks for: /Users/john/projects/console-viz/data.csv
console-viz ./data.csv                  # Same as above
console-viz ../other-project/data.csv   # Looks for: /Users/john/projects/other-project/data.csv
console-viz /absolute/path/data.csv     # Looks for: /absolute/path/data.csv (absolute)
```

---

## Troubleshooting

### "File not found" Error

**Check:**
1. ✅ File path is correct (typos?)
2. ✅ File exists at that location
3. ✅ You have read permissions
4. ✅ Use quotes if path has spaces: `"C:\My Documents\data.csv"`

**Test:**
```bash
# Verify file exists
ls /path/to/file.csv        # Mac/Linux
dir C:\path\to\file.csv     # Windows

# Then use the same path
console-viz /path/to/file.csv
```

### Relative Path Not Working

**Solution:** Use absolute path or check current directory:
```bash
# See where you are
pwd        # Mac/Linux
cd         # Windows

# Use absolute path instead
console-viz /full/path/to/file.csv
```

### Windows Path Issues

**Use forward slashes or quotes:**
```bash
# These all work:
console-viz C:\Users\John\data.csv
console-viz C:/Users/John/data.csv
console-viz "C:\Users\John\My Documents\data.csv"
```

---

## Best Practices

1. **Use absolute paths** for clarity when sharing commands
2. **Use relative paths** when working within a project
3. **Quote paths** that contain spaces
4. **Test with `ls`/`dir`** first to verify file exists

---

## Examples

### Example 1: File in Downloads

```bash
# Mac/Linux
console-viz ~/Downloads/sales_data.csv --widget=barchart

# Windows
console-viz C:\Users\YourName\Downloads\sales_data.csv --widget=barchart
```

### Example 2: File in Different Project

```bash
# From console-viz directory, access file in another project
console-viz ../other-project/data/analytics.csv --widget=plot
```

### Example 3: Network/External Drive

```bash
# Mac/Linux
console-viz /Volumes/ExternalDrive/data.csv

# Windows
console-viz D:\DataFiles\sales.csv
```

---

## Summary

✅ **YES** - You can specify **any path** on your local computer:
- Relative paths: `data.csv`, `./folder/file.csv`
- Absolute paths: `/full/path/to/file.csv` or `C:\full\path\to\file.csv`
- Home directory: `~/Documents/file.csv` (Unix)
- Any folder: `examples/data_visualizer.go /any/folder/path`

The tools use standard Go file operations that work with any valid file system path!
