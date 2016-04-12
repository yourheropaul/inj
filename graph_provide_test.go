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
	if g, e := len(g.nodes), 6; g != e {
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
	if g, e := len(g.nodes), 6; g != e {
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

	if g, e := c.String, DEFAULT_STRING; g == e {
		t.Errorf("Got %s, expected %s", g, e)
	}

	// Check index count
	if g, e := len(g.indexes), 2; g != e {
		t.Errorf("[1] Expected %d indexes, got %d", g, e)
	}
}

// Embedded structs should be parsed like any other
func Test_EmbeddedStructProvision(t *testing.T) {

	g := NewGraph()

	e, he := Embeddable{X: 10}, &HasEmbeddable{}

	if err := g.Provide(e, he); err != nil {
		t.Fatalf("g.Provide: %s", err.Error())
	}

	if v, errs := g.Assert(); !v {
		t.Error(errs)
	}

	if g, e := e.X, he.X; g != e {
		t.Errorf("Values don't match: got %d, expected %d", g, e)
	}
}

// Self-referential dependencies shouldn't be assigned
func Test_SelfReferencingDoesntHappen(t *testing.T) {

	g := NewGraph()

	o := Ouroboros1{}

	g.Provide(&o)

	valid, errs := g.Assert()

	if valid {
		t.Fatalf("g.Assert() is valid when it shouldn't be")
	}

	// There are two deps that should have missed
	if g, e := len(errs), 2; g != e {
		t.Fatalf("Expected %d errors, got %d (%v)", e, g, errs)
	}
}

// Self-referential dependencies shouldn't impede proper injection
func Test_SelfReferencingShouldntCircumentInjection(t *testing.T) {

	g := NewGraph()

	o1 := Ouroboros1{V: 1}
	o2 := Ouroboros2{V: 2}

	g.Provide(&o1, &o2)

	valid, errs := g.Assert()

	if !valid {
		t.Fatalf("g.Assert() is not valid when it should be (%v)", errs)
	}

	// The values should now be 'crossed'
	if o1.A.Value() != o1.B.Value() {
		t.Errorf("o1.A != o1.B")
	}

	if o1.A.Value() != o2.Value() {
		t.Errorf("o1.B and B aren't equal to o2")
	}

	if o2.A.Value() != o2.B.Value() {
		t.Errorf("o2.A != o2.B")
	}

	if o2.A.Value() != o1.Value() {
		t.Errorf("o2.B and B aren't equal to o1")
	}
}

// Self-referential prevention must extend to embedding
func Test_EmbeddedSelfReferencingDoesntHappen(t *testing.T) {

	g := NewGraph()

	o := Ouroboros4{}

	g.Provide(&o)

	valid, errs := g.Assert()

	if valid {
		t.Fatalf("g.Assert() is valid when it shouldn't be")
	}

	// There is one dep that should have missed
	if g, e := len(errs), 1; g != e {
		t.Fatalf("Expected %d error, got %d (%v)", e, g, errs)
	}
}

// Self-referential dependencies shouldn't impede proper injection
func Test_EmbeddedSelfReferencingShouldntCircumentInjection(t *testing.T) {

	g := NewGraph()

	o1 := Ouroboros3{V: 1}
	o2 := Ouroboros4{V: 2}

	g.Provide(o1, &o2)

	valid, errs := g.Assert()

	if !valid {
		t.Fatalf("g.Assert() is not valid when it should be (%v)", errs)
	}

	// The value should now be assigned
	if o2.Ouroboros3.Value() != o1.Value() {
		t.Errorf("o2.A != o2.B")
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
