package inj

import (
	"reflect"
	"testing"
)

// If the main type is a pointer, the specific type should be the
// indirect value
func Test_GetReflectionTypes(t *testing.T) {

	m, i := &ConcreteType{}, ConcreteType{}
	tm, ti := reflect.TypeOf(m), reflect.TypeOf(i)
	rm, ri := getReflectionTypes(m)

	if rm != tm {
		t.Errorf("Master type doesn't match expected type")
	}

	if ri != ti {
		t.Errorf("Specific type doesn't match expected type")
	}
}
