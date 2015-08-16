package inj

import (
	"fmt"
	"testing"
)

type StandardPasserType struct{}

func (s StandardPasserType) Fn(i1 InterfaceOne, i2 InterfaceTwo) {
	fmt.Println(i1.SayHello())
}

func StandardPasserFn(i1 InterfaceOne, i2 InterfaceTwo) {
	fmt.Println(i1.SayHello())
}

func Test_GraphSimpleInjection(t *testing.T) {

	g := NewGraph()

	// This is tested elsewhere
	if err := g.Provide(&helloSayer{}, &goodbyeSayer{}); err != nil {
		t.Fatalf("Graph.Provide: %s", err)
	}

	// Pass an anonymous function
	g.Inject(func(i1 InterfaceOne, i2 InterfaceTwo) {
		fmt.Println(i1.SayHello())
	})

	// Pass a first class function
	g.Inject(StandardPasserFn)

	// Pass a member of a struct
	spt := StandardPasserType{}
	g.Inject(spt.Fn)
}

func Test_GraphComplexInjection(t *testing.T) {

	// Don't provide anything for this graph
	g := NewGraph()

	// Pass an anonymous function
	g.Inject(func(i1 InterfaceOne, i2 InterfaceTwo) {
		fmt.Println(i1.SayHello())
	}, &helloSayer{}, &goodbyeSayer{})

	// Pass a first class function
	g.Inject(StandardPasserFn, &helloSayer{}, &goodbyeSayer{})

	// Pass a member of a struct
	spt := StandardPasserType{}
	g.Inject(spt.Fn, &helloSayer{}, &goodbyeSayer{})
}
