package inj

import (
	"fmt"
	"testing"
)

// Test a linked graph
func Test_Injection(t *testing.T) {

	g := NewGraph()

	c := ConcreteType{}

	// Register providers (can include non-providers, which will then be wired up)
	if err := g.Provide(&c, &helloSayer{}, &goodbyeSayer{}, funcInstance); err != nil {
		t.Fatalf("Graph.Provide: %s", err)
	}

	// Add a named provider
	//Provider(helloSayer, "some name")

	// Assert

	// Inject

	fmt.Println(c.Stringer("echo"))
	fmt.Println(c.Nested.Goodbye.SayGoodbye())
}
