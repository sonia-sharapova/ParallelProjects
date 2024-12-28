package scheduler

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func RunSequential(config Config) {
	// Open effects file
	effectsPathFile := "../data/effects.txt"
	effectsFile, err := os.Open(effectsPathFile)
	if err != nil {
		panic("Failed to open effects file")
	}
	defer effectsFile.Close()

	// Split data_dir argument by "+"
	dataDirs := strings.Split(config.DataDirs, "+")

	// JSON decoder
	reader := json.NewDecoder(effectsFile)
	for reader.More() {
		var effect struct {
			InPath  string   `json:"inPath"`
			OutPath string   `json:"outPath"`
			Effects []string `json:"effects"`
		}
		err := reader.Decode(&effect)
		if err != nil {
			panic("Failed to decode JSON")
		}

		// Process images in each specified directory
		for _, dir := range dataDirs {
			// Construct full input path based on each directory
			inFilePath := filepath.Join("../data/in", dir, effect.InPath)
			imgFile, err := os.Open(inFilePath)
			if err != nil {
				fmt.Printf("Failed to open image file %s: %v\n", inFilePath, err)
				continue // Skip to next image if opening fails
			}
			defer imgFile.Close()

			img, _, err := image.Decode(imgFile)
			if err != nil {
				fmt.Printf("Failed to decode image %s: %v\n", effect.InPath, err)
				continue
			}
			outImg := img

			// Apply each effect in sequence
			for _, ef := range effect.Effects {
				switch ef {
				case "S":
					outImg = ApplyKernel(outImg, []float64{0, -1, 0, -1, 5, -1, 0, -1, 0})
				case "E":
					outImg = ApplyKernel(outImg, []float64{-1, -1, -1, -1, 8, -1, -1, -1, -1})
				case "B":
					outImg = ApplyKernel(outImg, []float64{1.0 / 9, 1.0 / 9, 1.0 / 9, 1.0 / 9, 1.0 / 9, 1.0 / 9, 1.0 / 9, 1.0 / 9, 1.0 / 9})
				case "G":
					outImg = ApplyGrayscale(outImg)
				}
			}

			// Save the processed image with data_dir prefix in outPath
			outFilePath := filepath.Join("../data/out", fmt.Sprintf("%s_%s", dir, effect.OutPath))

			outFile, err := os.Create(outFilePath)
			if err != nil {
				fmt.Printf("Failed to create output file %s: %v\n", outFilePath, err)
				continue
			}
			defer outFile.Close()

			err = png.Encode(outFile, outImg)
			if err != nil {
				fmt.Printf("Failed to encode image %s: %v\n", outFilePath, err)
			}
		}
	}
}

// ApplyKernel applies a convolution kernel to an image and returns the processed image.
func ApplyKernel(img image.Image, kernel []float64) image.Image {
	bounds := img.Bounds()
	outImg := image.NewRGBA(bounds)
	kernelSize := 3 // We're using a fixed 3x3 kernel
	offset := kernelSize / 2

	// Iterate over each pixel in the image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var rSum, gSum, bSum float64

			// Convolve kernel with the 3x3 neighborhood around (x, y)
			for ky := -offset; ky <= offset; ky++ {
				for kx := -offset; kx <= offset; kx++ {
					// Determine neighboring pixel coordinates
					nx, ny := x+kx, y+ky

					// Zero-padding: Ignore out-of-bounds pixels
					if nx >= bounds.Min.X && nx < bounds.Max.X && ny >= bounds.Min.Y && ny < bounds.Max.Y {
						// Get color components from neighboring pixel
						r, g, b, _ := img.At(nx, ny).RGBA()
						weight := kernel[(ky+offset)*kernelSize+(kx+offset)]

						// Accumulate weighted sum for each color channel
						rSum += float64(r>>8) * weight
						gSum += float64(g>>8) * weight
						bSum += float64(b>>8) * weight
					}
				}
			}

			// Clamp color values to valid range [0, 255] and set pixel in output image
			outImg.Set(x, y, color.RGBA{
				R: clampToUint8(rSum),
				G: clampToUint8(gSum),
				B: clampToUint8(bSum),
				A: 255,
			})
		}
	}
	return outImg
}

// Helper function to clamp values to uint8 range [0, 255]
func clampToUint8(value float64) uint8 {
	if value < 0 {
		return 0
	} else if value > 255 {
		return 255
	}
	return uint8(value)
}

// ApplyGrayscale applies grayscale effect to an image and returns the modified image.
func ApplyGrayscale(img image.Image) image.Image {
	bounds := img.Bounds()
	grayImg := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			avg := uint8((r + g + b) / 3 / 256) // Convert to uint8
			grayImg.Set(x, y, color.RGBA{avg, avg, avg, 255})
		}
	}
	return grayImg
}
