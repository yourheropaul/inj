package inj

import "testing"

///////////////////////////////////////////////////////////////
// Basic datasource/graph integration tests
///////////////////////////////////////////////////////////////

func Test_ADatasourceCanBeAddedToAGraph(t *testing.T) {

	g := NewGraph()
	d := NewMockDatasource()

	if e := g.AddDatasource(d); e != nil {
		t.Fatalf("g.AddDatasource: %s", e)
	}

	if g, e := len(g.datasourceReaders), 1; g != e {
		t.Errorf("Expected %d datasource readers, got %d", e, g)
	}

	if g, e := len(g.datasourceWriters), 1; g != e {
		t.Errorf("Expected %d datasource writers, got %d", e, g)
	}
}

func Test_SomeNonDatasourceTypeCantBeAddedToGraph(t *testing.T) {

	g := NewGraph()
	d := struct{}{}

	if e := g.AddDatasource(d); e == nil {
		t.Fatalf("Expected nill error, got '%s'", e)
	}
}
