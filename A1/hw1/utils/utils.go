package utils

// Contains checks if an integer is inside of a integer slice.
// It returns True if the integer is; otherwise, false.
func Contains(slice []int, findMe int) bool {
	for _, element := range slice {
		if element == findMe {
			return true
		}
	}
	return false
}
