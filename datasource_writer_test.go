package inj

import "testing"

///////////////////////////////////////////////////////////////
// A mock DatasourceWriter implementation
///////////////////////////////////////////////////////////////

type MockDatasourceWriter struct {
	stack map[string]interface{}
}

func NewMockDatasourceWriter(data ...map[string]interface{}) *MockDatasourceWriter {

	d := &MockDatasourceWriter{}

	d.stack = make(map[string]interface{})

	for _, datum := range data {
		for k, v := range datum {
			d.stack[k] = v
		}
	}

	return d
}

func (d *MockDatasourceWriter) Write(key string, value interface{}) error {

	d.stack[key] = value

	return nil
}

func (d *MockDatasourceWriter) Assert(t *testing.T, key string, value interface{}) {

	v, exists := d.stack[key]

	if !exists {
		t.Fatalf("MockDatasourceWriter.Assert: Key '%s' doesn't exist", key)
	}

	if v != value {
		t.Fatalf("MockDatasourceWriter.Assert: key %s doesn't match (%v, %v)", key, v, value)
	}
}

func (d *MockDatasourceWriter) AssertMap(t *testing.T, data map[string]interface{}) {

	for k, v := range data {
		d.Assert(t, k, v)
	}
}

///////////////////////////////////////////////////////////////
// Unit tests for graph implementation
///////////////////////////////////////////////////////////////

type datasourceWriterDep struct {
	IntVal    int    `inj:"datasource.writer.int"`
	StringVal string `inj:"datasource.writer.string"`
}

func Test_DatasourceReaderWriterLoopInGraph(t *testing.T) {

	expected_values := map[string]interface{}{
		"datasource.writer.int":    10,
		"datasource.writer.string": DEFAULT_STRING,
	}

	reader := NewMockDatasourceReader(expected_values)
	writer := NewMockDatasourceWriter()
	dep := datasourceWriterDep{}
	g := NewGraph()

	g.AddDatasource(reader, writer)
	g.Provide(&dep)

	assertNoGraphErrors(t, g)

	writer.AssertMap(t, expected_values)
}

func Test_DatasourceWriterWritesWithoutAReader(t *testing.T) {

	expected_values := map[string]interface{}{
		"datasource.writer.int":    10,
		"datasource.writer.string": DEFAULT_STRING,
	}

	writer := NewMockDatasourceWriter()
	dep := datasourceWriterDep{}
	g := NewGraph()

	g.AddDatasource(writer)
	g.Provide(
		expected_values["datasource.writer.int"],
		expected_values["datasource.writer.string"],
		&dep,
	)

	assertNoGraphErrors(t, g)

	writer.AssertMap(t, expected_values)
}
