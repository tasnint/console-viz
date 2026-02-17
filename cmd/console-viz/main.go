package main

import (
	"console-viz/draw"
	"console-viz/styling"
	"console-viz/widgets"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config holds CLI configuration
type Config struct {
	DataFile   string
	Layout     string
	Columns    string
	Rows       string
	SkipRows   int
	Limit      int
	Theme      string
	XAxis      string
	YAxis      string
	Title      string
	Format     string
	ConfigFile string
}

// parseLayout parses layout string like "80:20" or "barchart:80,plot:20"
func parseLayout(layoutStr string, widgetCount int) ([]float64, error) {
	if layoutStr == "" {
		// Default: equal distribution
		ratio := 1.0 / float64(widgetCount)
		ratios := make([]float64, widgetCount)
		for i := range ratios {
			ratios[i] = ratio
		}
		return ratios, nil
	}

	// Check if it's a simple ratio like "80:20"
	if strings.Contains(layoutStr, ":") && !strings.Contains(layoutStr, ",") {
		parts := strings.Split(layoutStr, ":")
		ratios := make([]float64, len(parts))
		total := 0.0
		for i, part := range parts {
			val, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
			if err != nil {
				return nil, fmt.Errorf("invalid ratio: %s", part)
			}
			ratios[i] = val
			total += val
		}
		// Normalize to 0-1 range
		for i := range ratios {
			ratios[i] /= total
		}
		return ratios, nil
	}

	// Check if it's named ratios like "barchart:80,plot:20"
	if strings.Contains(layoutStr, ",") {
		parts := strings.Split(layoutStr, ",")
		ratios := make([]float64, len(parts))
		total := 0.0
		for i, part := range parts {
			part = strings.TrimSpace(part)
			if strings.Contains(part, ":") {
				subParts := strings.Split(part, ":")
				if len(subParts) != 2 {
					return nil, fmt.Errorf("invalid named ratio: %s", part)
				}
				val, err := strconv.ParseFloat(strings.TrimSpace(subParts[1]), 64)
				if err != nil {
					return nil, fmt.Errorf("invalid ratio value: %s", subParts[1])
				}
				ratios[i] = val
				total += val
			} else {
				val, err := strconv.ParseFloat(part, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid ratio: %s", part)
				}
				ratios[i] = val
				total += val
			}
		}
		// Normalize
		for i := range ratios {
			ratios[i] /= total
		}
		return ratios, nil
	}

	// Single value - distribute equally
	val, err := strconv.ParseFloat(layoutStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid layout: %s", layoutStr)
	}
	if val > 1.0 {
		val /= 100.0 // Assume percentage
	}
	ratio := val / float64(widgetCount)
	ratios := make([]float64, widgetCount)
	for i := range ratios {
		ratios[i] = ratio
	}
	return ratios, nil
}

// parseColumnSpec parses column specification like "1-3" or "name,value" or "1,3,5"
func parseColumnSpec(spec string, headers []string) ([]int, error) {
	if spec == "" {
		return nil, nil // All columns
	}

	// Check if it's a range like "1-3"
	if strings.Contains(spec, "-") {
		parts := strings.Split(spec, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid column range: %s", spec)
		}
		start, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		end, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid column range: %s", spec)
		}
		cols := []int{}
		for i := start; i <= end; i++ {
			cols = append(cols, i-1) // Convert to 0-based
		}
		return cols, nil
	}

	// Check if it's comma-separated
	if strings.Contains(spec, ",") {
		parts := strings.Split(spec, ",")
		cols := []int{}
		for _, part := range parts {
			part = strings.TrimSpace(part)
			// Try as number first
			if idx, err := strconv.Atoi(part); err == nil {
				cols = append(cols, idx-1) // Convert to 0-based
			} else {
				// Try as column name
				found := false
				for i, h := range headers {
					if h == part {
						cols = append(cols, i)
						found = true
						break
					}
				}
				if !found {
					return nil, fmt.Errorf("column not found: %s", part)
				}
			}
		}
		return cols, nil
	}

	// Single value
	if idx, err := strconv.Atoi(spec); err == nil {
		return []int{idx - 1}, nil // Convert to 0-based
	}
	// Try as column name
	for i, h := range headers {
		if h == spec {
			return []int{i}, nil
		}
	}
	return nil, fmt.Errorf("column not found: %s", spec)
}

