// +build ignore

package main

import (
	"console-viz/draw"
	"console-viz/styling"
	"console-viz/widgets"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// DataFile represents a data file in the data folder
type DataFile struct {
	Name     string
	Path     string
	Type     string // "csv", "json", "txt"
	Contents interface{}
}

// CSVData represents CSV file data
type CSVData struct {
	Headers []string
	Rows    [][]string
}

// JSONData represents JSON file data
type JSONData struct {
	Data interface{}
}

// LoadDataFolder loads all data files from a folder
func LoadDataFolder(folderPath string) ([]DataFile, error) {
	files := []DataFile{}
	
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			return nil
		}
		
		ext := strings.ToLower(filepath.Ext(path))
		var fileType string
		
		switch ext {
		case ".csv":
			fileType = "csv"
		case ".json":
			fileType = "json"
		case ".txt":
			fileType = "txt"
		default:
			return nil // Skip unknown file types
		}
		
		files = append(files, DataFile{
			Name: info.Name(),
			Path: path,
			Type: fileType,
		})
		
		return nil
	})
	
	return files, err
}

// LoadCSV loads a CSV file
func LoadCSV(filePath string) (*CSVData, error) {
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
	
	if len(records) == 0 {
		return &CSVData{Headers: []string{}, Rows: [][]string{}}, nil
	}
	
	return &CSVData{
		Headers: records[0],
		Rows:    records[1:],
	}, nil
}

// LoadJSON loads a JSON file
func LoadJSON(filePath string) (*JSONData, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, err
	}
	
	return &JSONData{Data: jsonData}, nil
}

// LoadText loads a text file
func LoadText(filePath string) ([]string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(string(data), "\n")
	// Remove empty lines at the end
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	
	return lines, nil
}

// ParseNumericColumn extracts numeric values from a CSV column
func ParseNumericColumn(csvData *CSVData, columnIndex int) []float64 {
	values := []float64{}
	for _, row := range csvData.Rows {
		if columnIndex < len(row) {
			if val, err := strconv.ParseFloat(strings.TrimSpace(row[columnIndex]), 64); err == nil {
				values = append(values, val)
			}
		}
	}
	return values
}

// stringValue helper for Tree widget
type stringValue string

func (s stringValue) String() string {
	return string(s)
}

