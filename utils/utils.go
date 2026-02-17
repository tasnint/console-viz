package utils

import (
	"console-viz/draw"
	"console-viz/styling"
	"fmt"
	"math"
	"reflect"

	rw "github.com/mattn/go-runewidth"
	wordwrap "github.com/mitchellh/go-wordwrap"
)

// ============================================================================
// Type Conversion Utilities
// ============================================================================

// InterfaceSlice converts an interface{} containing a slice to []interface{}
// Useful for working with variadic arguments and generic slice operations
// Example: InterfaceSlice([]int{1,2,3}) returns []interface{}{1,2,3}
func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}
	return ret
}

// ============================================================================
// String Utilities
// ============================================================================

// TrimString trims a string to a maximum width and adds ellipsis (…) if truncated
// Handles wide characters correctly using runewidth
// Returns empty string if width <= 0
func TrimString(s string, width int) string {
	if width <= 0 {
		return ""
	}
	ellipsis := "…" // Use ellipsis character directly
	if rw.StringWidth(s) > width {
		return rw.Truncate(s, width, ellipsis)
	}
	return s
}

// ============================================================================
// Color and Style Selection
// ============================================================================

// SelectColor cycles through a color slice, returning the color at the given index
// Uses modulo to wrap around if index exceeds slice length
// Useful for assigning colors to multiple data series
func SelectColor(colors []styling.Color, index int) styling.Color {
	if len(colors) == 0 {
		return styling.ColorWhite // fallback
	}
	return colors[index%len(colors)]
}

// SelectStyle cycles through a style slice, returning the style at the given index
// Uses modulo to wrap around if index exceeds slice length
// Useful for alternating styles in lists or tables
func SelectStyle(styles []styling.Style, index int) styling.Style {
	if len(styles) == 0 {
		return styling.StyleClear // fallback
	}
	return styles[index%len(styles)]
}

// ============================================================================
// Math Utilities
// ============================================================================

// SumIntSlice calculates the sum of all integers in a slice
func SumIntSlice(slice []int) int {
	sum := 0
	for _, val := range slice {
		sum += val
	}
	return sum
}

// SumFloat64Slice calculates the sum of all float64 values in a slice
func SumFloat64Slice(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum
}

// GetMaxIntFromSlice finds the maximum integer value in a slice
// Returns an error if the slice is empty
func GetMaxIntFromSlice(slice []int) (int, error) {
	if len(slice) == 0 {
		return 0, fmt.Errorf("cannot get max value from empty slice")
	}
	max := slice[0]
	for _, val := range slice[1:] {
		if val > max {
			max = val
		}
	}
	return max, nil
}

// GetMaxFloat64FromSlice finds the maximum float64 value in a slice
// Returns an error if the slice is empty
func GetMaxFloat64FromSlice(slice []float64) (float64, error) {
	if len(slice) == 0 {
		return 0, fmt.Errorf("cannot get max value from empty slice")
	}
	max := slice[0]
	for _, val := range slice[1:] {
		if val > max {
			max = val
		}
	}
	return max, nil
}

// GetMaxFloat64From2dSlice finds the maximum float64 value across all slices in a 2D slice
// Useful for finding the maximum value across multiple data series
// Returns an error if the slice is empty
func GetMaxFloat64From2dSlice(slices [][]float64) (float64, error) {
	if len(slices) == 0 {
		return 0, fmt.Errorf("cannot get max value from empty slice")
	}
	max := slices[0][0]
	for _, slice := range slices {
		for _, val := range slice {
			if val > max {
				max = val
			}
		}
	}
	return max, nil
}

// GetMinFloat64FromSlice finds the minimum float64 value in a slice
// Returns an error if the slice is empty
func GetMinFloat64FromSlice(slice []float64) (float64, error) {
	if len(slice) == 0 {
		return 0, fmt.Errorf("cannot get min value from empty slice")
	}
	min := slice[0]
	for _, val := range slice[1:] {
		if val < min {
			min = val
		}
	}
	return min, nil
}

