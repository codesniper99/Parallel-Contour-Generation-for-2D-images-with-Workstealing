package scheduler

import (
	"encoding/json"
	"image/color"
	"os"
	messagepackage "proj3/messagePackage"
	"proj3/png"
)

func ReadEffects() (effects []messagepackage.Message) {
	file, err := os.Open("./data/effects.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := json.NewDecoder(file)

	for {
		var msg messagepackage.Message
		err := reader.Decode(&msg)
		if err != nil {

			return effects
		}
		effects = append(effects, msg)
	}
}

func executeSequentialEffect(config messagepackage.Config, effects messagepackage.Message) {
	in_path := effects.InPath
	out_path := effects.OutPath

	pngImg, err := png.Load("./data/in/" + in_path)
	if err != nil {
		panic(err)
	}
	for _, j := range effects.Effects {

		if j == "S" {
			kernel := []float64{0, -1, 0, -1, 5, -1, 0, -1, 0}
			pngImg.ApplyKernel(kernel)

		} else if j == "E" {
			kernel := []float64{-1, -1, -1, -1, 8, -1, -1, -1, -1}
			pngImg.ApplyKernel(kernel)

		} else if j == "B" {
			kernel := []float64{1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0}
			pngImg.ApplyKernel(kernel)

		} else if j == "MS" {
			orange := color.RGBA{
				R: 255, // Red component
				G: 165, // Green component
				B: 0,   // Blue component
				A: 255, // Alpha component (255 means fully opaque)
			}
			pngImg.ProcessImage(config.Threshold, orange)
		} else {
			pngImg.Grayscale()
		}

	}

	err = pngImg.Save("./data/out/" + out_path)

	if err != nil {
		panic(err)
	}
}

func executeEffect(inputDir string, outputPrefix string, effects messagepackage.Message, mode string) {
	in_path := effects.InPath
	out_path := effects.OutPath

	pngImg, err := png.Load("../data/in/" + inputDir + in_path)
	if err != nil {
		panic(err)
	}
	for _, j := range effects.Effects {

		if j == "S" {
			kernel := []float64{0, -1, 0, -1, 5, -1, 0, -1, 0}
			pngImg.ApplyKernel(kernel)

		} else if j == "E" {
			kernel := []float64{-1, -1, -1, -1, 8, -1, -1, -1, -1}
			pngImg.ApplyKernel(kernel)

		} else if j == "B" {
			kernel := []float64{1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0}
			pngImg.ApplyKernel(kernel)

		} else {
			pngImg.Grayscale()
		}

	}

	err = pngImg.Save("../data/out/" + outputPrefix + out_path)

	if err != nil {
		panic(err)
	}
}

func executeParallelWorkStealingEffects(effects messagepackage.Message, threshold int, threadCount int) { //, timeForSaving *float64, timeForParallelParts *float64) {
	in_path := effects.InPath
	out_path := effects.OutPath

	pngImg, err := png.Load("./data/in/" + in_path)
	if err != nil {
		panic(err)
	}

	for _, j := range effects.Effects {
		if j == "S" {
			kernel := []float64{0, -1, 0, -1, 5, -1, 0, -1, 0}
			pngImg.ApplyKernel(kernel)

		} else if j == "E" {
			kernel := []float64{-1, -1, -1, -1, 8, -1, -1, -1, -1}
			pngImg.ApplyKernel(kernel)

		} else if j == "B" {
			kernel := []float64{1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0}
			pngImg.ApplyKernel(kernel)

		} else if j == "MS" {
			green := color.RGBA{
				R: 0,   // Red component
				G: 255, // Green component
				B: 0,   // Blue component
				A: 255, // Alpha component (255 means fully opaque)
			}
			yellow := color.RGBA{
				R: 0,   // Red component
				G: 255, // Green component
				B: 255, // Blue component
				A: 255, // Alpha component (255 means fully opaque)
			}
			blue := color.RGBA{
				R: 0,   // Red component
				G: 0,   // Green component
				B: 255, // Blue component
				A: 255, // Alpha component (255 means fully opaque)
			}
			red := color.RGBA{
				R: 255, // Red component
				G: 0,   // Green component
				B: 0,   // Blue component
				A: 255, // Alpha component (255 means fully opaque)
			}
			if threadCount == 2 {
				pngImg.ProcessImage(threshold, green)
			} else if threadCount == 4 {
				pngImg.ProcessImage(threshold, yellow)
			} else if threadCount == 6 {
				pngImg.ProcessImage(threshold, red)
			} else {
				pngImg.ProcessImage(threshold, blue)
			}
		} else {
			pngImg.Grayscale()
		}

	}
	err = pngImg.Save("./data/out/" + out_path)

	if err != nil {
		panic(err)
	}

}

func executeParallelChunkEffects(effects messagepackage.Message, threshold int, threadCount int) { //, timeForSaving *float64, timeForParallelParts *float64) {
	in_path := effects.InPath
	out_path := effects.OutPath

	pngImg, err := png.Load("./data/in/" + in_path)
	if err != nil {
		panic(err)
	}
	green := color.RGBA{
		R: 0,   // Red component
		G: 255, // Green component
		B: 0,   // Blue component
		A: 255, // Alpha component (255 means fully opaque)
	}
	yellow := color.RGBA{
		R: 0,   // Red component
		G: 255, // Green component
		B: 255, // Blue component
		A: 255, // Alpha component (255 means fully opaque)
	}
	blue := color.RGBA{
		R: 0,   // Red component
		G: 0,   // Green component
		B: 255, // Blue component
		A: 255, // Alpha component (255 means fully opaque)
	}
	red := color.RGBA{
		R: 255, // Red component
		G: 0,   // Green component
		B: 0,   // Blue component
		A: 255, // Alpha component (255 means fully opaque)
	}
	for _, j := range effects.Effects {
		if j == "S" {
			kernel := []float64{0, -1, 0, -1, 5, -1, 0, -1, 0}
			pngImg.ApplyKernel(kernel)

		} else if j == "E" {
			kernel := []float64{-1, -1, -1, -1, 8, -1, -1, -1, -1}
			pngImg.ApplyKernel(kernel)

		} else if j == "B" {
			kernel := []float64{1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0}
			pngImg.ApplyKernel(kernel)

		} else if j == "MS" {

			if threadCount == 2 {
				pngImg.ProcessImageParallely(threshold, threadCount, green)
			} else if threadCount == 4 {
				pngImg.ProcessImageParallely(threshold, threadCount, yellow)
			} else if threadCount == 6 {
				pngImg.ProcessImageParallely(threshold, threadCount, red)
			} else {
				pngImg.ProcessImageParallely(threshold, threadCount, blue)
			}
		} else {
			pngImg.GrayscaleParallely(threadCount, effects.InPath)
		}

	}
	err = pngImg.Save("./data/out/" + out_path)

	if err != nil {
		panic(err)
	}

}
