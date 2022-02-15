package helpers

// Checks whether or not the provided string
// array contains the provided needle.
func Contains(haystack []string, needle string) bool {
	for _, a := range haystack {
		if a == needle {
			return true
		}
	}

	return false
}
