// +build ignore

// mock_json.go loads MOCK_DATA.json: table, vertical barchart (ID), and gender distribution (vertical).
// Run from repo root in a real terminal:
//
//	go run ./examples/mock_json.go
//
// Press Escape to exit.
package main

import (
	"console-viz/draw"
	"console-viz/widgets"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	// Find MOCK_DATA.json (repo root or next to binary)
	dataPath := "MOCK_DATA.json"
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		// Try from examples dir
		dataPath = filepath.Join("examples", "..", "MOCK_DATA.json")
		dataPath = filepath.Clean(dataPath)
	}
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		log.Fatalf("MOCK_DATA.json not found. Run from repo root: go run ./examples/mock_json.go")
	}

	raw, err := os.ReadFile(dataPath)
	if err != nil {
		log.Fatalf("Failed to read MOCK_DATA.json: %v", err)
	}

	var arr []map[string]interface{}
	if err := json.Unmarshal(raw, &arr); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Build table rows: header + data (id, first_name, last_name, gender)
	rows := [][]string{
		{"ID", "First Name", "Last Name", "Gender"},
	}
	for _, obj := range arr {
		rows = append(rows, []string{
			fmt.Sprintf("%v", obj["id"]),
			fmt.Sprintf("%v", obj["first_name"]),
			fmt.Sprintf("%v", obj["last_name"]),
			fmt.Sprintf("%v", obj["gender"]),
		})
	}

	// Numeric series from id for charts
	ids := make([]float64, 0, len(arr))
	for _, obj := range arr {
		if v, ok := obj["id"]; ok {
			switch n := v.(type) {
			case float64:
				ids = append(ids, n)
			case int:
				ids = append(ids, float64(n))
			}
		}
	}

	// Gender distribution: count per gender
	genderCounts := make(map[string]int)
	for _, obj := range arr {
		if g, ok := obj["gender"]; ok {
			genderCounts[fmt.Sprintf("%v", g)]++
		}
	}
	// Sort genders by name for stable order; build labels and values
	type kv struct{ label string; count int }
	var gPairs []kv
	for label, count := range genderCounts {
		gPairs = append(gPairs, kv{label, count})
	}
	sort.Slice(gPairs, func(i, j int) bool { return gPairs[i].label < gPairs[j].label })
	genderLabels := make([]string, len(gPairs))
	genderValues := make([]float64, len(gPairs))
	for i, p := range gPairs {
		genderLabels[i] = p.label
		genderValues[i] = float64(p.count)
	}

	// Regular (vertical) barchart: first N IDs
	const barCount = 18
	barValues := ids
	if len(barValues) > barCount {
		barValues = barValues[:barCount]
	}
	barLabels := make([]string, len(barValues))
	for i := range barValues {
		barLabels[i] = fmt.Sprintf("%d", i+1)
	}

	if err := draw.Init(); err != nil {
		log.Fatalf("Failed to initialize terminal: %v", err)
	}
	defer draw.Close()

	draw.InitRenderer()
	width, height := draw.TerminalDimensions()

	// Layout: table ~50% | vertical barchart ~25% | gender distribution ~25%
	tableW := int(float64(width) * 0.5)
	barW := int(float64(width) * 0.25)

	table := widgets.NewTable()
	table.Title = fmt.Sprintf("MOCK_DATA.json (%d rows)", len(rows)-1)
	table.Rows = rows
	table.SetRect(0, 0, tableW, height)

	bc := widgets.NewBarChart()
	bc.Title = "ID (first 18)"
	bc.Data = barValues
	bc.Labels = barLabels
	bc.SetRect(tableW, 0, tableW+barW, height)

	genderChart := widgets.NewBarChart()
	genderChart.Title = "Gender Distribution"
	genderChart.Data = genderValues
	genderChart.Labels = genderLabels
	genderChart.SetRect(tableW+barW, 0, width, height)

	draw.Render(table, bc, genderChart)

	for e := range draw.PollEvents() {
		if e.Type == draw.KeyboardEvent && e.ID == "<Escape>" {
			break
		}
	}
}
