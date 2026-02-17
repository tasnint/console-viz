package widgets

import (
	"console-viz/draw"
	"console-viz/styling"
	"console-viz/utils"
	"fmt"
	"image"
	"strings"

	rw "github.com/mattn/go-runewidth"
)

const treeIndent = "  "

// TreeNode represents a node in a tree structure
type TreeNode struct {
	Value    fmt.Stringer // Value to display
	Expanded bool         // Whether node is expanded
	Nodes    []*TreeNode  // Child nodes
	level    int          // Node level in tree (internal)
}

// TreeWalkFn is a function used for walking a Tree
// Return false to interrupt the walking process
type TreeWalkFn func(*TreeNode) bool

// parseStyles converts node value to styled cells
func (tn *TreeNode) parseStyles(style styling.Style) []draw.Cell {
	var sb strings.Builder
	if len(tn.Nodes) == 0 {
		sb.WriteString(strings.Repeat(treeIndent, tn.level+1))
	} else {
		sb.WriteString(strings.Repeat(treeIndent, tn.level))
		theme := styling.GetTheme()
		if tn.Expanded {
			sb.WriteRune(theme.Tree.Expanded)
		} else {
			sb.WriteRune(theme.Tree.Collapsed)
		}
		sb.WriteByte(' ')
	}
	sb.WriteString(tn.Value.String())
	return draw.ParseStyles(sb.String(), style)
}

// Tree displays a hierarchical tree structure
type Tree struct {
	draw.Base
	TextStyle        styling.Style   // Default text style
	SelectedRowStyle styling.Style   // Style for selected row
	WrapText         bool            // Whether to wrap text
	SelectedRow      int             // Currently selected row index
	nodes            []*TreeNode     // Root nodes
	rows             []*TreeNode     // Flattened nodes for rendering
	topRow           int             // Top visible row (for scrolling)
}

// NewTree creates a new Tree widget with default settings
func NewTree() *Tree {
	theme := styling.GetTheme()
	return &Tree{
		Base:            *draw.NewBase(),
		TextStyle:       theme.Tree.Text,
		SelectedRowStyle: theme.Tree.Text,
		WrapText:        true,
	}
}

// SetNodes sets the root nodes of the tree
func (t *Tree) SetNodes(nodes []*TreeNode) {
	t.nodes = nodes
	t.prepareNodes()
}

// prepareNodes flattens the tree structure for rendering
func (t *Tree) prepareNodes() {
	t.rows = make([]*TreeNode, 0)
	for _, node := range t.nodes {
		t.prepareNode(node, 0)
	}
}

// prepareNode recursively flattens nodes
func (t *Tree) prepareNode(node *TreeNode, level int) {
	t.rows = append(t.rows, node)
	node.level = level

	if node.Expanded {
		for _, n := range node.Nodes {
			t.prepareNode(n, level+1)
		}
	}
}

// Walk walks through all nodes calling fn for each
func (t *Tree) Walk(fn TreeWalkFn) {
	for _, n := range t.nodes {
		if !t.walk(n, fn) {
			break
		}
	}
}

// walk recursively walks through nodes
func (t *Tree) walk(n *TreeNode, fn TreeWalkFn) bool {
	if !fn(n) {
		return false
	}

	for _, node := range n.Nodes {
		if !t.walk(node, fn) {
			return false
		}
	}

	return true
}