// parseRowSpec parses row specification like "1-100" or "10-50"
func parseRowSpec(spec string, totalRows int) (start, end int, err error) {
	if spec == "" {
		return 0, totalRows, nil // All rows
	}

	if strings.HasPrefix(spec, "-") {
		// Negative: last N rows
		n, err := strconv.Atoi(spec[1:])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid row spec: %s", spec)
		}
		start := totalRows - n
		if start < 0 {
			start = 0
		}
		return start, totalRows, nil
	}

	if strings.Contains(spec, "-") {
		parts := strings.Split(spec, "-")
		if len(parts) != 2 {
			return 0, 0, fmt.Errorf("invalid row range: %s", spec)
		}
		start, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		end, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err1 != nil || err2 != nil {
			return 0, 0, fmt.Errorf("invalid row range: %s", spec)
		}
		// Convert to 0-based
		return start - 1, end, nil
	}

	// Single value
	val, err := strconv.Atoi(spec)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid row spec: %s", spec)
	}
	return val - 1, val, nil
}

// loadCSV loads a CSV file
func loadCSV(filePath string, config Config) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Skip rows if specified
	start := config.SkipRows
	if start >= len(records) {
		return [][]string{}, nil
	}
	records = records[start:]

	// Apply row filter
	if config.Rows != "" {
		rowStart, rowEnd, err := parseRowSpec(config.Rows, len(records))
		if err != nil {
			return nil, err
		}
		if rowStart < len(records) {
			if rowEnd > len(records) {
				rowEnd = len(records)
			}
			records = records[rowStart:rowEnd]
		}
	}

	// Apply limit
	if config.Limit > 0 && config.Limit < len(records) {
		records = records[:config.Limit]
	}

	// Apply column filter
	if config.Columns != "" && len(records) > 0 {
		headers := records[0]
		colIndices, err := parseColumnSpec(config.Columns, headers)
		if err != nil {
			return nil, err
		}
		if len(colIndices) > 0 {
			filtered := [][]string{}
			for _, record := range records {
				filteredRow := []string{}
				for _, idx := range colIndices {
					if idx < len(record) {
						filteredRow = append(filteredRow, record[idx])
					}
				}
				filtered = append(filtered, filteredRow)
			}
			records = filtered
		}
	}

	return records, nil
}

// loadJSON loads a JSON file
func loadJSON(filePath string) (interface{}, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, err
	}

	return jsonData, nil
}

// numericFieldNames are JSON keys we try when extracting numbers from array-of-objects
var numericFieldNames = []string{"id", "value", "count", "amount", "number", "score", "sales", "price", "quantity", "total"}

// extractNumericArray extracts []float64 from JSON data
func extractNumericArray(jsonData interface{}) []float64 {
	values := []float64{}

	// Array of numbers: [1, 2, 3]
	if arr, ok := jsonData.([]interface{}); ok {
		allNumeric := true
		for _, v := range arr {
			switch val := v.(type) {
			case float64:
				values = append(values, val)
			case int:
				values = append(values, float64(val))
			case string:
				if num, err := strconv.ParseFloat(val, 64); err == nil {
					values = append(values, num)
				} else {
					allNumeric = false
				}
			default:
				allNumeric = false
			}
		}
		if allNumeric && len(values) > 0 {
			return values
		}
		// Array of objects: [{id: 1, name: "..."}, {id: 2, ...}] - extract numeric column
		if len(values) == 0 {
			for _, item := range arr {
				if obj, ok := item.(map[string]interface{}); ok {
					for _, fieldName := range numericFieldNames {
						if v, exists := obj[fieldName]; exists {
							switch val := v.(type) {
							case float64:
								values = append(values, val)
								break
							case int:
								values = append(values, float64(val))
								break
							case string:
								if num, err := strconv.ParseFloat(val, 64); err == nil {
									values = append(values, num)
									break
								}
							}
							break
						}
					}
				}
			}
		}
		if len(values) > 0 {
			return values
		}
	}

	// Nested structure: {"data": [1,2,3]} or {"time_series": [1,2,3]}
	if obj, ok := jsonData.(map[string]interface{}); ok {
		fieldNames := []string{"data", "values", "time_series", "series", "array"}
		for _, field := range fieldNames {
			if arr, exists := obj[field]; exists {
				if nums := extractNumericArray(arr); len(nums) > 0 {
					return nums
				}
			}
		}
		for _, v := range obj {
			if nums := extractNumericArray(v); len(nums) > 0 {
				return nums
			}
		}
	}

	return values
}

