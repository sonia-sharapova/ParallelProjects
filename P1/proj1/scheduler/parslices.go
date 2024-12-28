/*
Parallelize Each Image
Parallelize the processing of individual images
- Assume that only one image is processed at a time.
*/
package scheduler

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// max returns the larger of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func RunParallelSlices(config Config) {
	// Split the data directories by "+" and process each one
	dataDirs := strings.Split(config.DataDirs, "+")

	// Create task queue
	queue := &TaskQueue{}

	// Populate the queue with tasks from each specified directory
	for _, dir := range dataDirs {
		effectsPathFile := "../data/effects.txt"
		effectsFile, err := os.Open(effectsPathFile)
		if err != nil {
			panic("Failed to open effects file")
		}
		defer effectsFile.Close()

		// JSON decoder
		decoder := json.NewDecoder(effectsFile)
		for decoder.More() {
			var effect struct {
				InPath  string   `json:"inPath"`
				OutPath string   `json:"outPath"`
				Effects []string `json:"effects"`
			}
			err := decoder.Decode(&effect)
			if err != nil {
				panic("Failed to decode JSON")
			}
			// Prefix output path with the current directory name
			task := &Task{
				inPath:  filepath.Join("../data/in", dir, effect.InPath),
				outPath: filepath.Join("../data/out", fmt.Sprintf("%s_%s", dir, effect.OutPath)),
				effects: effect.Effects,
			}
			queue.Enqueue(task)
		}
	}

	// Iterate over tasks in the queue. For each image, create goroutines to process slices.
	for {
		task := queue.Dequeue()
		if task == nil {
			break
		}

		// Start timer for the parallel section of this image processing
		//startParallel := time.Now()

		processImageInSlices(task, config.ThreadCount)

		// End timer and calculate duration for this image
		//parallelDuration := time.Since(startParallel).Seconds()
		//fmt.Print(parallelDuration, "\n")

	}
}

// Process an image in parallel slices
func processImageInSlices(task *Task, threadCount int) {
	// Load the image
	imgFile, err := os.Open(task.inPath)
	if err != nil {
		fmt.Printf("Failed to open image file %s: %v\n", task.inPath, err)
		return
	}
	defer imgFile.Close()

	// Decode the input image
	decodedImg, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Printf("Failed to decode image %s: %v\n", task.inPath, err)
		return
	}

	// Convert to *image.RGBA if necessary
	var inImg *image.RGBA
	if rgba, ok := decodedImg.(*image.RGBA); ok {
		inImg = rgba
	} else {
		// Convert to RGBA
		bounds := decodedImg.Bounds()
		inImg = image.NewRGBA(bounds)
		draw.Draw(inImg, bounds, decodedImg, bounds.Min, draw.Src)
	}

	// Prepare an output buffer
	outImg := image.NewRGBA(inImg.Bounds())

	height := inImg.Bounds().Dy()
	sliceHeight := height / threadCount

	// Apply each effect in sequence, reusing the output buffer by swapping pointers
	for _, effect := range task.effects {
		var wg sync.WaitGroup

		for i := 0; i < threadCount; i++ {
			// Define slice bounds, with overlap
			startY := i * sliceHeight
			endY := startY + sliceHeight
			if i == threadCount-1 {
				endY = height
			}
			overlapStartY := max(0, startY-1)
			overlapEndY := min(height, endY+1)

			wg.Add(1)
			go func(start, end, overlapStart, overlapEnd int) {
				defer wg.Done()
				sliceWorker(inImg, outImg, start, end, overlapStart, overlapEnd, effect)
			}(startY, endY, overlapStartY, overlapEndY)
		}

		wg.Wait()

		// Swap input and output images for the next effect
		inImg, outImg = outImg, inImg
	}

	// Save the final processed image (now in inImg after the last swap)
	outFile, err := os.Create(task.outPath)
	if err != nil {
		fmt.Printf("Failed to create output file %s: %v\n", task.outPath, err)
		return
	}
	defer outFile.Close()

	err = png.Encode(outFile, inImg)
	if err != nil {
		fmt.Printf("Failed to encode image %s: %v\n", task.outPath, err)
	}
}

/*
Worker function for processing a slice of the image
Applies effects sequentially for each slice and handles slice boundaries
*/
func sliceWorker(inImg *image.RGBA, outImg *image.RGBA, startY, endY, overlapStartY, overlapEndY int, effect string) {
	bounds := inImg.Bounds()
	width := bounds.Dx()

	for y := startY; y < endY; y++ {
		for x := 0; x < width; x++ {
			switch effect {
			case "S":
				outImg.SetRGBA(x, y, applyKernelDirect(inImg, x, y, overlapStartY, overlapEndY, []float64{0, -1, 0, -1, 5, -1, 0, -1, 0}))
			case "E":
				outImg.SetRGBA(x, y, applyKernelDirect(inImg, x, y, overlapStartY, overlapEndY, []float64{-1, -1, -1, -1, 8, -1, -1, -1, -1}))
			case "B":
				outImg.SetRGBA(x, y, applyKernelDirect(inImg, x, y, overlapStartY, overlapEndY, []float64{1.0 / 9, 1.0 / 9, 1.0 / 9, 1.0 / 9, 1.0 / 9, 1.0 / 9, 1.0 / 9, 1.0 / 9, 1.0 / 9}))
			case "G":
				outImg.SetRGBA(x, y, applyGrayscaleDirect(inImg, x, y))
			}
		}
	}
}

// Helper function for handling Overlapping Boundaries:
func applyKernelDirect(img *image.RGBA, x, y, overlapStartY, overlapEndY int, kernel []float64) color.RGBA {
	var rSum, gSum, bSum float64
	offset := 1 // 3x3 kernel

	for ky := -offset; ky <= offset; ky++ {
		for kx := -offset; kx <= offset; kx++ {
			nx, ny := x+kx, y+ky

			// Ensure ny is within bounds
			if nx >= img.Bounds().Min.X && nx < img.Bounds().Max.X && ny >= overlapStartY && ny < overlapEndY {
				i := img.PixOffset(nx, ny)
				r := float64(img.Pix[i])
				g := float64(img.Pix[i+1])
				b := float64(img.Pix[i+2])
				weight := kernel[(ky+offset)*3+(kx+offset)]
				rSum += r * weight
				gSum += g * weight
				bSum += b * weight
			}
		}
	}

	return color.RGBA{
		R: clampToUint8(rSum),
		G: clampToUint8(gSum),
		B: clampToUint8(bSum),
		A: 255,
	}
}

// Helper function to apply grayscale at a specific pixel
func applyGrayscaleDirect(img *image.RGBA, x, y int) color.RGBA {
	i := img.PixOffset(x, y)
	r := img.Pix[i]
	g := img.Pix[i+1]
	b := img.Pix[i+2]
	avg := uint8((int(r) + int(g) + int(b)) / 3)
	return color.RGBA{avg, avg, avg, 255}
}
