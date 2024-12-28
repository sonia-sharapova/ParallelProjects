/*
Process multiple images in parallel.
Each individual image is handled by only one thread.
*/
package scheduler

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Hold each image processing job's information
type Task struct {
	inPath  string
	outPath string
	effects []string
}

// Queue with enqueue and dequeue methods
type TaskQueue struct {
	tasks []*Task
}

func (q *TaskQueue) Enqueue(task *Task) {
	q.tasks = append(q.tasks, task)
}

func (q *TaskQueue) Dequeue() *Task {
	if len(q.tasks) == 0 {
		return nil
	}
	task := q.tasks[0]
	q.tasks = q.tasks[1:]
	return task
}

/*
TAS lock implementation (to safeguard accesses to the queue)
Items can only be taken out of the queue by a Go routine that holds the lock.
*/
type TASLock struct {
	state int32
}

func (lock *TASLock) Lock() {
	for !atomic.CompareAndSwapInt32(&lock.state, 0, 1) {
		// Busy wait until the lock is acquired
	}
}

func (lock *TASLock) Unlock() {
	atomic.StoreInt32(&lock.state, 0)
}

/*
RunParallelFiles Function
This function populates the task queue, spawns goroutines, and uses the TAS lock to synchronize access.
*/
func RunParallelFiles(config Config) {
	// Split the data directories by "+" and process each one
	dataDirs := strings.Split(config.DataDirs, "+")

	// Open effects file
	effectsPathFile := "../data/effects.txt"
	effectsFile, err := os.Open(effectsPathFile)
	if err != nil {
		panic("Failed to open effects file")
	}
	defer effectsFile.Close()

	// Create task queue and TAS lock
	queue := &TaskQueue{}
	lock := &TASLock{}

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

	// Spawn Go routines
	numGoroutines := min(config.ThreadCount, len(queue.tasks)) // min(command line threads, mun of images the queue)
	var wg sync.WaitGroup                                      // wait group from Go
	wg.Add(numGoroutines)

	// Start timer for the parallel section
	startParallel := time.Now()

	// Worker function for each goroutine
	worker := func() {
		defer wg.Done()
		for {
			// Acquire TAS lock to get the next task
			lock.Lock()
			task := queue.Dequeue()
			lock.Unlock()

			// Start timer for the parallel section of this image processing
			//startParallel := time.Now()

			if task == nil {
				return // No more tasks in the queue
			}

			// Have goroutine process the image
			processImage(task)

			// End timer and calculate duration for this image
			//parallelDuration := time.Since(startParallel).Seconds()
			//fmt.Print(parallelDuration, "\n")
		}
	}

	// Start goroutines
	// run until all tasks from the queue are processed
	for i := 0; i < numGoroutines; i++ {
		go worker()
	}

	// Wait for all goroutines to terminate
	wg.Wait()

	// End timer for the parallel section and calculate duration
	//parallelDuration := time.Since(startParallel).Seconds()
	//fmt.Printf("Parallel Section Execution Time: %.2f seconds\n", parallelDuration)

}

// Helper to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Process an image based on the task
func processImage(task *Task) {
	// Open the input image
	imgFile, err := os.Open(task.inPath)
	if err != nil {
		fmt.Printf("Failed to open image file %s: %v\n", task.inPath, err)
		return
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Printf("Failed to decode image %s: %v\n", task.inPath, err)
		return
	}
	outImg := img

	// Apply each effect in sequence
	for _, ef := range task.effects {
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

	// Save the processed image
	outFile, err := os.Create(task.outPath)
	if err != nil {
		fmt.Printf("Failed to create output file %s: %v\n", task.outPath, err)
		return
	}
	defer outFile.Close()

	err = png.Encode(outFile, outImg)
	if err != nil {
		fmt.Printf("Failed to encode image %s: %v\n", task.outPath, err)
	}
}
