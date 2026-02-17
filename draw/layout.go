package draw

import (
	"reflect"
)

// Alignment defines text alignment within a layout item
type Alignment uint

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
)

// LayoutItemType represents whether a layout item is a row or column
type LayoutItemType uint

const (
	LayoutItemColumn LayoutItemType = iota // Column - divides space horizontally
	LayoutItemRow                          // Row - divides space vertically
)

// LayoutItem represents either a Row or Column in a layout, or a widget
// Holds sizing information and either nested LayoutItems or a widget
type LayoutItem struct {
	Type        LayoutItemType // Column or Row
	XRatio      float64        // X position ratio (0.0 to 1.0)
	YRatio      float64        // Y position ratio (0.0 to 1.0)
	WidthRatio  float64        // Width ratio (0.0 to 1.0)
	HeightRatio float64        // Height ratio (0.0 to 1.0)
	Entry       interface{}    // Either a Drawable (widget) or []LayoutItem (nested)
	IsLeaf      bool           // True if Entry is a widget, false if nested items
	Ratio       float64        // Size ratio for this item (used during construction)
	Align       Alignment      // Alignment for the item (only used for leaf items)
}

// Layout is an advanced layout system that supports nested rows and columns
// Similar to CSS Grid or Flexbox, allowing complex widget arrangements
type Layout struct {
	Base
	Items []*LayoutItem // Flattened list of all leaf items (widgets) with calculated positions
}

// NewLayout creates a new Layout container
func NewLayout() *Layout {
	l := &Layout{
		Base:  *NewBase(),
		Items: make([]*LayoutItem, 0),
	}
	l.Border = false
	return l
}

// NewLayoutColumn creates a column layout item that divides space horizontally
// Takes a width ratio (0.0 to 1.0) and either a widget or nested layout items
// Example: NewLayoutColumn(0.5, widget) creates a column taking 50% width
func NewLayoutColumn(ratio float64, items ...interface{}) LayoutItem {
	if len(items) == 0 {
		return LayoutItem{
			Type:   LayoutItemColumn,
			Ratio:  ratio,
			IsLeaf: false,
			Entry:  []LayoutItem{},
		}
	}

	// Check if first item is a Drawable (widget)
	_, isDrawable := items[0].(Drawable)
	entry := items[0]

	// If not a Drawable, treat all items as nested LayoutItems
	if !isDrawable {
		entry = items
	}

	return LayoutItem{
		Type:   LayoutItemColumn,
		Entry:  entry,
		IsLeaf: isDrawable,
		Ratio:  ratio,
	}
}

// NewLayoutRow creates a row layout item that divides space vertically
// Takes a height ratio (0.0 to 1.0) and either a widget or nested layout items
// Example: NewLayoutRow(0.33, widget) creates a row taking 33% height
func NewLayoutRow(ratio float64, items ...interface{}) LayoutItem {
	if len(items) == 0 {
		return LayoutItem{
			Type:   LayoutItemRow,
			Ratio:  ratio,
			IsLeaf: false,
			Entry:  []LayoutItem{},
		}
	}

	// Check if first item is a Drawable (widget)
	_, isDrawable := items[0].(Drawable)
	entry := items[0]

	// If not a Drawable, treat all items as nested LayoutItems
	if !isDrawable {
		entry = items
	}

	return LayoutItem{
		Type:   LayoutItemRow,
		Entry:  entry,
		IsLeaf: isDrawable,
		Ratio:  ratio,
	}
}

// Set configures the layout with nested rows and columns
// Recursively processes LayoutItems, calculating positions for all widgets
// Example:
//   layout.Set(
//     NewLayoutRow(0.5,
//       NewLayoutColumn(0.5, widget1),
//       NewLayoutColumn(0.5, widget2),
//     ),
//     NewLayoutRow(0.5, widget3),
//   )
func (l *Layout) Set(items ...interface{}) {
	// Clear existing items
	l.Items = make([]*LayoutItem, 0)

	// Create root row containing all items
	rootItem := LayoutItem{
		Type:   LayoutItemRow,
		Entry:  items,
		IsLeaf: false,
		Ratio:  1.0,
	}

	// Recursively process the layout tree
	l.layoutHelper(rootItem, 1.0, 1.0, 0.0, 0.0)
}

