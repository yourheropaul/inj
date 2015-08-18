package inj

import (
	"reflect"
	"testing"
)

// All empty types should return zero
func Test_Zero(t *testing.T) {

	var zeroConcrete *ConcreteType
	var zeroNested *NestedType

	inputs := []interface{}{
		ConcreteType{},
		zeroConcrete,
		NestedType{},
		zeroNested,
	}

	for i, input := range inputs {
		if !zero(reflect.ValueOf(input)) {
			t.Errorf("[%d] Value not zero", i)
		}
	}
}
