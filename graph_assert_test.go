package inj

import (
	"fmt"
	"testing"
)

type AssertionTester struct {
	S string `inj:""`
	I int    `inj:""`
}

func Test_AssertionHappyPath(t *testing.T) {

	g := NewGraph()

	g.Provide(&AssertionTester{}, "hello", 1)

	if v, m := g.Assert(); !v {
		fmt.Println(m)
		t.Error("Assert() failed")
	}
}

func Test_AssertionSadPath(t *testing.T) {

	g := NewGraph()

	// Only provide the concrete type
	g.Provide(&AssertionTester{})

	v, m := g.Assert()

	if v {
		t.Error("Assert() didn't fail")
	}

	if g, e := len(m), 2; g != e {
		t.Errorf("Expected %d errors, got %d", e, g)
	}
}
