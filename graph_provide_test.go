package inj

import "testing"

//////////////////////////////////////////
// Unit tests
//////////////////////////////////////////

// Objects passed to a graph should appear in the graph and be populated
func Test_ProvideHappyPath1(t *testing.T) {

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

	// There should be exactly six nodes in the graph now
	if g, e := len(g.Nodes), 6; g != e {
		t.Errorf("Got %d nodes, expected %d", g, e)
	}

	// Check the whole type
	assertConcreteValue(c, t)

	// Check index count
	if g, e := len(g.indexes), 6; g != e {
		t.Errorf("Expected %d indexes, got %d", g, e)
	}
}

// Multiple calls to Provide should have the same effect as a single call
func Test_ProvideHappyPath2(t *testing.T) {

	g, c := NewGraph(), ConcreteType{}

	vs := []interface{}{
		&c,
		&helloSayer{},
		&goodbyeSayer{},
		funcInstance,
		ichannel,
		DEFAULT_STRING,
	}

	// Register providers individually
	for _, v := range vs {
		if err := g.Provide(v); err != nil {
			t.Fatalf("Graph.Provide: %s", err)
		}
	}

	// There should be exactly six nodes in the graph now
	if g, e := len(g.Nodes), 6; g != e {
		t.Errorf("Got %d nodes, expected %d", g, e)
	}

	// Check the whole type
	assertConcreteValue(c, t)

	// Check index count
	if g, e := len(g.indexes), 6; g != e {
		t.Errorf("Expected %d indexes, got %d", g, e)
	}
}

// New dependency provisions shouldn't overwrite previously set ones
func Test_ProvideOverride1(t *testing.T) {

	g, c := NewGraph(), ConcreteType{}

	if err := g.Provide(
		&c,
		DEFAULT_STRING,
	); err != nil {
		t.Fatalf("Graph.Provide: %s", err)
	}

	// Check index count
	if g, e := len(g.indexes), 2; g != e {
		t.Errorf("[1] Expected %d indexes, got %d", g, e)
	}

	// The graph now includes DEFAULT_STRING as its
	// only met dependency (missing dependencies covered
	// by graph_assert_test.go). Adding another string
	// to the graph shouldn't alter the value of the
	// concrete type.

	if err := g.Provide(
		&c,
		"some other string",
	); err != nil {
		t.Fatalf("Graph.Provide: %s", err)
	}

	if g, e := c.String, DEFAULT_STRING; g != e {
		t.Errorf("Got %s, expected %s", g, e)
	}

	// Check index count
	if g, e := len(g.indexes), 2; g != e {
		t.Errorf("[1] Expected %d indexes, got %d", g, e)
	}
}

//////////////////////////////////////////
// Benchmark tests
//////////////////////////////////////////

func BenchmarkComplexProvide(b *testing.B) {

	for n := 0; n < b.N; n++ {
		NewGraph(
			&ConcreteType{},
			&helloSayer{},
			&goodbyeSayer{},
			funcInstance,
			ichannel,
			DEFAULT_STRING,
		)
	}
}

func BenchmarkManyProvisions(b *testing.B) {

	for n := 0; n < b.N; n++ {
		g := NewGraph()
		g.Provide(&ConcreteType{})
		g.Provide(&helloSayer{})
		g.Provide(&goodbyeSayer{})
		g.Provide(funcInstance)
		g.Provide(ichannel)
		g.Provide(DEFAULT_STRING)
	}
}
