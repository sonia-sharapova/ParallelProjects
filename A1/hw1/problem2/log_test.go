package problem2

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

func TestLog(t *testing.T) {

	var tests = []struct {
		file        string
		maxDuration int
		expected    int
	}{
		{"test0.txt", 5, 2},
		{"test1.txt", 100, 1},
		{"test2.txt", 1, 1},
		{"test3.txt", 1, 18},
		{"test4.txt", 10, 18},
		{"test5.txt", 100000, 17},
		{"test6.txt", 1000000000, 33},
		{"test7.txt", 1, 18491},
		{"test8.txt", 1, 18487},
		{"test9.txt", 10, 18516},
		{"test10.txt", 10, 18463},
		{"test11.txt", 100000, 18496},
		{"test12.txt", 100000, 18482},
		{"test13.txt", 1000000000, 33333},
		{"test14.txt", 1000000000, 33333},
	}

	for num, test := range tests {
		testname := fmt.Sprintf("T=%v", num)
		t.Run(testname, func(t *testing.T) {
			input := readFile(test.file)
			got := ProcessLog(input, test.maxDuration)
			if got != test.expected {
				t.Errorf("Called ProcessLog for test file=%v, with max duration = %v)\nExpected:%v\nGot:%v", test.file, test.maxDuration, test.expected, got)
			}
		})
	}
}
