package inj

import (
	"fmt"
	"testing"
)

///////////////////////////////////////////////////////////////
// A mock DatasourceReader implementation
///////////////////////////////////////////////////////////////

type MockDatasourceReader struct {
	stack map[string]interface{}
}

func NewMockDatasourceReader(data ...map[string]interface{}) *MockDatasourceReader {

	d := &MockDatasourceReader{}

	d.stack = make(map[string]interface{})

	for _, datum := range data {
		for k, v := range datum {
			d.stack[k] = v
		}
	}

	return d
}

func (d *MockDatasourceReader) Read(key string) (interface{}, error) {

	if value, exists := d.stack[key]; exists {
		return value, nil
	}

	return nil, fmt.Errorf("No stack entry for '%s'", key)
}

///////////////////////////////////////////////////////////////
// Unit tests for graph implementation
///////////////////////////////////////////////////////////////

func Test_TheDatasourceReaderWritesToTheDepdendency(t *testing.T) {

	dep := dataSourceDep{}
	ds := newMockDataSourceWithValues(t)
	g := NewGraph()

	g.AddDatasource(ds)
	g.Provide(&dep)

	assertNoGraphErrors(t, g.(*graph))

	if g, e := dep.StringValue, DEFAULT_STRING; g != e {
		t.Errorf("Expected string '%s', got '%s'", e, g)
	}

	if dep.FuncValue == nil {
		t.Errorf("Didn't get expected function instance")
	}
}