// RoundFloat64 rounds a float64 to the nearest integer
func RoundFloat64(x float64) float64 {
	return math.Floor(x + 0.5)
}

// FloorFloat64 returns the floor of a float64 value
func FloorFloat64(x float64) float64 {
	return math.Floor(x)
}

// CeilFloat64 returns the ceiling of a float64 value
func CeilFloat64(x float64) float64 {
	return math.Ceil(x)
}

// AbsInt returns the absolute value of an integer
func AbsInt(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

// AbsFloat64 returns the absolute value of a float64
func AbsFloat64(x float64) float64 {
	return math.Abs(x)
}

// MinFloat64 returns the minimum of two float64 values
func MinFloat64(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}

// MaxFloat64 returns the maximum of two float64 values
func MaxFloat64(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}

// MaxInt returns the maximum of two integers
func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// MinInt returns the minimum of two integers
func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// ClampFloat64 clamps a value between min and max
func ClampFloat64(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampInt clamps an integer value between min and max
func ClampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ============================================================================
// Cell Manipulation Utilities
// ============================================================================

// WrapCells wraps a slice of cells to fit within a specified width
// Inserts cells containing '\n' wherever a linebreak should occur
// Preserves cell styles during wrapping
func WrapCells(cells []draw.Cell, width uint) []draw.Cell {
	if len(cells) == 0 {
		return cells
	}

	str := CellsToString(cells)
	wrapped := wordwrap.WrapString(str, width)
	wrappedRunes := []rune(wrapped)

	wrappedCells := make([]draw.Cell, 0, len(wrappedRunes))
	cellIndex := 0

	for _, r := range wrappedRunes {
		if r == '\n' {
			wrappedCells = append(wrappedCells, draw.Cell{
				Rune:  r,
				Style: styling.StyleClear,
			})
		} else {
			// Preserve style from original cell
			if cellIndex < len(cells) {
				wrappedCells = append(wrappedCells, draw.Cell{
					Rune:  r,
					Style: cells[cellIndex].Style,
				})
				cellIndex++
			} else {
				// Fallback if we run out of cells (shouldn't happen)
				wrappedCells = append(wrappedCells, draw.Cell{
					Rune:  r,
					Style: styling.StyleClear,
				})
			}
		}
	}

	return wrappedCells
}

// RunesToStyledCells converts a slice of runes to styled cells with a given style
// Useful for converting plain text to cells for rendering
func RunesToStyledCells(runes []rune, style styling.Style) []draw.Cell {
	cells := make([]draw.Cell, len(runes))
	for i, r := range runes {
		cells[i] = draw.Cell{
			Rune:  r,
			Style: style,
		}
	}
	return cells
}

// StringToStyledCells converts a string to styled cells with a given style
// Convenience wrapper around RunesToStyledCells
func StringToStyledCells(s string, style styling.Style) []draw.Cell {
	return RunesToStyledCells([]rune(s), style)
}

// CellsToString converts a slice of cells back to a plain string
// Useful for extracting text content from styled cells
func CellsToString(cells []draw.Cell) string {
	runes := make([]rune, len(cells))
	for i, cell := range cells {
		runes[i] = cell.Rune
	}
	return string(runes)
}

// TrimCells trims cells to fit within a specified width
// Adds ellipsis if the content was truncated
// Preserves cell styles
func TrimCells(cells []draw.Cell, width int) []draw.Cell {
	if len(cells) == 0 {
		return cells
	}

	s := CellsToString(cells)
	trimmed := TrimString(s, width)
	trimmedRunes := []rune(trimmed)

	newCells := make([]draw.Cell, 0, len(trimmedRunes))
	for i, r := range trimmedRunes {
		if i < len(cells) {
			newCells = append(newCells, draw.Cell{
				Rune:  r,
				Style: cells[i].Style,
			})
		} else {
			// Use last cell's style for ellipsis if needed
			style := styling.StyleClear
			if len(cells) > 0 {
				style = cells[len(cells)-1].Style
			}
			newCells = append(newCells, draw.Cell{
				Rune:  r,
				Style: style,
			})
		}
	}

	return newCells
}

// SplitCells splits cells by a delimiter rune
// Returns a slice of cell slices, one for each segment
// Useful for splitting text into lines or columns
func SplitCells(cells []draw.Cell, delimiter rune) [][]draw.Cell {
	splitCells := make([][]draw.Cell, 0)
	temp := make([]draw.Cell, 0)

	for _, cell := range cells {
		if cell.Rune == delimiter {
			if len(temp) > 0 {
				splitCells = append(splitCells, temp)
			}
			temp = make([]draw.Cell, 0)
		} else {
			temp = append(temp, cell)
		}
	}

	if len(temp) > 0 {
		splitCells = append(splitCells, temp)
	}

	return splitCells
}

// CellWithX represents a cell with its X coordinate
// Used for positioning cells with variable-width characters
type CellWithX struct {
	X    int        // X coordinate accounting for character width
	Cell draw.Cell // The cell itself
}

// BuildCellWithXArray builds an array of CellWithX from a slice of cells
// Calculates X positions accounting for variable-width characters (e.g., CJK)
// Essential for proper text alignment and rendering
func BuildCellWithXArray(cells []draw.Cell) []CellWithX {
	cellWithXArray := make([]CellWithX, len(cells))
	xPos := 0

	for i, cell := range cells {
		cellWithXArray[i] = CellWithX{
			X:    xPos,
			Cell: cell,
		}
		xPos += rw.RuneWidth(cell.Rune)
	}

	return cellWithXArray
}

// ============================================================================
// Additional Utility Functions (Project-Specific Enhancements)
// ============================================================================

// NormalizeFloat64Slice normalizes a float64 slice to 0-1 range
// Useful for scaling data for visualization
func NormalizeFloat64Slice(slice []float64) []float64 {
	if len(slice) == 0 {
		return slice
	}

	max, err := GetMaxFloat64FromSlice(slice)
	if err != nil || max == 0 {
		return slice
	}

	normalized := make([]float64, len(slice))
	for i, val := range slice {
		normalized[i] = val / max
	}

	return normalized
}

// ScaleFloat64Slice scales a float64 slice to fit within a target range
// Useful for mapping data values to display coordinates
func ScaleFloat64Slice(slice []float64, min, max float64) []float64 {
	if len(slice) == 0 {
		return slice
	}

	dataMin, errMin := GetMinFloat64FromSlice(slice)
	dataMax, errMax := GetMaxFloat64FromSlice(slice)

	if errMin != nil || errMax != nil || dataMax == dataMin {
		return slice
	}

	scale := (max - min) / (dataMax - dataMin)
	scaled := make([]float64, len(slice))

	for i, val := range slice {
		scaled[i] = min + (val-dataMin)*scale
	}

	return scaled
}

// RepeatCells repeats a cell or cell slice a specified number of times
// Useful for creating borders, separators, or padding
func RepeatCells(cell draw.Cell, count int) []draw.Cell {
	cells := make([]draw.Cell, count)
	for i := range cells {
		cells[i] = cell
	}
	return cells
}

// PadCells pads a cell slice to a target width with a padding cell
// Useful for aligning text or creating fixed-width displays
func PadCells(cells []draw.Cell, width int, padCell draw.Cell) []draw.Cell {
	currentWidth := 0
	for _, cell := range cells {
		currentWidth += rw.RuneWidth(cell.Rune)
	}

	if currentWidth >= width {
		return TrimCells(cells, width)
	}

	padded := make([]draw.Cell, len(cells))
	copy(padded, cells)

	for currentWidth < width {
		padded = append(padded, padCell)
		currentWidth += rw.RuneWidth(padCell.Rune)
	}

	return padded
}
