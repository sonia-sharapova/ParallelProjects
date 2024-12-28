// References:
// Making a field accessible from outside of the package: https://www.digitalocean.com/community/tutorials/understanding-package-visibility-in-go
// Struct tutorial: https://gobyexample.com/structs

package set

import "hw1/utils"

// IntSet represents a set of integers
type IntSet struct {
	IntSeq []int // Uses an []int slice to hold the collection of integers that are part of the set.
}

// NewIntSet returns a newlly allocated IntSet with its contents initilaized to zero
func NewIntSet() *IntSet {
	return &IntSet{[]int{}}
}

// adds a number to the set (no duplicates)
func (recv *IntSet) Add(num int) {
	for _, n := range recv.IntSeq {
		if n == num {
			return
		}
	}
	recv.IntSeq = append(recv.IntSeq, num)
}

// set union operation (values in other or recv)
// using 'Add' function
func (recv *IntSet) Union(other *IntSet) *IntSet {
	result := NewIntSet()
	// for receiver
	for _, n := range recv.IntSeq {
		result.Add(n)
	}
	// for other
	for _, n := range other.IntSeq {
		result.Add(n)
	}
	return result
}

// intersect operation (values in both other and recv)
func (recv *IntSet) Intersect(other *IntSet) *IntSet {
	result := NewIntSet()
	for _, n := range recv.IntSeq {
		if utils.Contains(other.IntSeq, n) {
			result.Add(n)
		}
	}
	return result
}

// diff function: set difference between recv and other
func (recv *IntSet) Diff(other *IntSet) *IntSet {
	result := NewIntSet()
	for _, n := range recv.IntSeq {
		if !utils.Contains(other.IntSeq, n) {
			result.Add(n)
		}
	}
	return result
}
