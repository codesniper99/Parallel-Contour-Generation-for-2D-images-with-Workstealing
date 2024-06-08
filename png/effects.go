// Package png allows for loading png images and applying
// image flitering effects on them.
package png

import (
	"image"
	"image/color"
)

func (img *Image) SetBackgroundBlack() {
	bounds := img.out.Bounds()
	for x := 0; x < bounds.Dx()-1; x++ {
		for y := 0; y < bounds.Dy()-1; y++ {
			img.out.Set(x, y, color.Black)
			img.out.Set(x+1, y, color.Black)
			img.out.Set(x, y+1, color.Black)
			img.out.Set(x+1, y+1, color.Black)
		}
	}
}

// processImage applies the Marching Squares algorithm to the given RGBA image.
func (img *Image) ProcessImage(threshold int, c color.Color) {
	bounds := img.out.Bounds()

	// Iterate through each 2x2 square in the image.
	for x := 0; x < bounds.Dx()-1; x++ {
		for y := 0; y < bounds.Dy()-1; y++ {
			// Get the intensity values of the four corners of the square.
			topLeft := intensity(img.in.At(x, y))
			topRight := intensity(img.in.At(x+1, y))
			bottomLeft := intensity(img.in.At(x, y+1))
			bottomRight := intensity(img.in.At(x+1, y+1))
			// Determine the contour value based on the intensity values.
			contourValue := calculateContourValue(topLeft, topRight, bottomLeft, bottomRight, threshold)

			// Draw the contour line based on the contour value.
			drawContourLine(img, x, y, contourValue, c)
		}
	}

	img.in = img.out
	img.out = image.NewRGBA64(img.Bounds)
}

// Grayscale applies a grayscale filtering effect to the image
func (img *Image) Grayscale() {

	// Bounds returns defines the dimensions of the image. Always
	// use the bounds Min and Max fields to get out the width
	// and height for the image
	bounds := img.out.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			//Returns the pixel (i.e., RGBA) value at a (x,y) position
			// Note: These get returned as int32 so based on the math you'll
			// be performing you'll need to do a conversion to float64(..)
			r, g, b, a := img.in.At(x, y).RGBA()

			//Note: The values for r,g,b,a for this assignment will range between [0, 65535].
			//For certain computations (i.e., convolution) the values might fall outside this
			// range so you need to clamp them between those values.
			greyC := clamp(float64(r+g+b) / 3)

			//Note: The values need to be stored back as uint16 (I know weird..but there's valid reasons
			// for this that I won't get into right now).
			img.out.Set(x, y, color.RGBA64{greyC, greyC, greyC, uint16(a)})
		}
	}
	img.in = img.out
	img.out = image.NewRGBA64(img.Bounds)
}

// Blur applies a blur to the image
func (img *Image) ApplyKernel(kernel []float64) {

	idx_array_x := []int{-1, 0, 1}
	idx_array_y := []int{-1, 0, 1}
	bounds := img.out.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var r_tmp, g_tmp, b_tmp float64
			_, _, _, a := img.in.At(x, y).RGBA()
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					idx := i*3 + j
					mul_val := kernel[idx]

					new_x := x + idx_array_x[j]
					new_y := y + idx_array_y[i]

					if new_x < bounds.Min.X || new_x >= bounds.Max.X {
						continue
					} else if new_y < bounds.Min.Y || new_y >= bounds.Max.Y {
						continue
					} else {

						r_orig, g_orig, b_orig, _ := img.in.At(new_x, new_y).RGBA()
						r_tmp += (mul_val) * float64(r_orig)
						g_tmp += (mul_val) * float64(g_orig)
						b_tmp += (mul_val) * float64(b_orig)

					}

				}
			}
			r_out_int := clamp(r_tmp)
			g_out_int := clamp(g_tmp)
			b_out_int := clamp(b_tmp)
			img.out.Set(x, y, color.RGBA64{uint16(r_out_int), uint16(g_out_int), uint16(b_out_int), uint16(a)})
		}
	}

	img.in = img.out
	img.out = image.NewRGBA64(img.Bounds)
}

// intensity calculates the intensity value of a pixel.
func intensity(c color.Color) int {
	r, g, b, _ := c.RGBA()
	r_out := clamp(float64(r))
	g_out := clamp(float64(g))
	b_out := clamp(float64(b))

	// Use a simple intensity calculation (you may need to adjust this based on your image characteristics).
	return int((uint16(r_out) + uint16(g_out) + uint16(b_out)) / 3)
}

// calculateContourValue determines the contour value based on the intensity values of the square corners.
func calculateContourValue(topLeft, topRight, bottomLeft, bottomRight int, threshold int) int {
	contourValue := 0

	if topLeft > threshold {
		contourValue += 1
	}
	if topRight > threshold {
		contourValue += 2
	}
	if bottomLeft > threshold {
		contourValue += 8
	}
	if bottomRight > threshold {
		contourValue += 4
	}

	return contourValue
}

// drawContourLine draws the contour line based on the contour value.
func drawContourLine(img *Image, x, y, contourValue int, c color.Color) {
	// Draw line segments based on the contour value.
	img.out.Set(x, y, color.Black)
	img.out.Set(x+1, y, color.Black)
	img.out.Set(x+1, y+1, color.Black)
	img.out.Set(x, y+1, color.Black)

	switch contourValue {
	case 1, 14:
		// Draw line from top left to top right
		img.out.Set(x, y, c)
		img.out.Set(x+1, y, c)
	case 2, 13:
		// Draw line from top right to bottom right
		img.out.Set(x+1, y, c)
		img.out.Set(x+1, y+1, c)
	case 4, 11:
		// Draw line from bottom right to bottom left
		img.out.Set(x+1, y+1, c)
		img.out.Set(x, y+1, c)
	case 8, 7:
		// Draw line from bottom left to top left
		img.out.Set(x, y+1, c)
		img.out.Set(x, y, c)
	case 3, 12:
		// Draw lines from top left to bottom right and top right to bottom left
		img.out.Set(x, y, c)
		img.out.Set(x+1, y+1, c)
	case 9, 6:
		// Draw lines from top left to bottom left and top right to bottom right
		img.out.Set(x, y, c)
		img.out.Set(x+1, y+1, c)
	case 5, 10:
		// Draw lines from top right to top left and bottom right to bottom left
		img.out.Set(x+1, y, c)
		img.out.Set(x, y+1, c)
	default:
		// No contour line
	}
}
