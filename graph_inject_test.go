package inj

import "testing"

//////////////////////////////////////////
// Standard injection testers
//////////////////////////////////////////

type StandardPasserType struct{}

func (s StandardPasserType) Fn(i1 InterfaceOne, i2 InterfaceTwo, t *testing.T) {
	assertPasserInterfaceValues(i1, i2, t)
}

func StandardPasserFn(i1 InterfaceOne, i2 InterfaceTwo, t *testing.T) {
	assertPasserInterfaceValues(i1, i2, t)
}

//////////////////////////////////////////
// Assertions for passer types
//////////////////////////////////////////

func assertPasserInterfaceValues(i1 InterfaceOne, i2 InterfaceTwo, t *testing.T) {

	if i1 == nil {
		t.Errorf("i1 is nil")
	}

	if g, e := i1.SayHello(), HELLO_SAYER_MESSAGE; g != e {
		t.Errorf("i2.SayHello(): got %s, expected %s", g, e)
	}

	if i2 == nil {
		t.Errorf("i1 is nil")
	}

	if g, e := i2.SayGoodbye(), GOODBYE_SAYER_MESSAGE; g != e {
		t.Errorf("i2.SayHello(): got %s, expected %s", g, e)
	}
}

//////////////////////////////////////////
// Unit tests
//////////////////////////////////////////

// A basic test of the entire injection feature
func Test_GraphSimpleInjectionHappyPath(t *testing.T) {

	g := NewGraph()

	// This is tested elsewhere
	if err := g.Provide(
		&helloSayer{},
		&goodbyeSayer{},
		t,
	); err != nil {
		t.Fatalf("Graph.Provide: %s", err)
	}

	// Pass an anonymous function
	g.Inject(func(i1 InterfaceOne, i2 InterfaceTwo) {
		assertPasserInterfaceValues(i1, i2, t)
	})

	// Pass a first class function
	g.Inject(StandardPasserFn)

	// Pass a member of a struct
	spt := StandardPasserType{}
	g.Inject(spt.Fn)
}

// Graph.Inject should panic if a non-function is passed
func Test_GraphSimpleInjectionSadPath1(t *testing.T) {

	defer func() {
		if recover() != nil {
			// The test has succeeded
		}
	}()

	g := NewGraph()
	g.Inject("not a func")

	t.Error("Inject failed to panic with non-func type")
}

// Should panic on missing dependencies
func Test_GraphSimpleInjectionSadPath2(t *testing.T) {

	defer func() {
		if recover() != nil {
			// The test has succeeded
		}
	}()

	g := NewGraph()
	g.Inject(func(s string) {})

	t.Error("Inject failed to panic with no depdency provided")
}

// Should panic if a variadic function is passed
func Test_GraphSimpleInjectionSadPath3(t *testing.T) {

	defer func() {
		if recover() != nil {
			// The test has succeeded
		}
	}()

	g := NewGraph()
	g.Inject(func(s ...string) {})

	t.Error("Inject failed to panic with a variadic function argument")
}

// Complex injection is essentially passing additional variables
func Test_GraphComplexInjectionHappyPath1(t *testing.T) {

	// Don't provide anything for this graph
	g := NewGraph()

	// Pass an anonymous function
	g.Inject(func(i1 InterfaceOne, i2 InterfaceTwo) {
		assertPasserInterfaceValues(i1, i2, t)
	}, &helloSayer{}, &goodbyeSayer{}, t)

	// Pass a first class function
	g.Inject(StandardPasserFn, &helloSayer{}, &goodbyeSayer{}, t)

	// Pass a member of a struct
	spt := StandardPasserType{}
	g.Inject(spt.Fn, &helloSayer{}, &goodbyeSayer{}, t)
}

// Dependencies should come from the graph before the xargs
func Test_GraphComplexInjectionHappyPath2(t *testing.T) {

	// Don't provide anything for this graph
	g := NewGraph()
	g.Provide("string one")

	// Pass an anonymous function
	g.Inject(func(s string) {
		if s != "string one" {
			t.Fatalf("Expected 'string one', got %s", s)
		}
	}, "string two")
}
