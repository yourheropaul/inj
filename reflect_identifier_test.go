package inj

import (
	"reflect"
	"testing"
)

// Different instances of the same type should have the same identifier.
// This is a a bit like testing the reflect package, which is unncessary,
// but since it's an assumption of the inj package it's a useful sanity
// check.
func Test_Identifier(t *testing.T) {

	if reflect.TypeOf(ConcreteType{}) != reflect.TypeOf(ConcreteType{}) {
		t.Errorf("Different instance of the same type don't match")
	}

	if reflect.TypeOf(&ConcreteType{}) == reflect.TypeOf(ConcreteType{}) {
		t.Errorf("Different instance of the different types match")
	}
}