func main() {
	// Parse command line arguments
	dataFolder := "data"
	if len(os.Args) > 1 {
		dataFolder = os.Args[1]
	}
	
	// Initialize theme
	cliTheme := styling.ParseThemeFlag()
	if err := styling.InitThemeFromCLI(cliTheme); err != nil {
		log.Printf("Warning: Theme initialization failed: %v", err)
	}
	
	// Initialize terminal
	if err := draw.Init(); err != nil {
		log.Fatalf("Failed to initialize terminal: %v", err)
	}
	defer draw.Close()
	
	draw.InitRenderer()
	
	// Load data files
	log.Printf("Loading data from folder: %s", dataFolder)
	files, err := LoadDataFolder(dataFolder)
	if err != nil {
		log.Fatalf("Failed to load data folder: %v", err)
	}
	
	if len(files) == 0 {
		log.Fatalf("No data files found in %s. Please add CSV, JSON, or TXT files.", dataFolder)
	}
	
	log.Printf("Found %d data files", len(files))
	
	// Get terminal dimensions
	width, height := draw.TerminalDimensions()
	
	// Create widgets
	widgetsList := []draw.Drawable{}
	
	// Title widget
	title := widgets.NewParagraph()
	title.Text = fmt.Sprintf("Data Visualizer - %d files loaded", len(files))
	title.Title = "Console-Viz Data Visualizer"
	title.SetRect(0, 0, width, 3)
	widgetsList = append(widgetsList, title)
	
	// File list widget
	fileList := widgets.NewList()
	fileList.Title = "Data Files"
	fileRows := []string{}
	for i, f := range files {
		fileRows = append(fileRows, fmt.Sprintf("[%d] %s (%s)", i+1, f.Name, f.Type))
	}
	fileList.Rows = fileRows
	fileList.SetRect(0, 3, width/3, height-5)
	widgetsList = append(widgetsList, fileList)
	
	// Try to load and visualize the first CSV file
	var tableWidget *widgets.Table
	var chartWidget *widgets.BarChart
	var plotWidget *widgets.Plot
	
	for _, file := range files {
		if file.Type == "csv" {
			csvData, err := LoadCSV(file.Path)
			if err != nil {
				log.Printf("Error loading CSV %s: %v", file.Name, err)
				continue
			}
			
			// Create table widget
			tableWidget = widgets.NewTable()
			tableWidget.Title = fmt.Sprintf("Table: %s", file.Name)
			tableWidget.Rows = append([][]string{csvData.Headers}, csvData.Rows...)
			tableWidget.SetRect(width/3, 3, width*2/3, height/2)
			
			// Create bar chart if we have numeric data
			if len(csvData.Headers) > 1 && len(csvData.Rows) > 0 {
				// Try to parse first numeric column
				numericData := ParseNumericColumn(csvData, 1)
				if len(numericData) > 0 {
					chartWidget = widgets.NewBarChart()
					chartWidget.Title = fmt.Sprintf("Chart: %s", file.Name)
					chartWidget.Data = numericData
					labels := []string{}
					for i := range numericData {
						if i < len(csvData.Rows) && len(csvData.Rows[i]) > 0 {
							labels = append(labels, csvData.Rows[i][0])
						} else {
							labels = append(labels, fmt.Sprintf("Item %d", i+1))
						}
					}
					chartWidget.Labels = labels
					chartWidget.SetRect(width*2/3, 3, width, height/2)
				}
				
				// Create plot widget
				if len(numericData) > 0 {
					plotWidget = widgets.NewPlot()
					plotWidget.Title = fmt.Sprintf("Plot: %s", file.Name)
					plotWidget.Data = [][]float64{numericData}
					plotWidget.SetRect(width/3, height/2, width, height-5)
				}
			}
			break
		} else if file.Type == "txt" {
			// Load text file for list display
			textLines, err := LoadText(file.Path)
			if err != nil {
				log.Printf("Error loading text file %s: %v", file.Name, err)
				continue
			}
			
			textList := widgets.NewList()
			textList.Title = fmt.Sprintf("Text: %s", file.Name)
			textList.Rows = textLines
			textList.SetRect(width/3, 3, width, height-5)
			widgetsList = append(widgetsList, textList)
			break
		}
	}
	
	// Add widgets to list
	if tableWidget != nil {
		widgetsList = append(widgetsList, tableWidget)
	}
	if chartWidget != nil {
		widgetsList = append(widgetsList, chartWidget)
	}
	if plotWidget != nil {
		widgetsList = append(widgetsList, plotWidget)
	}
	
	// Instructions widget
	instructions := widgets.NewParagraph()
	instructions.Text = "Press ESC to exit | Use arrow keys to navigate file list"
	instructions.Border = false
	instructions.SetRect(0, height-2, width, height)
	widgetsList = append(widgetsList, instructions)
	
	// Render all widgets
	draw.Render(widgetsList...)
	
	// Event loop
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	
	selectedFile := 0
	
	for {
		select {
		case e := <-draw.PollEvents():
			if e.Type == draw.KeyboardEvent {
				if e.ID == "<Escape>" {
					return
				} else if e.ID == "<Up>" && selectedFile > 0 {
					selectedFile--
					fileList.SelectedRow = selectedFile
					draw.Render(widgetsList...)
				} else if e.ID == "<Down>" && selectedFile < len(files)-1 {
					selectedFile++
					fileList.SelectedRow = selectedFile
					draw.Render(widgetsList...)
				} else if e.ID == "<Enter>" {
					// Load selected file
					if selectedFile < len(files) {
						file := files[selectedFile]
						log.Printf("Loading file: %s", file.Name)
						// You can add logic here to reload and re-render the selected file
					}
				}
			} else if e.Type == draw.ResizeEvent {
				// Handle resize
				if resize, ok := e.Payload.(draw.Resize); ok {
					draw.ResizeRenderer(resize.Width, resize.Height)
					// Recalculate widget positions
					width, height = resize.Width, resize.Height
					// Update widget positions (simplified - you'd want to recalculate all)
					draw.Render(widgetsList...)
				}
			}
		case <-ticker.C:
			// Auto-refresh (optional)
			draw.Render(widgetsList...)
		}
	}
}
