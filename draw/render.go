package draw

import (
	"console-viz/styling"
	"image"
	"sync"

	tb "github.com/nsf/termbox-go"
)

// Drawable interface that all widgets implement
type Drawable interface {
	GetRect() image.Rectangle   // return widget's rectangle
	SetRect(int, int, int, int) // set widget's rectangle
	Draw(*Buffer)               // draw widget into buffer
	sync.Locker                 // embeds locking methods for the widget
}

// Renderer manages the rendering state and implements diff-based rendering
type Renderer struct {
	frameBuffer *FrameBuffer
	enabled     bool // whether diff-based rendering is enabled
}

// NewRenderer creates a new renderer with diff-based rendering enabled
func NewRenderer() *Renderer {
	tb.Sync() // sync termbox to get current terminal size
	w, h := tb.Size()
	return &Renderer{
		frameBuffer: NewFrameBuffer(image.Rect(0, 0, w, h)),
		enabled:     true,
	}
}

// Resize updates the renderer's frame buffer bounds (call this on terminal resize events)
func (r *Renderer) Resize(w, h int) {
	r.frameBuffer.Resize(image.Rect(0, 0, w, h))
}

// Render renders widgets using diff-based rendering (only updates changed cells)
func (r *Renderer) Render(items ...Drawable) {
	if !r.enabled {
		// Fallback to full redraw if diff-based rendering is disabled
		renderFull(items...)
		return
	}

	// Collect all buffers from widgets
	buffers := make([]*Buffer, 0, len(items))
	for _, item := range items {
		buf := NewBuffer(item.GetRect())
		item.Lock()
		item.Draw(buf)
		item.Unlock()
		buffers = append(buffers, buf)
	}

	// Compute diffs and render only changed cells
	changedSet := make(map[image.Point]Cell)
	removedSet := make(map[image.Point]bool)

	for _, buf := range buffers {
		changed, removed := r.frameBuffer.Diff(buf)
		for p, cell := range changed {
			changedSet[p] = cell
		}
		for _, p := range removed {
			removedSet[p] = true
		}
		// Update frame buffer with current state
		r.frameBuffer.Update(buf)
	}

	// Render changed cells
	for p, cell := range changedSet {
		if p.In(r.frameBuffer.Bounds) {
			r.renderCell(p, cell)
		}
	}

	// Clear removed cells
	for p := range removedSet {
		if p.In(r.frameBuffer.Bounds) {
			// Clear cell by rendering a blank cell
			r.renderCell(p, CellClear)
		}
	}

	tb.Flush()
}

// renderCell renders a single cell to the terminal
func (r *Renderer) renderCell(p image.Point, c Cell) {
	fg := tb.Attribute(c.Style.Fg + 1)
	if c.Style.Fg == styling.ColorClear {
		fg = tb.ColorDefault
	}
	bg := tb.Attribute(c.Style.Bg + 1)
	if c.Style.Bg == styling.ColorClear {
		bg = tb.ColorDefault
	}
	tb.SetCell(p.X, p.Y, c.Rune, fg|tb.Attribute(c.Style.Modifier), bg)
}

// renderFull is the full redraw implementation (fallback when diff is disabled)
func renderFull(items ...Drawable) {
	for _, item := range items {
		buf := NewBuffer(item.GetRect())
		item.Lock()
		item.Draw(buf)
		item.Unlock()
		for point, cell := range buf.CellMap {
			if point.In(buf.Rectangle) {
				fg := tb.Attribute(cell.Style.Fg + 1)
				if cell.Style.Fg == styling.ColorClear {
					fg = tb.ColorDefault
				}
				bg := tb.Attribute(cell.Style.Bg + 1)
				if cell.Style.Bg == styling.ColorClear {
					bg = tb.ColorDefault
				}
				tb.SetCell(
					point.X, point.Y,
					cell.Rune,
					fg|tb.Attribute(cell.Style.Modifier),
					bg,
				)
			}
		}
	}
	tb.Flush()
}

// Global renderer instance (for backward compatibility)
var globalRenderer *Renderer

// InitRenderer initializes the global renderer (should be called after termbox.Init())
func InitRenderer() {
	globalRenderer = NewRenderer()
}

// Render is the convenience function that uses the global renderer
// This maintains backward compatibility with the original API
func Render(items ...Drawable) {
	if globalRenderer == nil {
		// Auto-initialize if not already done
		InitRenderer()
	}
	globalRenderer.Render(items...)
}

// ResizeRenderer updates the renderer on terminal resize (call this from resize event handler)
func ResizeRenderer(w, h int) {
	if globalRenderer != nil {
		globalRenderer.Resize(w, h)
	}
}
