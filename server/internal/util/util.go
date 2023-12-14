package util

// StringPtr returns a pointer to the passed string.
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns a pointer to the passed int.
func IntPtr(s int64) *int64 {
	return &s
}
