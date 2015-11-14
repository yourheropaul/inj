package inj

import "strings"

type structPath string

// Create a new, empty struct path
func emptyStructPath() structPath {
	return ""
}

// Create a seprate struct path with a parent name
func (s structPath) Branch(parent string) structPath {
	return s + structPath("."+parent)
}

// Remove the first element from the path, and return it
// and the resultant path
func (s structPath) Shift() (string, structPath) {

	if s.Empty() {
		return "", emptyStructPath()
	}

	parts := strings.Split(string(s), ".")

	// Make a new structpath
	var s2 structPath

	if len(parts) > 2 {
		s2 = structPath("." + strings.Join(parts[2:], "."))
	} else {
		s2 = structPath(strings.Join(parts[2:], "."))
	}

	return parts[1], s2
}

// Returns true if the path is empty
func (s structPath) Empty() bool {

	if string(s) == "" {
		return true
	}

	return false
}

// Implementation of the Stringer interface
func (s structPath) String() string {
	return string(s)
}
