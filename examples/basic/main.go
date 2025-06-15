package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"

	"github.com/dmahlow/imgrid"
)

func main() {
	// Create a sample image (600x400 pixels)
	img := image.NewRGBA(image.Rect(0, 0, 600, 400))
	
	// Fill with gradient background
	for y := 0; y < 400; y++ {
		for x := 0; x < 600; x++ {
			r := uint8(x * 255 / 600)
			g := uint8(y * 255 / 400)
			b := uint8(128)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	
	// Add some colored rectangles
	draw.Draw(img, image.Rect(50, 50, 150, 150), &image.Uniform{color.RGBA{255, 0, 0, 255}}, image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(200, 100, 350, 200), &image.Uniform{color.RGBA{0, 255, 0, 255}}, image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(400, 200, 550, 350), &image.Uniform{color.RGBA{0, 0, 255, 255}}, image.Point{}, draw.Src)
	
	// Apply grid with default configuration
	config := imgrid.DefaultConfig()
	gridBytes, err := imgrid.AddGrid(img, config)
	if err != nil {
		log.Fatalf("Failed to add grid: %v", err)
	}
	
	// Save the result
	if err := os.WriteFile("sample_with_grid.png", gridBytes, 0644); err != nil {
		log.Fatalf("Failed to save image: %v", err)
	}
	
	fmt.Println("Grid applied successfully! Check sample_with_grid.png")
	
	// Demonstrate coordinate conversion
	fmt.Println("\nCoordinate conversion examples:")
	
	// Convert some cell numbers to pixel coordinates
	for _, cellNum := range []int{0, 1, 6, 7, 12, 13} {
		x, y, err := imgrid.CellToPixel(cellNum, 600, 100)
		if err != nil {
			log.Printf("Error converting cell %d: %v", cellNum, err)
			continue
		}
		fmt.Printf("Cell %d -> Pixel (%d, %d)\n", cellNum, x, y)
	}
	
	// Convert some pixel coordinates to cell numbers
	pixels := [][]int{{50, 50}, {150, 150}, {250, 250}, {450, 350}}
	for _, pixel := range pixels {
		cellNum := imgrid.PixelToCell(pixel[0], pixel[1], 600, 100)
		fmt.Printf("Pixel (%d, %d) -> Cell %d\n", pixel[0], pixel[1], cellNum)
	}
}