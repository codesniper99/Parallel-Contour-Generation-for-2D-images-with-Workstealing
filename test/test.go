package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"proj3/png"
)

// Point represents a point in 2D space.
type Point struct {
	X, Y int
}

// Cell represents a square cell in the image.
type Cell struct {
	TopLeft, TopRight, BottomLeft, BottomRight Point
}

func main() {
	// Load your RGBA image (replace this with your image loading logic).
	//Loads the png image and returns the image or an error
	pngImg, err := png.Load("../data/in/IMG_1.png")

	if err != nil {
		panic(err)
	}

	//Performs a grayscale filtering effect on the image
	pngImg.Grayscale()

	//Saves the image to a new file
	err = pngImg.Save("../data/in/IMG_1_Grey.png")

	//Checks to see if there were any errors when saving.
	if err != nil {
		panic(err)
	}
	pngImg, err = png.Load("../data/in/IMG_1_Grey.png")

	if err != nil {
		panic(err)
	}

	pngImg.SetBackgroundBlack()

	threshold := 5000
	orange := color.RGBA{
		R: 255, // Red component
		G: 165, // Green component
		B: 0,   // Blue component
		A: 255, // Alpha component (255 means fully opaque)
	}

	pngImg.ProcessImage(threshold, orange)

	// threshold = 5000

	// pngImg.ProcessImage(threshold, color.White)
	err = pngImg.Save("../data/out/IMG_1_OUT.png")

	//Checks to see if there were any errors when saving.
	if err != nil {
		panic(err)
	}
}

// loadImage loads an image from a file and returns its RGBA representation.
func loadImage(filename string) *image.RGBA {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening image file:", err)
		os.Exit(1)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		os.Exit(1)
	}

	rgba := image.NewRGBA(img.Bounds())
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	return rgba
}

// clamp will clamp the comp parameter to zero if it is less than zero or to 65535 if the comp parameter
// is greater than 65535.
func clamp(comp float64) uint16 {
	return uint16(math.Min(65535, math.Max(0, comp)))
}
