package inj

import "strings"

type StructPath string

// Create a new, empty struct path
func EmptyStructPath() StructPath {
	return ""
}

// Create a seprate struct path with a parent name
func (s StructPath) Branch(parent string) StructPath {
	return s + StructPath("."+parent)
}

// Remove the first element from the path, and return it
// and the resultant path
func (s StructPath) Shift() (string, StructPath) {

	if s.Empty() {
		return "", EmptyStructPath()
	}

	parts := strings.Split(string(s), ".")

	// Make a new structpath
	var s2 StructPath

	if len(parts) > 2 {
		s2 = StructPath("." + strings.Join(parts[2:], "."))
	} else {
		s2 = StructPath(strings.Join(parts[2:], "."))
	}

	return parts[1], s2
}

// Returns true if the path is empty
func (s StructPath) Empty() bool {

	if string(s) == "" {
		return true
	}

	return false
}

// Implementation of the Stringer interface
func (s StructPath) String() string {
	return string(s)
}
