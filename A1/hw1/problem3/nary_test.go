package problem3

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func readFile(filePath string) []string {

	inFile, _ := os.Open(fmt.Sprintf("./tests/%v", filePath))
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	var input []string

	for scanner.Scan() {
		line := scanner.Text()
		input = append(input, line)
	}
	return input
}

func TestNAry(t *testing.T) {

	var tests = []struct {
		file     string
		nary     int
		expected int
	}{
		{"test0.txt", 2, 1},
		{"test1.txt", 3, 0},
		{"test2.txt", 0, 1},
		{"test3.txt", 2, 3},
		{"test4.txt", 3, 0},
		{"test5.txt", 1, 5},
		{"test6.txt", 1, 0},
		{"test7.txt", 0, 5},
		{"test8.txt", 1, 6},
		{"test9.txt", 1, 1},
		{"test10.txt", 0, 4},
		{"test11.txt", 1, 0},
		{"test12.txt", 0, 4},
		{"test13.txt", 1, 5},
		{"test14.txt", 0, 7},
	}

	for num, test := range tests {
		testname := fmt.Sprintf("T=%v", num)
		t.Run(testname, func(t *testing.T) {
			input := readFile(test.file)
			got := SearchTree(input, test.nary)
			if got != test.expected {
				t.Errorf("Called SearchTree for test file=%v, with max duration = %v)\nExpected:%v\nGot:%v", test.file, test.nary, test.expected, got)
			}
		})
	}
}
