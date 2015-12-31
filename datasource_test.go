package inj

import (
	"fmt"
	"testing"
)

///////////////////////////////////////////////////////////////
// A mock DatasourceReader and DatasourceWriter implementation
///////////////////////////////////////////////////////////////

type MockDatasource struct {
	stack map[string]interface{}
}

func NewMockDatasource() Datasource {

	d := &MockDatasource{}

	d.stack = make(map[string]interface{})

	return d
}

func (d *MockDatasource) Read(key string) (interface{}, error) {

	if value, exists := d.stack[key]; exists {
		return value, nil
	}

	return nil, fmt.Errorf("No stack entry for '%s'", key)
}

func (d *MockDatasource) Write(key string, value interface{}) error {

	d.stack[key] = value

	return nil
}

///////////////////////////////////////////////////////////////
// A mock specific dependency implementation for testing
///////////////////////////////////////////////////////////////

type dataSourceDep struct {
	StringValue string   `inj:"datasource.string"`
	FuncValue   FuncType `inj:"datasource.func"`
}

func newMockDataSourceWithValues(t *testing.T) Datasource {

	d := NewMockDatasource()

	if e := d.Write("datasource.string", DEFAULT_STRING); e != nil {
		t.Fatalf("newMockDataSourceWithValues: Datasource.Write: %s", e)
	}

	if e := d.Write("datasource.func", funcInstance); e != nil {
		t.Fatalf("newMockDataSourceWithValues: Datasource.Write: %s", e)
	}

	return d
}
