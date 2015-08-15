package inj

import "reflect"

// A unique indentifier for a given reflect.Type
func identifier(t reflect.Type) string {
	return t.PkgPath() + "/" + t.String()
}
