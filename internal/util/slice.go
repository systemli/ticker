package util

// Contains returns true when s is in slice.
func Contains(slice []int, s int) bool {
	for _, value := range slice {
		if value == s {
			return true
		}
	}
	return false
}

// Append appends the integer when it not already in the slice.
func Append(slice []int, s int) []int {
	if Contains(slice, s) {
		return slice
	}

	return append(slice, s)
}

// Remove element from slice.
func Remove(slice []int, s int) []int {
	for i := range slice {
		if slice[i] == s {
			slice = append(slice[:i], slice[i+1:]...)
		}
	}

	return slice
}

// ContainsString returns true when s in slice.
func ContainsString(list []string, s string) bool {
	for _, b := range list {
		if b == s {
			return true
		}
	}
	return false
}
