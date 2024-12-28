/*
References:
- Command line arguments: https://gobyexample.com/command-line-arguments
- Implementing Leibniz formula:
	https://en.wikipedia.org/wiki/Leibniz_formula_for_%CF%80
	https://proofwiki.org/wiki/Leibniz's_Formula_for_Pi
- Types in Go: https://go.dev/tour/basics/11
- GoLang Concurrency and fork-join examples: https://rogerwelin.github.io/golang/go/concurrency/2018/09/04/golang-concurrency-fundamentals.html
- WaitGroups: https://gobyexample.com/waitgroups
*/

package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
)

// The main goroutine
func main() {

	var pi float64

	// Read interval and threads arguments
	interval, _ := strconv.Atoi(os.Args[1])
	threads, _ := strconv.Atoi(os.Args[2])

	// check if interval < 1000
	if interval <= 1000 {
		// Call sequential implementation
		pi = sequentialPi(interval)
	} else {
		// Call parallel implementation
		pi = parallelPi(interval, threads)
	}
	fmt.Printf("%.10f\n", pi)
	return
}

/* Estimating Pi using the Leibniz formula */

// Sequential:
func sequentialPi(interval int) float64 {
	var sum float64
	// Sum over calculations in work interval
	for k := 0; k < interval; k++ {
		sum += math.Pow(-1, float64(k)) / (2.0*float64(k) + 1.0)
	}
	pi := sum * 4 // multiply result by 4
	return pi
}

// Parallel:
func parallelPi(interval int, threads int) float64 {
	var sum float64
	var wg sync.WaitGroup

	// Divide work equally among threads
	threadWork := interval / threads
	remainder := interval % threads // remainder (given to last thread)

	// Launch goroutines:
	// 		- based on number of threads
	// 		- static task distribution technique
	for t := 0; t < threads; t++ {
		wg.Add(1)

		// Create goroutine: the concurrent calcultaions
		go func(threadID int) {
			defer wg.Done()
			localSum := 0.0
			// Get work interval based on thread ID
			start := threadID * threadWork
			end := start + threadWork

			// Any remaining work is done by last thread
			if threadID == threads-1 {
				end += remainder
			}

			// Sum over calculations in work interval
			for k := start; k < end; k++ {
				localSum += math.Pow(-1, float64(k)) / (2.0*float64(k) + 1.0)
			}

			// accumulate results from saparate threads
			sum = sum + localSum
		}(t)
	}

	// Wait for all goroutines to finish
	wg.Wait() // join point

	pi := sum * 4 // multiply result by 4
	return pi
}
