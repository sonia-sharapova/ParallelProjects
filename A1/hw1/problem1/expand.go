// References:
// trimming whitespaces: https://www.codecademy.com/resources/docs/go/strings/trim
// Contains function: https://www.geeksforgeeks.org/string-contains-function-in-golang-with-examples/
// Splitting range with split function: https://pkg.go.dev/strings#Split
// Go structs: https://go.dev/tour/moretypes/2

package problem1

import (
	"strconv"
	"strings"
)

func Expand(intList string) []int {
	// Input: intList - a comma-separated list of integers and/or integer ranges,
	// Returns: Slice of all the integers with the ranges expanded into individual integers.

	// make map (struct) to ensure unique nums
	unique := make(map[int]struct{})

	// Split the input by commas to process individual components
	comp := strings.Split(intList, ",")

	// Iterate over each component
	for _, comp := range comp {
		// Deal with whitespaces (leading and tailing)
		comp := strings.Trim(comp, " ")

		// Range: Check if the compomponent contains '-'
		if strings.Contains(comp, "-") {

			// Split the range into MIN and MAX
			ranges := strings.Split(comp, "-") // []string

			if len(ranges) == 2 {
				// Convert MIN and MAX to integers
				min, minErr := strconv.Atoi(ranges[0])
				max, maxErr := strconv.Atoi(ranges[1])

				// Check for no errors, no negatives, and that MIN < MAX
				if minErr == nil && maxErr == nil && min > 0 && max > 0 && min < max {
					// Add the expanded range to the result
					for i := min; i <= max; i++ {
						unique[i] = struct{}{}
					}
				}
			}
		} else {
			// convert single component to integer
			num, err := strconv.Atoi(comp)
			// check for no error and sign
			if err == nil && num > 0 {
				unique[num] = struct{}{}
			}
		}
	}

	// Convert map to slice of integers
	result := []int{}
	for num := range unique {
		result = append(result, num)
	}

	return result

}