// layoutHelper recursively processes layout items and calculates positions
// This is the core algorithm that converts nested layout structure into flat item list
func (l *Layout) layoutHelper(item LayoutItem, parentWidthRatio, parentHeightRatio, parentXRatio, parentYRatio float64) {
	var heightRatio float64
	var widthRatio float64

	// Calculate width and height ratios based on item type
	switch item.Type {
	case LayoutItemColumn:
		// Column divides space horizontally
		heightRatio = 1.0
		widthRatio = item.Ratio
	case LayoutItemRow:
		// Row divides space vertically
		heightRatio = item.Ratio
		widthRatio = 1.0
	}

	// Calculate absolute ratios
	item.WidthRatio = parentWidthRatio * widthRatio
	item.HeightRatio = parentHeightRatio * heightRatio
	item.XRatio = parentXRatio
	item.YRatio = parentYRatio

	// If this is a leaf (widget), add it to the items list
	if item.IsLeaf {
		l.Items = append(l.Items, &item)
		return
	}

	// Process nested items
	xRatio := 0.0
	yRatio := 0.0
	hasColumns := false
	hasRows := false

	// Convert entry to slice of interfaces
	children := interfaceSlice(item.Entry)

	for i := 0; i < len(children); i++ {
		if children[i] == nil {
			continue
		}

		// Try to convert to LayoutItem
		child, ok := children[i].(LayoutItem)
		if !ok {
			// If not a LayoutItem, try to treat as Drawable
			if drawable, ok := children[i].(Drawable); ok {
				// Create a leaf item for the widget
				leafItem := LayoutItem{
					Type:        item.Type, // Inherit type from parent
					Entry:       drawable,
					IsLeaf:      true,
					Ratio:       1.0 / float64(len(children)), // Distribute evenly
					WidthRatio:  item.WidthRatio / float64(len(children)),
					HeightRatio: item.HeightRatio / float64(len(children)),
					XRatio:      item.XRatio + xRatio*item.WidthRatio,
					YRatio:      item.YRatio + yRatio*item.HeightRatio,
				}
				l.Items = append(l.Items, &leafItem)
				if item.Type == LayoutItemColumn {
					xRatio += leafItem.Ratio
				} else {
					yRatio += leafItem.Ratio
				}
			}
			continue
		}

		// Calculate child position
		child.XRatio = item.XRatio + (item.WidthRatio * xRatio)
		child.YRatio = item.YRatio + (item.HeightRatio * yRatio)

		// Track what types of children we have
		switch child.Type {
		case LayoutItemColumn:
			hasColumns = true
			xRatio += child.Ratio
			// If we have both rows and columns, adjust ratios
			if hasRows {
				item.HeightRatio /= 2.0
			}
		case LayoutItemRow:
			hasRows = true
			yRatio += child.Ratio
			// If we have both rows and columns, adjust ratios
			if hasColumns {
				item.WidthRatio /= 2.0
			}
		}

		// Recursively process child
		l.layoutHelper(child, item.WidthRatio, item.HeightRatio, child.XRatio, child.YRatio)
	}
}

// AddItem adds a single widget to the layout with manual positioning
// This is the simple API for basic layouts (backward compatible)
func (l *Layout) AddItem(
	drawable Drawable,
	xRatio float64,
	yRatio float64,
	widthRatio float64,
	heightRatio float64,
	align Alignment,
) {
	item := &LayoutItem{
		XRatio:      xRatio,
		YRatio:      yRatio,
		WidthRatio:  widthRatio,
		HeightRatio: heightRatio,
		Align:       align,
		IsLeaf:      true,
		Entry:       drawable,
	}

	l.Items = append(l.Items, item)
}

// Draw renders the layout and all its child widgets
func (l *Layout) Draw(buf *Buffer) {
	l.Base.Draw(buf) // Draw base (border, background, etc.)

	width := float64(l.Inner.Dx())
	height := float64(l.Inner.Dy())

	for _, item := range l.Items {
		// Get the drawable widget
		drawable, ok := item.Entry.(Drawable)
		if !ok {
			continue
		}

		// Calculate absolute coordinates
		x := int(width*item.XRatio) + l.Inner.Min.X
		y := int(height*item.YRatio) + l.Inner.Min.Y
		w := int(width * item.WidthRatio)
		h := int(height * item.HeightRatio)

		// Ensure we don't exceed bounds
		if x+w > l.Inner.Max.X {
			w = l.Inner.Max.X - x
		}
		if y+h > l.Inner.Max.Y {
			h = l.Inner.Max.Y - y
		}
		if w <= 0 || h <= 0 {
			continue
		}

		// Apply alignment (only for leaf items)
		if item.IsLeaf {
			switch item.Align {
			case AlignCenter:
				x += w / 4
			case AlignRight:
				x += w / 2
			}
		}

		// Set widget rectangle
		drawable.SetRect(x, y, x+w, y+h)

		// Draw widget
		drawable.Lock()
		drawable.Draw(buf)
		drawable.Unlock()
	}
}

// Clear removes all items from the layout
func (l *Layout) Clear() {
	l.Items = make([]*LayoutItem, 0)
}

// GetItemCount returns the number of items in the layout
func (l *Layout) GetItemCount() int {
	return len(l.Items)
}

// interfaceSlice converts an interface{} containing a slice to []interface{}
// Helper function to avoid import cycle with utils package
func interfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("interfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}
	return ret
}