// Draw renders the tree widget
func (t *Tree) Draw(buf *draw.Buffer) {
	t.Base.Draw(buf)
	point := t.Inner.Min

	// Adjust view to show selected row
	if t.SelectedRow >= t.Inner.Dy()+t.topRow {
		t.topRow = t.SelectedRow - t.Inner.Dy() + 1
	} else if t.SelectedRow < t.topRow {
		t.topRow = t.SelectedRow
	}

	// Draw rows
	for row := t.topRow; row < len(t.rows) && point.Y < t.Inner.Max.Y; row++ {
		cells := t.rows[row].parseStyles(t.TextStyle)
		if t.WrapText {
			cells = utils.WrapCells(cells, uint(t.Inner.Dx()))
		}

		for j := 0; j < len(cells) && point.Y < t.Inner.Max.Y; j++ {
			style := cells[j].Style
			if row == t.SelectedRow {
				style = t.SelectedRowStyle
			}

			if point.X+1 == t.Inner.Max.X+1 && len(cells) > t.Inner.Dx() {
				buf.SetCell(draw.NewCell(styling.ELLIPSES, style), point.Add(image.Pt(-1, 0)))
			} else {
				buf.SetCell(draw.NewCell(cells[j].Rune, style), point)
				point = point.Add(image.Pt(rw.RuneWidth(cells[j].Rune), 0))
			}
		}
		point = image.Pt(t.Inner.Min.X, point.Y+1)
	}

	// Draw scroll indicators
	if t.topRow > 0 {
		buf.SetCell(
			draw.NewCell(styling.UP_ARROW, styling.NewStyle(styling.ColorWhite)),
			image.Pt(t.Inner.Max.X-1, t.Inner.Min.Y),
		)
	}

	if len(t.rows) > int(t.topRow)+t.Inner.Dy() {
		buf.SetCell(
			draw.NewCell(styling.DOWN_ARROW, styling.NewStyle(styling.ColorWhite)),
			image.Pt(t.Inner.Max.X-1, t.Inner.Max.Y-1),
		)
	}
}

// ScrollAmount scrolls by the given amount
func (t *Tree) ScrollAmount(amount int) {
	if len(t.rows)-int(t.SelectedRow) <= amount {
		t.SelectedRow = len(t.rows) - 1
	} else if int(t.SelectedRow)+amount < 0 {
		t.SelectedRow = 0
	} else {
		t.SelectedRow += amount
	}
}

// SelectedNode returns the currently selected node
func (t *Tree) SelectedNode() *TreeNode {
	if len(t.rows) == 0 {
		return nil
	}
	return t.rows[t.SelectedRow]
}

// ScrollUp scrolls up one row
func (t *Tree) ScrollUp() {
	t.ScrollAmount(-1)
}

// ScrollDown scrolls down one row
func (t *Tree) ScrollDown() {
	t.ScrollAmount(1)
}

// ScrollPageUp scrolls up one page
func (t *Tree) ScrollPageUp() {
	if t.SelectedRow > t.topRow {
		t.SelectedRow = t.topRow
	} else {
		t.ScrollAmount(-t.Inner.Dy())
	}
}

// ScrollPageDown scrolls down one page
func (t *Tree) ScrollPageDown() {
	t.ScrollAmount(t.Inner.Dy())
}

// ScrollTop scrolls to the top
func (t *Tree) ScrollTop() {
	t.SelectedRow = 0
}

// ScrollBottom scrolls to the bottom
func (t *Tree) ScrollBottom() {
	t.SelectedRow = len(t.rows) - 1
}

// Collapse collapses the selected node
func (t *Tree) Collapse() {
	if t.SelectedRow < len(t.rows) {
		t.rows[t.SelectedRow].Expanded = false
		t.prepareNodes()
	}
}

// Expand expands the selected node
func (t *Tree) Expand() {
	if t.SelectedRow < len(t.rows) {
		node := t.rows[t.SelectedRow]
		if len(node.Nodes) > 0 {
			node.Expanded = true
		}
		t.prepareNodes()
	}
}

// ToggleExpand toggles expansion of the selected node
func (t *Tree) ToggleExpand() {
	if t.SelectedRow < len(t.rows) {
		node := t.rows[t.SelectedRow]
		if len(node.Nodes) > 0 {
			node.Expanded = !node.Expanded
		}
		t.prepareNodes()
	}
}

// ExpandAll expands all nodes
func (t *Tree) ExpandAll() {
	t.Walk(func(n *TreeNode) bool {
		if len(n.Nodes) > 0 {
			n.Expanded = true
		}
		return true
	})
	t.prepareNodes()
}

// CollapseAll collapses all nodes
func (t *Tree) CollapseAll() {
	t.Walk(func(n *TreeNode) bool {
		n.Expanded = false
		return true
	})
	t.prepareNodes()
}
