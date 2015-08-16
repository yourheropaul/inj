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

func Test_AssertionComplexHappyPath(t *testing.T) {

	g, c := NewGraph(), ConcreteType{}

	// Register providers (can include non-providers, which will then be wired up)
	if err := g.Provide(
		&c,
		&helloSayer{},
		&goodbyeSayer{},
		funcInstance,
		ichannel,
		DEFAULT_STRING,
	); err != nil {
		t.Fatalf("Graph.Provide: %s", err)
	}

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
