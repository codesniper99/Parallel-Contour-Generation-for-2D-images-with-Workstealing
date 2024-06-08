package png

import (
	"image"
	"image/color"
	"proj3/customwg"
)

// --------------------------------- PARALLEL

func swapImageInOut(img *Image) {
	img.in = img.out
	img.out = image.NewRGBA64(img.Bounds)
}

// Grayscale applies a grayscale filtering effect to the image
func (img *Image) GrayscaleParallely(threadCount int, imageName string) {

	bounds := img.Bounds
	YSideLength := (bounds.Max.Y - bounds.Min.Y) / threadCount
	var chunkYMin, chunkYMax int
	chunkYMin = -1
	chunkYMax = -1
	var wg3 customwg.CustomWaitGroup = *customwg.NewWaitGroup()
	for i := 0; i < threadCount; i++ {
		wg3.Add(1)

		chunkYMin = chunkYMax + 1
		chunkYMax = chunkYMin + YSideLength

		if chunkYMin > bounds.Max.Y {
			chunkYMin = bounds.Max.Y
		}
		if chunkYMax > bounds.Max.Y {
			chunkYMax = bounds.Max.Y
		}
		chunkMax := image.Point{bounds.Max.X, chunkYMax}
		chunkMin := image.Point{bounds.Min.X, chunkYMin}

		chunkBounds := image.Rectangle{chunkMin, chunkMax}
		go GreyscaleSlice(img, chunkBounds, &wg3)
	}

	wg3.Wait()
	swapImageInOut(img)

}

func GreyscaleSlice(img *Image, bounds image.Rectangle, wg3 *customwg.CustomWaitGroup) {

	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		for x := bounds.Min.X; x <= bounds.Max.X; x++ {

			r, g, b, a := img.in.At(x, y).RGBA()
			greyC := clamp(float64(r+g+b) / 3)
			img.out.Set(x, y, color.RGBA64{greyC, greyC, greyC, uint16(a)})
		}
	}
	wg3.Done()
}

func WorkOnSlice(img *Image, kernel []float64, bounds image.Rectangle, wg3 *customwg.CustomWaitGroup) {

	idx_array_x := []int{-1, 0, 1}
	idx_array_y := []int{-1, 0, 1}

	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		for x := bounds.Min.X; x <= bounds.Max.X; x++ {
			var r_tmp, g_tmp, b_tmp float64
			_, _, _, a := img.in.At(x, y).RGBA()
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					idx := i*3 + j
					mul_val := kernel[idx]

					new_x := x + idx_array_x[j]
					new_y := y + idx_array_y[i]

					if new_x < img.Bounds.Min.X || new_x > img.Bounds.Max.X {
						continue
					} else if new_y < img.Bounds.Min.Y || new_y > img.Bounds.Max.Y {
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
	wg3.Done()
}

// Blur applies a blur to the image
func (img *Image) ApplyKernelParallely(kernel []float64, threadCount int, imageName string) {

	bounds := img.Bounds
	YSideLength := (bounds.Max.Y - bounds.Min.Y) / threadCount
	var chunkYMin, chunkYMax int
	chunkYMin = -1
	chunkYMax = -1
	var wg3 customwg.CustomWaitGroup = *customwg.NewWaitGroup()
	for i := 0; i < threadCount; i++ {
		wg3.Add(1)
		//start := time.Now()
		chunkYMin = chunkYMax + 1
		chunkYMax = chunkYMin + YSideLength

		if chunkYMin > bounds.Max.Y {
			chunkYMin = bounds.Max.Y
		}
		if chunkYMax > bounds.Max.Y {
			chunkYMax = bounds.Max.Y
		}
		chunkMax := image.Point{bounds.Max.X, chunkYMax}
		chunkMin := image.Point{bounds.Min.X, chunkYMin}

		chunkBounds := image.Rectangle{chunkMin, chunkMax}
		go WorkOnSlice(img, kernel, chunkBounds, &wg3)
	}

	wg3.Wait()

	swapImageInOut(img)
}

func (img *Image) ProcessImageParallely(threshold int, threadCount int, c color.Color) {

	bounds := img.Bounds
	YSideLength := (bounds.Max.Y - bounds.Min.Y) / threadCount
	var chunkYMin, chunkYMax int
	chunkYMin = -1
	chunkYMax = -1
	var wg3 customwg.CustomWaitGroup = *customwg.NewWaitGroup()
	for i := 0; i < threadCount; i++ {
		wg3.Add(1)
		chunkYMin = chunkYMax + 1
		chunkYMax = chunkYMin + YSideLength

		if chunkYMin > bounds.Max.Y {
			chunkYMin = bounds.Max.Y
		}
		if chunkYMax > bounds.Max.Y {
			chunkYMax = bounds.Max.Y
		}
		chunkMax := image.Point{bounds.Max.X, chunkYMax}
		chunkMin := image.Point{bounds.Min.X, chunkYMin}

		chunkBounds := image.Rectangle{chunkMin, chunkMax}
		go WorkOnSliceMSParallel(img, c, threshold, chunkBounds, &wg3)
	}

	wg3.Wait()

	swapImageInOut(img)

}

func WorkOnSliceMSParallel(img *Image, c color.Color, threshold int, bounds image.Rectangle, wg3 *customwg.CustomWaitGroup) {
	// Iterate through each 2x2 square in the image.
	//fmt.Println("Bounds are ", bounds)
	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		for x := bounds.Min.X; x <= bounds.Max.X; x++ {
			// Get the intensity values of the four corners of the square.
			topLeft := intensity(img.in.At(x, y))
			topRight := intensity(img.in.At(x+1, y))
			bottomLeft := intensity(img.in.At(x, y+1))
			bottomRight := intensity(img.in.At(x+1, y+1))
			// Determine the contour value based on the intensity values.
			contourValue := calculateContourValue(topLeft, topRight, bottomLeft, bottomRight, threshold)

			// Draw the contour line based on the contour value.
			//fmt.Println(x, " ", y, " ", contourValue)
			drawContourLine(img, x, y, contourValue, c)
		}
	}
	wg3.Done()
}
