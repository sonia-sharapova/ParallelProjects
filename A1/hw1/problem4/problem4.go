// References:
// Command line arguments: https://gobyexample.com/command-line-arguments
// Reading-n file lines: https://stackoverflow.com/questions/8757389/reading-a-file-line-by-line-in-go
// Go with OS: https://pkg.go.dev/os
package main

import (
	"bufio"
	"fmt"
	"hw1/set"
	"os"
	"sort"
	"strconv"
	"strings"
)

// for set notation syntax
func formatSet(intSet *set.IntSet) string {
	// sorted for test cases
	sorted := intSet.IntSeq
	sort.Ints(sorted)

	// formatting
	strs := make([]string, len(sorted))
	for i, num := range sorted {
		strs[i] = strconv.Itoa(num)
	}
	return "{" + strings.Join(strs, ",") + "}" // have right brackets
}

// read a set from a line of integers
func readSetFromLine(line string) *set.IntSet {
	intSet := set.NewIntSet()
	vals := strings.Fields(line)
	for _, nums := range vals {
		num, err := strconv.Atoi(nums)
		if err != nil {
			fmt.Println("Couldn't read file lines:", err)
			os.Exit(0)
		}
		intSet.Add(num)
	}
	return intSet
}

func main() {
	// check arguments
	if len(os.Args) != 3 {
		fmt.Println("Usage: problem4  <union | intersect | diff> file")
		return
	}

	op := os.Args[1]
	file_name := os.Args[2]

	if op != "union" && op != "intersect" && op != "diff" {
		fmt.Println("Usage: problem4  <union | intersect | diff> file")
		return
	}

	// Check if file name is empty or invalid
	if file_name == "" {
		fmt.Println("Usage: problem4  <union | intersect | diff> file")
		return
	}

	// open file, give usage if error
	file, err := os.Open(file_name)
	if err != nil {
		fmt.Println("Usage: problem4  <union | intersect | diff> file")
		return
	}
	defer file.Close()

	// read file lines using scanner
	scanner := bufio.NewScanner(file)
	var sets [2]*set.IntSet
	for i := 0; i < 2; i++ {
		if scanner.Scan() {
			sets[i] = readSetFromLine(scanner.Text())
		}
	}

	// for the given operations
	var result *set.IntSet
	switch op {
	case "union":
		result = sets[0].Union(sets[1])
	case "intersect":
		result = sets[0].Intersect(sets[1])
	case "diff":
		result = sets[0].Diff(sets[1])
	default:
		fmt.Println("Usage: problem4  <union | intersect | diff> file")
		return
	}

	fmt.Println(formatSet(result))
}
