package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

var epsilon float64 = 0.00000001

func floatEquals(aStr, bStr string) bool {
	a, _ := strconv.ParseFloat(aStr, 64)
	b, _ := strconv.ParseFloat(bStr, 64)
	return (a-b) < epsilon && (b-a) < epsilon
}

func TestPi(t *testing.T) {

	var tests = []struct {
		threads  string
		interval string
		expected string
	}{
		{"1", "100", "3.1315929036"},
		{"2", "100", "3.1315929036"},
		{"3", "100", "3.1315929036"},
		{"4", "100", "3.1315929036"},
		{"1", "1000000", "3.1415916536"},
		{"2", "1000000", "3.1415916536"},
		{"3", "1000000", "3.1415916536"},
		{"4", "1000000", "3.1415916536"},
		{"1", "10000000", "3.1415925536"},
		{"2", "10000000", "3.1415925536"},
		{"3", "10000000", "3.1415925536"},
		{"4", "10000000", "3.1415925536"},
		{"1", "100000000", "3.1415926436"},
		{"2", "100000000", "3.1415926436"},
		{"3", "100000000", "3.1415926436"},
		{"4", "100000000", "3.1415926436"},
	}
	for num, test := range tests {
		testname := fmt.Sprintf("T=%v", num)
		t.Run(testname, func(t *testing.T) {
			var err error
			cmd := exec.Command("go", "run", "hw2/pi", test.interval, test.threads)
			out, err := cmd.Output()
			sOut := strings.TrimSpace(string(out))

			if err != nil || !floatEquals(test.expected, sOut) {
				t.Errorf("\nRan:%s\nExpected:%s\nGot:%s", cmd, test.expected, sOut)
			}
		})
	}

}
