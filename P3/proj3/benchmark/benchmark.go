package benchmark

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"proj3/dicom"
	"proj3/proc"
	"proj3/proc/workstealing"
)

type BenchmarkResult struct {
	Mode        proc.ExecutionMode
	Workers     int
	Duration    time.Duration
	FolderCount int
	FrameCount  int
}

// RunBenchmark runs a single benchmark with specified parameters
func RunBenchmark(inputDir, outputDir string, mode proc.ExecutionMode, workers int) (*BenchmarkResult, error) {
	start := time.Now()

	var processor proc.Processor
	switch mode {
	case proc.Sequential:
		processor = proc.NewSequentialProcessor()
	case proc.Pipeline:
		processor = proc.NewPipelineProcessor(workers, 10)
	case proc.WorkStealing:
		processor = workstealing.NewWorkStealingProcessor(workers)
	default:
		return nil, fmt.Errorf("unsupported mode: %s", mode)
	}

	// Count folders and frames before processing
	folderCount, frameCount, err := countFiles(inputDir)
	if err != nil {
		return nil, fmt.Errorf("error counting files: %w", err)
	}

	err = processor.ProcessDataset(inputDir, outputDir)
	if err != nil {
		return nil, fmt.Errorf("processing error: %w", err)
	}

	duration := time.Since(start)

	return &BenchmarkResult{
		Mode:        mode,
		Workers:     workers,
		Duration:    duration,
		FolderCount: folderCount,
		FrameCount:  frameCount,
	}, nil
}

// RunExperiments runs all benchmarks with different worker counts
func RunExperiments(inputDir string, workerCounts []int) error {
	// Create results directory
	resultsDir := "benchmark_results"
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return fmt.Errorf("failed to create results directory: %w", err)
	}

	// Open CSV file for results
	csvFile, err := os.Create(fmt.Sprintf("%s/results.csv", resultsDir))
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Write CSV header
	header := []string{"Mode", "Workers", "Duration_ms", "Folders", "Frames", "Speedup"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Run sequential benchmark first
	seqResult, err := RunBenchmark(inputDir, fmt.Sprintf("%s/sequential", resultsDir), proc.Sequential, 1)
	if err != nil {
		return fmt.Errorf("sequential benchmark failed: %w", err)
	}

	// Write sequential result
	if err := writeResult(writer, seqResult, seqResult.Duration); err != nil {
		return fmt.Errorf("failed to write sequential result: %w", err)
	}

	// Run parallel benchmarks
	modes := []proc.ExecutionMode{proc.Pipeline, proc.WorkStealing}
	for _, mode := range modes {
		for _, workers := range workerCounts {
			result, err := RunBenchmark(
				inputDir,
				fmt.Sprintf("%s/%s_%d", resultsDir, mode, workers),
				mode,
				workers,
			)
			if err != nil {
				fmt.Printf("Warning: benchmark failed for mode=%s, workers=%d: %v\n", mode, workers, err)
				continue
			}

			if err := writeResult(writer, result, seqResult.Duration); err != nil {
				return fmt.Errorf("failed to write result: %w", err)
			}
		}
	}

	return nil
}

func writeResult(writer *csv.Writer, result *BenchmarkResult, seqDuration time.Duration) error {
	speedup := float64(seqDuration) / float64(result.Duration)
	record := []string{
		string(result.Mode),
		strconv.Itoa(result.Workers),
		strconv.FormatInt(result.Duration.Milliseconds(), 10),
		strconv.Itoa(result.FolderCount),
		strconv.Itoa(result.FrameCount),
		strconv.FormatFloat(speedup, 'f', 4, 64),
	}
	return writer.Write(record)
}

func countFiles(inputDir string) (folders, frames int, err error) {
	entries, err := filepath.Glob(filepath.Join(inputDir, "*"))
	if err != nil {
		return 0, 0, err
	}

	folderCount := 0
	totalFrames := 0

	for _, entry := range entries {
		if info, err := os.Stat(entry); err != nil || !info.IsDir() {
			continue
		}

		files, err := dicom.LoadDICOMFiles(entry)
		if err != nil {
			continue
		}

		if len(files) >= 2 { // Only count folders with enough frames
			folderCount++
			totalFrames += len(files)
		}
	}

	return folderCount, totalFrames, nil
}