// jsonToTable converts JSON to table format
func jsonToTable(jsonData interface{}) [][]string {
	rows := [][]string{}
	
	if arr, ok := jsonData.([]interface{}); ok {
		// Array of objects
		for i, item := range arr {
			if obj, ok := item.(map[string]interface{}); ok {
				if i == 0 {
					// Header row
					header := []string{}
					for k := range obj {
						header = append(header, k)
					}
					rows = append(rows, header)
				}
				// Data row
				row := []string{}
				for _, k := range rows[0] {
					val := fmt.Sprintf("%v", obj[k])
					row = append(row, val)
				}
				rows = append(rows, row)
			}
		}
	} else if obj, ok := jsonData.(map[string]interface{}); ok {
		// Single object - create key-value table
		rows = append(rows, []string{"Key", "Value"})
		for k, v := range obj {
			rows = append(rows, []string{k, fmt.Sprintf("%v", v)})
		}
	}
	
	return rows
}

// createWidget creates a widget based on type and data
func createWidget(widgetType string, data interface{}, config Config) (draw.Drawable, error) {
	switch widgetType {
	case "table":
		var tableRows [][]string
		if csvData, ok := data.([][]string); ok {
			tableRows = csvData
		} else {
			// Try to convert JSON to table
			tableRows = jsonToTable(data)
			if len(tableRows) == 0 {
				return nil, fmt.Errorf("table widget: unable to convert data to table format")
			}
		}
		table := widgets.NewTable()
		table.Rows = tableRows
		if config.Title != "" {
			table.Title = config.Title
		}
		return table, nil

	case "barchart":
		var values []float64
		var labels []string
		
		if csvData, ok := data.([][]string); ok && len(csvData) > 1 {
			// CSV data
			for i, row := range csvData {
				if i == 0 {
					continue // Skip header
				}
				if len(row) > 0 {
					labels = append(labels, row[0])
				}
				if len(row) > 1 {
					if val, err := strconv.ParseFloat(row[1], 64); err == nil {
						values = append(values, val)
					}
				}
			}
		} else {
			// JSON data
			values = extractNumericArray(data)
			if len(values) == 0 {
				return nil, fmt.Errorf("barchart widget: no numeric data found in JSON")
			}
			// Generate labels
			for i := range values {
				labels = append(labels, fmt.Sprintf("Item %d", i+1))
			}
		}
		
		if len(values) == 0 {
			return nil, fmt.Errorf("barchart widget: no numeric data found")
		}
		
		chart := widgets.NewBarChart()
		chart.Data = values
		chart.Labels = labels
		if config.Title != "" {
			chart.Title = config.Title
		}
		return chart, nil

	case "plot":
		var values []float64
		
		if csvData, ok := data.([][]string); ok && len(csvData) > 1 {
			// CSV data
			for i, row := range csvData {
				if i == 0 {
					continue // Skip header
				}
				if len(row) > 1 {
					if val, err := strconv.ParseFloat(row[1], 64); err == nil {
						values = append(values, val)
					}
				}
			}
		} else {
			// JSON data
			values = extractNumericArray(data)
		}
		
		if len(values) == 0 {
			return nil, fmt.Errorf("plot widget: no numeric data found")
		}
		
		plot := widgets.NewPlot()
		plot.Data = [][]float64{values}
		if config.Title != "" {
			plot.Title = config.Title
		}
		return plot, nil

	case "sparkline":
		var values []float64
		
		if csvData, ok := data.([][]string); ok && len(csvData) > 1 {
			// CSV data
			for i, row := range csvData {
				if i == 0 {
					continue // Skip header
				}
				if len(row) > 1 {
					if val, err := strconv.ParseFloat(row[1], 64); err == nil {
						values = append(values, val)
					}
				}
			}
		} else {
			// JSON data
			values = extractNumericArray(data)
		}
		
		if len(values) == 0 {
			return nil, fmt.Errorf("sparkline widget: no numeric data found")
		}
		
		sparkline := widgets.NewSparkline()
		sparkline.Data = values
		if config.Title != "" {
			sparkline.Title = config.Title
		}
		sparklineGroup := widgets.NewSparklineGroup(sparkline)
		return sparklineGroup, nil

	case "horizontal", "hbar", "horizontal-barchart":
		// Horizontal bar chart - perfect for distributions
		var values []float64
		var labels []string
		
		if csvData, ok := data.([][]string); ok && len(csvData) > 1 {
			// CSV data: first column = labels, second column = values
			for i, row := range csvData {
				if i == 0 {
					continue // Skip header
				}
				if len(row) > 0 {
					labels = append(labels, row[0])
				}
				if len(row) > 1 {
					if val, err := strconv.ParseFloat(row[1], 64); err == nil {
						values = append(values, val)
					} else {
						// If second column isn't numeric, try to extract from JSON structure
						values = append(values, 0)
					}
				}
			}
		} else {
			// JSON data - try to extract category/value pairs
			values = extractNumericArray(data)
			// Try to extract labels from JSON
			if arr, ok := data.([]interface{}); ok {
				for _, item := range arr {
					if obj, ok := item.(map[string]interface{}); ok {
						// Try to find a label field
						for _, labelField := range []string{"name", "label", "category", "key", "type"} {
							if v, exists := obj[labelField]; exists {
								labels = append(labels, fmt.Sprintf("%v", v))
								break
							}
						}
						// If no label field, use first string field
						if len(labels) < len(values) {
							for k, v := range obj {
								if _, ok := v.(string); ok && k != "id" {
									labels = append(labels, fmt.Sprintf("%v", v))
									break
								}
							}
						}
					}
				}
			}
			// Generate labels if we have values but no labels
			if len(labels) == 0 && len(values) > 0 {
				for i := range values {
					labels = append(labels, fmt.Sprintf("Item %d", i+1))
				}
			}
		}
		
		if len(values) == 0 {
			return nil, fmt.Errorf("horizontal-barchart widget: no numeric data found")
		}
		
		hbc := widgets.NewHorizontalBarChart()
		hbc.Data = values
		hbc.Labels = labels
		if config.Title != "" {
			hbc.Title = config.Title
		}
		return hbc, nil

	case "list":
		var rows []string
		
		if csvData, ok := data.([][]string); ok {
			// CSV data
			for _, row := range csvData {
				rows = append(rows, strings.Join(row, " | "))
			}
		} else {
			// JSON data - convert to string representation
			if arr, ok := data.([]interface{}); ok {
				for _, item := range arr {
					rows = append(rows, fmt.Sprintf("%v", item))
				}
			} else if obj, ok := data.(map[string]interface{}); ok {
				for k, v := range obj {
					rows = append(rows, fmt.Sprintf("%s: %v", k, v))
				}
			} else {
				rows = append(rows, fmt.Sprintf("%v", data))
			}
		}
		
		list := widgets.NewList()
		list.Rows = rows
		if config.Title != "" {
			list.Title = config.Title
		}
		return list, nil

	default:
		return nil, fmt.Errorf("unknown widget type: %s", widgetType)
	}
}

