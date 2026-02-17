// +build ignore

package main

import (
	"console-viz/draw"
	"console-viz/styling"
	"console-viz/widgets"
	"log"
)

// stringValue is a helper type to convert strings to fmt.Stringer
type stringValue string

func (s stringValue) String() string {
	return string(s)
}

// Minimal test to verify the project is runnable
func main() {
	// Initialize theme
	cliTheme := styling.ParseThemeFlag()
	if err := styling.InitThemeFromCLI(cliTheme); err != nil {
		log.Printf("Warning: Theme initialization failed: %v", err)
	}

	// Initialize termbox
	if err := draw.Init(); err != nil {
		log.Fatalf("Failed to initialize terminal: %v", err)
	}
	defer draw.Close()

	// Initialize renderer
	draw.InitRenderer()

	// Test creating and rendering each widget type
	log.Println("Testing widget creation...")

	// Test Paragraph
	p := widgets.NewParagraph()
	p.Text = "Hello World!"
	p.Title = "Paragraph Test"
	p.SetRect(0, 0, 30, 5)
	log.Println("✓ Paragraph created")

	// Test List
	l := widgets.NewList()
	l.Title = "List Test"
	l.Rows = []string{"Item 1", "Item 2", "Item 3"}
	l.SetRect(0, 5, 25, 12)
	log.Println("✓ List created")

	// Test Gauge
	g := widgets.NewGauge()
	g.Title = "Gauge Test"
	g.Percent = 50
	g.SetRect(0, 12, 30, 15)
	log.Println("✓ Gauge created")

	// Test BarChart
	bc := widgets.NewBarChart()
	bc.Title = "BarChart Test"
	bc.Data = []float64{3, 2, 5, 3, 9}
	bc.Labels = []string{"A", "B", "C", "D", "E"}
	bc.SetRect(30, 0, 60, 10)
	log.Println("✓ BarChart created")

	// Test Table
	t := widgets.NewTable()
	t.Title = "Table Test"
	t.Rows = [][]string{
		{"Name", "Age"},
		{"Alice", "25"},
		{"Bob", "30"},
	}
	t.SetRect(30, 10, 60, 17)
	log.Println("✓ Table created")

	// Test Sparkline
	sl := widgets.NewSparkline()
	sl.Title = "Sparkline Test"
	sl.Data = []float64{4, 2, 1, 6, 3, 9, 1, 4, 2, 15}
	slg := widgets.NewSparklineGroup(sl)
	slg.Title = "Sparkline Group"
	slg.SetRect(0, 15, 30, 22)
	log.Println("✓ Sparkline created")

	// Test Plot
	plot := widgets.NewPlot()
	plot.Title = "Plot Test"
	plot.Data = [][]float64{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}
	plot.SetRect(30, 17, 60, 25)
	log.Println("✓ Plot created")

	// Test PieChart
	pc := widgets.NewPieChart()
	pc.Title = "PieChart Test"
	pc.Data = []float64{25, 30, 20, 15, 10}
	pc.SetRect(0, 22, 30, 30)
	log.Println("✓ PieChart created")

	// Test Tree
	tree := widgets.NewTree()
	tree.Title = "Tree Test"
	// Create tree nodes with Stringer values
	rootNode := &widgets.TreeNode{
		Value:    stringValue("Root"),
		Expanded: true,
		Nodes: []*widgets.TreeNode{
			{Value: stringValue("Child 1")},
			{Value: stringValue("Child 2")},
		},
	}
	tree.SetNodes([]*widgets.TreeNode{rootNode})
	tree.SetRect(30, 25, 60, 35)
	log.Println("✓ Tree created")

	// Test Tabs
	tabs := widgets.NewTabPane("Tab 1", "Tab 2", "Tab 3")
	tabs.Title = "Tabs Test"
	tabs.SetRect(0, 30, 60, 35)
	log.Println("✓ Tabs created")

	// Test StackedBarChart
	sbc := widgets.NewStackedBarChart()
	sbc.Title = "StackedBarChart Test"
	sbc.Data = [][]float64{{3, 2, 5}, {1, 2, 3}}
	sbc.Labels = []string{"A", "B", "C"}
	sbc.SetRect(0, 35, 60, 42)
	log.Println("✓ StackedBarChart created")

	log.Println("\nAll widgets created successfully!")
	log.Println("Attempting to render...")

	// Try to render all widgets
	draw.Render(p, l, g, bc, t, slg, plot, pc, tree, tabs, sbc)

	log.Println("✓ Rendering successful!")
	log.Println("\nProject is runnable! Press ESC to exit.")

	// Wait for ESC
	for e := range draw.PollEvents() {
		if e.Type == draw.KeyboardEvent && e.ID == "<Escape>" {
			break
		}
	}
}
