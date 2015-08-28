package util

// ContainsString returns true if a string slice contains the target string
// and false otherwise.
func ContainsString(haystack []string, needle string) bool {
	for _, straw := range haystack {
		if straw == needle {
			return true
		}
	}
	return false
}
