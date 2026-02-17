package draw

import (
	"image"
	"sync"
)

// FrameBuffer stores the last rendered state of the terminal screen.
// This enables diff-based rendering by comparing current frame with previous frame.
type FrameBuffer struct {
	// CellMap stores the last rendered state of each cell
	CellMap map[image.Point]Cell
	// Bounds represents the terminal dimensions
	Bounds image.Rectangle
	mu     sync.RWMutex
}

// NewFrameBuffer creates a new frame buffer with the given terminal bounds
func NewFrameBuffer(bounds image.Rectangle) *FrameBuffer {
	return &FrameBuffer{
		CellMap: make(map[image.Point]Cell),
		Bounds:  bounds,
	}
}

// Resize updates the frame buffer bounds (called on terminal resize)
func (fb *FrameBuffer) Resize(bounds image.Rectangle) {
	fb.mu.Lock()
	defer fb.mu.Unlock()
	
	// Clear cells outside new bounds
	for p := range fb.CellMap {
		if !p.In(bounds) {
			delete(fb.CellMap, p)
		}
	}
	fb.Bounds = bounds
}

// GetCell returns the cell at the given point from the last frame
func (fb *FrameBuffer) GetCell(p image.Point) (Cell, bool) {
	fb.mu.RLock()
	defer fb.mu.RUnlock()
	cell, exists := fb.CellMap[p]
	return cell, exists
}

// SetCell stores a cell in the frame buffer
func (fb *FrameBuffer) SetCell(p image.Point, c Cell) {
	fb.mu.Lock()
	defer fb.mu.Unlock()
	fb.CellMap[p] = c
}

// Diff compares the current buffer with the frame buffer and returns:
// - changed: cells that have changed (need to be updated)
// - removed: cells that were in the previous frame but not in current (need to be cleared)
func (fb *FrameBuffer) Diff(current *Buffer) (changed map[image.Point]Cell, removed []image.Point) {
	fb.mu.RLock()
	defer fb.mu.RUnlock()
	
	changed = make(map[image.Point]Cell)
	removed = make([]image.Point, 0)
	
	// Track which cells exist in current frame
	currentCells := make(map[image.Point]bool)
	
	// Find changed and new cells
	for p, newCell := range current.CellMap {
		if !p.In(current.Rectangle) {
			continue
		}
		currentCells[p] = true
		
		oldCell, exists := fb.CellMap[p]
		if !exists || !cellsEqual(oldCell, newCell) {
			changed[p] = newCell
		}
	}
	
	// Find removed cells (were in previous frame but not in current)
	for p := range fb.CellMap {
		if !currentCells[p] && p.In(current.Rectangle) {
			removed = append(removed, p)
		}
	}
	
	return changed, removed
}

// Update updates the frame buffer with the current buffer state
func (fb *FrameBuffer) Update(buf *Buffer) {
	fb.mu.Lock()
	defer fb.mu.Unlock()
	
	// Update all cells in the buffer
	for p, cell := range buf.CellMap {
		if p.In(buf.Rectangle) {
			fb.CellMap[p] = cell
		}
	}
	
	// Remove cells that are outside the buffer's rectangle
	for p := range fb.CellMap {
		if !p.In(buf.Rectangle) && p.In(fb.Bounds) {
			// Only remove if it's outside current buffer but was in previous frame
			// This handles the case where a widget moved or was removed
			delete(fb.CellMap, p)
		}
	}
}

// Clear clears the entire frame buffer
func (fb *FrameBuffer) Clear() {
	fb.mu.Lock()
	defer fb.mu.Unlock()
	fb.CellMap = make(map[image.Point]Cell)
}

// cellsEqual compares two cells for equality
func cellsEqual(a, b Cell) bool {
	return a.Rune == b.Rune &&
		a.Style.Fg == b.Style.Fg &&
		a.Style.Bg == b.Style.Bg &&
		a.Style.Modifier == b.Style.Modifier
}
