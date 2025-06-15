// Package imgrid provides utilities for overlaying coordinate grids on images.
// It creates numbered grid cells of configurable size and provides coordinate conversion functions.
package imgrid

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

// Config holds grid overlay configuration.
type Config struct {
	CellSize    int         // Size of each grid cell in pixels (default: 100)
	GridColor   color.Color // Color of grid lines (default: semi-transparent cyan)
	NumberColor color.Color // Color of cell numbers (default: white)
	NumberBG    color.Color // Background color for cell numbers (default: semi-transparent black)
	LineWidth   int         // Width of grid lines in pixels (default: 2)
	NumberScale int         // Scale factor for number size (default: 3)
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		CellSize:    100,
		GridColor:   color.RGBA{0, 255, 255, 100}, // Semi-transparent cyan
		NumberColor: color.RGBA{255, 255, 255, 255}, // White
		NumberBG:    color.RGBA{0, 0, 0, 200}, // Semi-transparent black
		LineWidth:   2,
		NumberScale: 3,
	}
}

// AddGrid overlays a numbered grid on the provided image using the given configuration.
// Returns the modified image as PNG bytes.
func AddGrid(img image.Image, config Config) ([]byte, error) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Create a new RGBA image to draw on
	overlay := image.NewRGBA(bounds)
	draw.Draw(overlay, bounds, img, bounds.Min, draw.Src)

	// Draw vertical lines
	for x := config.CellSize; x < width; x += config.CellSize {
		for y := 0; y < height; y++ {
			for i := 0; i < config.LineWidth && x-i >= 0; i++ {
				overlay.Set(x-i, y, config.GridColor)
			}
		}
	}

	// Draw horizontal lines
	for y := config.CellSize; y < height; y += config.CellSize {
		for x := 0; x < width; x++ {
			for i := 0; i < config.LineWidth && y-i >= 0; i++ {
				overlay.Set(x, y-i, config.GridColor)
			}
		}
	}

	// Add sequential numbers in center of each cell
	cellNumber := 0
	for gridY := 0; gridY*config.CellSize < height; gridY++ {
		for gridX := 0; gridX*config.CellSize < width; gridX++ {
			// Calculate center of the cell
			centerX := gridX*config.CellSize + config.CellSize/2
			centerY := gridY*config.CellSize + config.CellSize/2

			// Only draw if center is within bounds
			if centerX < width && centerY < height {
				drawLargeNumber(overlay, centerX, centerY, cellNumber, config)
			}
			cellNumber++
		}
	}

	// Encode to PNG bytes
	var buf bytes.Buffer
	if err := png.Encode(&buf, overlay); err != nil {
		return nil, fmt.Errorf("failed to encode image with grid: %v", err)
	}

	return buf.Bytes(), nil
}

// CellToPixel converts a cell number to pixel coordinates (center of the cell).
func CellToPixel(cellNumber int, imageWidth int, cellSize int) (int, int, error) {
	if cellNumber < 0 {
		return 0, 0, fmt.Errorf("invalid cell number: %d", cellNumber)
	}

	// Calculate columns per row based on image width
	columnsPerRow := imageWidth / cellSize
	if columnsPerRow == 0 {
		columnsPerRow = 1
	}

	// Convert cell number to grid coordinates
	gridX := cellNumber % columnsPerRow
	gridY := cellNumber / columnsPerRow

	// Calculate pixel coordinates (center of the cell)
	pixelX := gridX*cellSize + cellSize/2
	pixelY := gridY*cellSize + cellSize/2

	return pixelX, pixelY, nil
}

// PixelToCell converts pixel coordinates to the corresponding cell number.
func PixelToCell(x, y int, imageWidth int, cellSize int) int {
	columnsPerRow := imageWidth / cellSize
	if columnsPerRow == 0 {
		columnsPerRow = 1
	}

	gridX := x / cellSize
	gridY := y / cellSize

	return gridY*columnsPerRow + gridX
}

// getDigitPattern returns a 5x7 bitmap pattern for digits 0-9.
func getDigitPattern(digit rune) []string {
	patterns := map[rune][]string{
		'0': {
			" ### ",
			"#   #",
			"#   #",
			"#   #",
			"#   #",
			"#   #",
			" ### ",
		},
		'1': {
			"  #  ",
			" ##  ",
			"  #  ",
			"  #  ",
			"  #  ",
			"  #  ",
			"#####",
		},
		'2': {
			" ### ",
			"#   #",
			"    #",
			"   # ",
			"  #  ",
			" #   ",
			"#####",
		},
		'3': {
			" ### ",
			"#   #",
			"    #",
			"  ## ",
			"    #",
			"#   #",
			" ### ",
		},
		'4': {
			"   # ",
			"  ## ",
			" # # ",
			"#  # ",
			"#####",
			"   # ",
			"   # ",
		},
		'5': {
			"#####",
			"#    ",
			"#    ",
			"#### ",
			"    #",
			"#   #",
			" ### ",
		},
		'6': {
			" ### ",
			"#   #",
			"#    ",
			"#### ",
			"#   #",
			"#   #",
			" ### ",
		},
		'7': {
			"#####",
			"    #",
			"   # ",
			"  #  ",
			" #   ",
			" #   ",
			" #   ",
		},
		'8': {
			" ### ",
			"#   #",
			"#   #",
			" ### ",
			"#   #",
			"#   #",
			" ### ",
		},
		'9': {
			" ### ",
			"#   #",
			"#   #",
			" ####",
			"    #",
			"#   #",
			" ### ",
		},
	}

	if pattern, ok := patterns[digit]; ok {
		return pattern
	}
	return []string{} // Return empty pattern for unknown digits
}

// drawLargeNumber draws a number at the specified position with large, readable digits.
func drawLargeNumber(img draw.Image, x, y int, number int, config Config) {
	numStr := fmt.Sprintf("%d", number)

	// Size settings
	digitWidth := 5 * config.NumberScale
	digitHeight := 7 * config.NumberScale
	spacing := 2 * config.NumberScale
	padding := 2 * config.NumberScale

	// Calculate total width needed
	totalWidth := len(numStr)*digitWidth + (len(numStr)-1)*spacing + 2*padding
	totalHeight := digitHeight + 2*padding

	// Center the number block
	startX := x - totalWidth/2
	startY := y - totalHeight/2

	// Draw background rectangle
	for dx := 0; dx < totalWidth; dx++ {
		for dy := 0; dy < totalHeight; dy++ {
			px := startX + dx
			py := startY + dy
			if px >= 0 && py >= 0 && px < img.Bounds().Max.X && py < img.Bounds().Max.Y {
				img.Set(px, py, config.NumberBG)
			}
		}
	}

	// Draw each digit
	for i, digit := range numStr {
		pattern := getDigitPattern(digit)
		digitX := startX + padding + i*(digitWidth+spacing)
		digitY := startY + padding

		// Draw the pattern
		for row, line := range pattern {
			for col, char := range line {
				if char == '#' {
					// Draw a scaled block for each '#'
					for sx := 0; sx < config.NumberScale; sx++ {
						for sy := 0; sy < config.NumberScale; sy++ {
							px := digitX + col*config.NumberScale + sx
							py := digitY + row*config.NumberScale + sy
							if px >= 0 && py >= 0 && px < img.Bounds().Max.X && py < img.Bounds().Max.Y {
								img.Set(px, py, config.NumberColor)
							}
						}
					}
				}
			}
		}
	}
}