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
	// Create a larger sample image (800x600 pixels)
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))
	
	// Fill with white background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.Point{}, draw.Src)
	
	// Add some shapes for visual reference
	drawCircle(img, 200, 150, 50, color.RGBA{255, 100, 100, 255})
	drawCircle(img, 400, 300, 75, color.RGBA{100, 255, 100, 255})
	drawCircle(img, 600, 450, 60, color.RGBA{100, 100, 255, 255})
	
	// Custom configuration with smaller cells and different colors
	config1 := imgrid.Config{
		CellSize:    50,
		GridColor:   color.RGBA{255, 0, 0, 150}, // Semi-transparent red
		NumberColor: color.RGBA{255, 255, 255, 255},
		NumberBG:    color.RGBA{255, 0, 0, 200}, // Red background
		LineWidth:   1,
		NumberScale: 2,
	}
	
	gridBytes1, err := imgrid.AddGrid(img, config1)
	if err != nil {
		log.Fatalf("Failed to add grid: %v", err)
	}
	
	if err := os.WriteFile("custom_50px_grid.png", gridBytes1, 0644); err != nil {
		log.Fatalf("Failed to save image: %v", err)
	}
	
	// Another custom configuration with larger cells
	config2 := imgrid.Config{
		CellSize:    150,
		GridColor:   color.RGBA{0, 255, 0, 100}, // Semi-transparent green
		NumberColor: color.RGBA{0, 0, 0, 255}, // Black numbers
		NumberBG:    color.RGBA{255, 255, 255, 200}, // White background
		LineWidth:   3,
		NumberScale: 4,
	}
	
	gridBytes2, err := imgrid.AddGrid(img, config2)
	if err != nil {
		log.Fatalf("Failed to add grid: %v", err)
	}
	
	if err := os.WriteFile("custom_150px_grid.png", gridBytes2, 0644); err != nil {
		log.Fatalf("Failed to save image: %v", err)
	}
	
	fmt.Println("Custom grids created!")
	fmt.Println("- custom_50px_grid.png (50x50 cells, red)")
	fmt.Println("- custom_150px_grid.png (150x150 cells, green)")
	
	// Demonstrate coordinate conversion with different cell sizes
	fmt.Println("\nCoordinate conversion with 50px cells:")
	for _, cellNum := range []int{0, 16, 32} {
		x, y, err := imgrid.CellToPixel(cellNum, 800, 50)
		if err != nil {
			continue
		}
		fmt.Printf("Cell %d -> Pixel (%d, %d)\n", cellNum, x, y)
	}
	
	fmt.Println("\nCoordinate conversion with 150px cells:")
	for _, cellNum := range []int{0, 1, 2, 5, 6, 7} {
		x, y, err := imgrid.CellToPixel(cellNum, 800, 150)
		if err != nil {
			continue
		}
		fmt.Printf("Cell %d -> Pixel (%d, %d)\n", cellNum, x, y)
	}
}

// drawCircle draws a filled circle on the image
func drawCircle(img draw.Image, centerX, centerY, radius int, col color.Color) {
	for y := centerY - radius; y <= centerY + radius; y++ {
		for x := centerX - radius; x <= centerX + radius; x++ {
			dx := x - centerX
			dy := y - centerY
			if dx*dx + dy*dy <= radius*radius {
				if x >= 0 && y >= 0 && x < img.Bounds().Max.X && y < img.Bounds().Max.Y {
					img.Set(x, y, col)
				}
			}
		}
	}
}