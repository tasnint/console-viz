Sample Data Files for Console-Viz Visualizer

This folder contains sample data files that can be visualized using the data_visualizer example.

Supported file types:
- CSV files (.csv) - Will be displayed as tables and charts
- JSON files (.json) - Structured data visualization
- Text files (.txt) - Displayed as lists

To add your own data:
1. Place CSV, JSON, or TXT files in this folder
2. Run: go run examples/data_visualizer.go
3. The visualizer will automatically detect and display your files

CSV Format:
- First row should contain headers
- Subsequent rows contain data
- Numeric columns will be used for charts

JSON Format:
- Any valid JSON structure
- Arrays of numbers will be plotted as charts

Text Format:
- One item per line
- Displayed as a scrollable list
