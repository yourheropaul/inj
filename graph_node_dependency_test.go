package inj

import (
	"reflect"
	"testing"
)

func compareGraphNodeDeps(d1 graphNodeDependency, d2 graphNodeDependency, t *testing.T) {

	if d1.Path != d2.Path {
		t.Errorf("compareGraphNodeDeps: paths don't match (%s,%s)", d1.Path, d2.Path)
	}

	if d1.Type != nil && d1.Type != d2.Type {
		t.Errorf("compareGraphNodeDeps: types don't match (%s,%s)", d1.Type, d2.Type)
	}
}

// findDependencies should populate a known number of deps
func Test_FindDependencies(t *testing.T) {

	c := ConcreteType{}
	d := make([]graphNodeDependency, 0)
	s := emptyStructPath()

	if e := findDependencies(reflect.TypeOf(c), &d, &s); e != nil {
		t.Errorf("Unexpected error: %s", e)
	}

	eds := c.expectedDeps()

	if g, e := len(d), len(eds); g != e {
		t.Errorf("Expected %d deps in c, got %d", e, g)
	}

	for i, ed := range eds {
		compareGraphNodeDeps(ed, d[i], t)
	}
}

// findDependencies should populate a known number of deps
func Test_FindDependenciesInEmbeddedStructs(t *testing.T) {

	c := HasEmbeddable{}
	d := make([]graphNodeDependency, 0)
	s := emptyStructPath()

	if e := findDependencies(reflect.TypeOf(c), &d, &s); e != nil {
		t.Errorf("Unexpected error: %s", e)
	}

	eds := c.expectedDeps()

	if g, e := len(d), len(eds); g != e {
		t.Errorf("Expected %d deps in c, got %d", e, g)
	}

	for i, ed := range eds {
		compareGraphNodeDeps(ed, d[i], t)
	}
}

// parseStructTag should return expected struct tags
func Test_ParseStructTag(t *testing.T) {

	inps := []struct {
		tag            reflect.StructTag
		expectedValues []string
	}{
		{
			tag: "inj:\"\" someothertag:\"\"",
		},
		{
			tag:            "inj:\"some.datasource.path\"",
			expectedValues: []string{"some.datasource.path"},
		},
		{
			tag:            "inj:\"one,two\"",
			expectedValues: []string{"one", "two"},
		},
	}

	for _, inp := range inps {
		d := parseStructTag(inp.tag)

		if !reflect.DeepEqual(inp.expectedValues, d.DatasourcePaths) {
			t.Errorf("inp.expectedValues != d.DatasourcePaths")
		}
	}
}
