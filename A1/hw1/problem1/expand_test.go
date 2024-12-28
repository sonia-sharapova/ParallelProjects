package problem1

import (
	"fmt"
	"testing"
)

func _containsP3test(slice []int, findMe int) bool {
	for _, element := range slice {
		if element == findMe {
			return true
		}
	}
	return false
}

func check(t *testing.T, in1 string, got, expected []int) {

	if len(got) != len(expected) {
		t.Errorf("\nCalled:Expand(%v)\nExpected:%v\nGot:%v", in1, expected, got)
		return
	}
	for _, item := range got {
		if !_containsP3test(expected, item) {
			t.Errorf("\nCalled:Expand(%v)\nExpected:%v\nGot:%v", in1, expected, got)
			return
		}
	}
}

func TestExpand(t *testing.T) {

	var tests = []struct {
		input    string
		expected []int
	}{
		{"", []int{}},
		{"2", []int{2}},
		{"1,2,3,3,3,3-4", []int{1, 2, 3, 4}},
		{"1-2,1-2,1-2,3,3,3,3-4,1-2,1-2,3-4", []int{1, 2, 3, 4}},
		{"Invalid,3,4", []int{3, 4}},
		{"4 43, 4,", []int{4}},
		{"1-4,7,3-5,10,12-14", []int{1, 2, 3, 4, 5, 7, 10, 12, 13, 14}},
		{"1-4, 4  , 5-6  , 5", []int{1, 2, 3, 4, 5, 6}},
		{"            1-4", []int{1, 2, 3, 4}},
		{"1-", []int{}},
		{"1-,45,-4,1-2", []int{45, 1, 2}},
		{"-4,----,4-,25  7, 55     ", []int{55}},
		{"17,1--4,", []int{17}},
		{"19-3,1-,Bob,1-2-3,4,", []int{4}},
		{"3-3,1-,Bob,1-2-3,4-4-4,          ", []int{}},
		{"3    -      4,1-,Bob,1-2-3,4-4-4,          ", []int{}},
	}

	for num, test := range tests {
		testname := fmt.Sprintf("T=%v", num)
		t.Run(testname, func(t *testing.T) {
			got := Expand(test.input)
			check(t, test.input, got, test.expected)
		})
	}
}
