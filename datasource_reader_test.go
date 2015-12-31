package inj

import "testing"

func Test_TheDatasourceReaderWritesToTheDepdendency(t *testing.T) {

	dep := dataSourceDep{}
	ds := newMockDataSourceWithValues(t)
	g := NewGraph()

	g.AddDatasource(ds)
	g.Provide(&dep)

	if len(g.Errors) != 0 {
		t.Errorf("Graph was initialised with errors > 0")
	}

	if g.UnmetDependencies != 0 {
		t.Errorf("Graph was initialised with UnmetDependencies > 0")
	}

	if g, e := dep.StringValue, DEFAULT_STRING; g != e {
		t.Errorf("Expected string '%s', got '%s'", e, g)
	}

	if dep.FuncValue == nil {
		t.Errorf("Didn't get expected function instance")
	}
}
