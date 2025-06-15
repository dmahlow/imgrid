# imgrid

A Go library for overlaying numbered coordinate grids on images.

## Installation

```bash
go get github.com/danielkesler/imgrid
```

## Features

- Overlay numbered grid cells on any image
- Configurable cell size, colors, and line width
- Convert between cell numbers and pixel coordinates
- PNG output support

## Usage

### Basic Grid Overlay

```go
package main

import (
    "image"
    "os"
    "github.com/danielkesler/imgrid"
)

func main() {
    // Load an image
    file, _ := os.Open("input.png")
    defer file.Close()
    img, _, _ := image.Decode(file)
    
    // Apply grid with default settings
    config := imgrid.DefaultConfig()
    gridBytes, _ := imgrid.AddGrid(img, config)
    
    // Save result
    os.WriteFile("output_grid.png", gridBytes, 0644)
}
```

### Custom Configuration

```go
config := imgrid.Config{
    CellSize:    50,  // 50x50 pixel cells
    GridColor:   color.RGBA{255, 0, 0, 128}, // Semi-transparent red
    NumberColor: color.RGBA{255, 255, 255, 255}, // White numbers
    NumberBG:    color.RGBA{0, 0, 0, 200}, // Black number background
    LineWidth:   1,  // 1-pixel wide lines
    NumberScale: 2,  // Smaller numbers
}
gridBytes, err := imgrid.AddGrid(img, config)
```

### Coordinate Conversion

```go
// Convert cell number to pixel coordinates
x, y, err := imgrid.CellToPixel(42, imageWidth, 100)
// Returns center coordinates of cell 42

// Convert pixel coordinates to cell number
cellNum := imgrid.PixelToCell(250, 150, imageWidth, 100)
// Returns cell number containing pixel (250, 150)
```

## API Reference

### Types

#### Config
```go
type Config struct {
    CellSize    int         // Size of each grid cell in pixels
    GridColor   color.Color // Color of grid lines
    NumberColor color.Color // Color of cell numbers
    NumberBG    color.Color // Background color for cell numbers
    LineWidth   int         // Width of grid lines in pixels
    NumberScale int         // Scale factor for number size
}
```

### Functions

#### DefaultConfig() Config
Returns a Config with sensible defaults:
- CellSize: 100 pixels
- GridColor: Semi-transparent cyan
- NumberColor: White
- NumberBG: Semi-transparent black
- LineWidth: 2 pixels
- NumberScale: 3

#### AddGrid(img image.Image, config Config) ([]byte, error)
Overlays a numbered grid on the provided image. Returns PNG-encoded bytes.

#### CellToPixel(cellNumber int, imageWidth int, cellSize int) (int, int, error)
Converts a cell number to pixel coordinates (center of the cell).

#### PixelToCell(x, y int, imageWidth int, cellSize int) int
Converts pixel coordinates to the corresponding cell number.

## Grid Layout

Cells are numbered sequentially starting from 0, left-to-right, top-to-bottom:

```
+-----+-----+-----+
|  0  |  1  |  2  |
+-----+-----+-----+
|  3  |  4  |  5  |
+-----+-----+-----+
|  6  |  7  |  8  |
+-----+-----+-----+
```

## Examples

See the `examples/` directory for complete working examples.