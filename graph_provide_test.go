package inj

import "testing"

// Objects passed to a graph should appear in the graph and be populated
func Test_ProvideHappyPath(t *testing.T) {

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
}

// New dependency provisions shouldn't overwrite previously set ones
func Test_ProvideOverride(t *testing.T) {

	g, c := NewGraph(), ConcreteType{}

	if err := g.Provide(
		&c,
		DEFAULT_STRING,
	); err != nil {
		t.Fatalf("Graph.Provide: %s", err)
	}

	// The graph now includes DEFAULT_STRING as its
	// only met depdendency (missing dependencies covered
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
}