func main() {
	config := Config{}

	// Parse flags
	var widgetStr string
	flag.StringVar(&config.DataFile, "file", "", "Data file path (CSV, JSON, TXT)")
	flag.StringVar(&widgetStr, "widget", "table", "Widget type: table, barchart, horizontal, horizontal-barchart, plot, sparkline, list (comma-separated for multiple)")
	flag.StringVar(&config.Layout, "layout", "", "Layout ratios: '80:20' or 'barchart:80,plot:20'")
	flag.StringVar(&config.Columns, "columns", "", "Column selection: '1-3' or 'name,value'")
	flag.StringVar(&config.Rows, "rows", "", "Row selection: '1-100' or '-50' (last 50)")
	flag.IntVar(&config.SkipRows, "skip-rows", 0, "Skip first N rows")
	flag.IntVar(&config.Limit, "limit", 0, "Limit number of rows")
	flag.StringVar(&config.Theme, "theme", "", "Theme: dark, light, default")
	flag.StringVar(&config.Title, "title", "", "Widget title")
	flag.StringVar(&config.Format, "format", "", "Force format: csv, json, txt")
	flag.Parse()

	// Get data file from args if not in flag
	if config.DataFile == "" && len(flag.Args()) > 0 {
		config.DataFile = flag.Args()[0]
	}

	if config.DataFile == "" {
		fmt.Fprintf(os.Stderr, "Usage: console-viz <data-file> [options]\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Initialize theme
	// Check if theme was set via flag first
	if config.Theme != "" {
		if err := styling.InitThemeFromCLI(config.Theme); err != nil {
			log.Printf("Warning: Theme initialization failed: %v", err)
		}
	} else {
		// Try to get theme from environment variable or config file
		// Don't call ParseThemeFlag() as it will conflict with our flag definition
		if err := styling.InitThemeFromCLI(""); err != nil {
			log.Printf("Warning: Theme initialization failed: %v", err)
		}
	}

	// Initialize terminal
	if err := draw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to initialize terminal: %v\n", err)
		fmt.Fprintf(os.Stderr, "Make sure you're running in a real terminal (not piping output)\n")
		fmt.Fprintf(os.Stderr, "Try running: go run ./cmd/console-viz/main.go <file> [options]\n")
		os.Exit(1)
	}
	defer draw.Close()

	draw.InitRenderer()

	// Detect file format
	ext := strings.ToLower(filepath.Ext(config.DataFile))
	if config.Format == "" {
		config.Format = ext[1:] // Remove dot
	}

	// Load data
	var data interface{}
	var err error

	switch config.Format {
	case "csv":
		data, err = loadCSV(config.DataFile, config)
		if err != nil {
			log.Fatalf("Failed to load CSV: %v", err)
		}
	case "json":
		data, err = loadJSON(config.DataFile)
		if err != nil {
			log.Fatalf("Failed to load JSON: %v", err)
		}
	default:
		log.Fatalf("Unsupported format: %s", config.Format)
	}

	// Parse widgets (normalize to lowercase so "Horizontal-Barchart" matches)
	widgetTypes := []string{}
	if widgetStr != "" {
		widgetTypes = strings.Split(widgetStr, ",")
		for i := range widgetTypes {
			widgetTypes[i] = strings.TrimSpace(strings.ToLower(widgetTypes[i]))
		}
	} else {
		widgetTypes = []string{"table"}
	}

	// Debug: Print what widget was requested
	if len(os.Getenv("DEBUG")) > 0 {
		log.Printf("DEBUG: Requested widget(s): %v", widgetTypes)
		log.Printf("DEBUG: widgetStr value: '%s'", widgetStr)
	}

	// Create widgets
	widgetList := []draw.Drawable{}
	for _, widgetType := range widgetTypes {
		widget, err := createWidget(widgetType, data, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to create widget '%s': %v\n", widgetType, err)
			fmt.Fprintf(os.Stderr, "Available widgets: table, barchart, horizontal, horizontal-barchart, plot, sparkline, list\n")
			fmt.Fprintf(os.Stderr, "\nTip: Make sure your data contains numeric values for chart widgets.\n")
			os.Exit(1)
		}
		widgetList = append(widgetList, widget)
	}
	
	if len(widgetList) == 0 {
		log.Fatalf("Error: No widgets were created. Check your data format and widget type.")
	}

	// Parse layout ratios
	ratios, err := parseLayout(config.Layout, len(widgetList))
	if err != nil {
		log.Printf("Warning: Failed to parse layout: %v, using equal distribution", err)
		ratios = make([]float64, len(widgetList))
		for i := range ratios {
			ratios[i] = 1.0 / float64(len(widgetList))
		}
	}

	// Get terminal dimensions
	width, height := draw.TerminalDimensions()

	// Position widgets based on ratios
	currentX := 0
	for i, widget := range widgetList {
		widgetWidth := int(float64(width) * ratios[i])
		widget.SetRect(currentX, 0, currentX+widgetWidth, height-1)
		currentX += widgetWidth
	}

	// Render
	draw.Render(widgetList...)

	// Event loop
	for e := range draw.PollEvents() {
		if e.Type == draw.KeyboardEvent && e.ID == "<Escape>" {
			break
		}
	}
}
