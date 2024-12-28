package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func CreateFile(filePath string, nums1, nums2 string) {

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create: %s file.", filePath)
		os.Exit(1)
	}

	for _, num := range strings.Split(nums1, ",") {
		fmt.Fprintf(file, "%s ", num)
	}
	fmt.Fprintf(file, "\n")

	for _, num := range strings.Split(nums2, ",") {
		fmt.Fprintf(file, "%s ", num)
	}

	err = file.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not close: %s file.", filePath)
		os.Exit(1)
	}

}
func TestGoodArgs(t *testing.T) {

	var tests = []struct {
		nums1     string
		nums2     string
		operation string
		expected  string
	}{
		{"2,3", "1,2", "union", "{1,2,3}"},
		{"2,3", "1,2", "intersect", "{2}"},
		{"2,3", "1,2", "diff", "{3}"},
		{"1,2", "1,2", "union", "{1,2}"},
		{"1,2", "1,2", "intersect", "{1,2}"},
		{"1,2", "1,2", "diff", "{}"},
		{"1,2,3", "2,3,4", "union", "{1,2,3,4}"},
		{"5,2,3,9", "2,3,4", "diff", "{5,9}"},
		{"5,2,3,9", "2,3,4", "intersect", "{2,3}"},
	}
	filePath := "test.txt"
	for num, test := range tests {
		testname := fmt.Sprintf("T=%v", num)
		t.Run(testname, func(t *testing.T) {
			var err error
			CreateFile(filePath, test.nums1, test.nums2)
			cmd := exec.Command("go", "run", "hw1/problem4", test.operation, filePath)
			out, err := cmd.Output()
			sOut := strings.TrimSpace(string(out))

			if err != nil || test.expected != sOut {
				t.Errorf("\nRan:%s\ntest.txt --> A={%s}  B={%s}\nExpected:%s\nGot:%s", cmd, test.nums1, test.nums2,
					test.expected, sOut)
			}
		})
	}

}
func TestBadArgs(t *testing.T) {

	var tests = []struct {
		operation string
		file      string
	}{
		{"", ""},
		{"u", ""},
		{"-d", ""},
		{"union", ""},
		{"intersect", ""},
		{"diff", ""},
		{"-d", ""},
		{"intersect", "f"},
		{"diff", "t.txt"},
		{"-d", "file.txt"},
	}

	for num, test := range tests {
		testname := fmt.Sprintf("T=%v", num)
		t.Run(testname, func(t *testing.T) {
			var err error
			usage := "Usage: problem4  <union | intersect | diff> file"
			cmd := exec.Command("go", "run", "hw1/problem4", test.operation, test.file)
			out, err := cmd.Output()
			sOut := strings.TrimSpace(string(out))
			if err != nil || sOut != usage {
				t.Errorf("\nRan:%s\nExpected:%s\nGot:%s", cmd,
					usage, sOut)
			}
		})
	}
}
