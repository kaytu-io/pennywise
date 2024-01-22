package region

// Code represents an AWS region code.
type Code string

// Valid returns true if the region exists and is supported, false otherwise.
func (c Code) Valid() bool {
	if c == "" {
		return false
	}
	_, ok := codeToName[c]
	return ok
}

// String returns the code of the region as a string.
func (c Code) String() string {
	return string(c)
}
