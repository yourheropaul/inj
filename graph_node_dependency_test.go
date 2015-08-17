package inj

import (
	"reflect"
	"testing"
)

func compareGraphNodeDeps(d1 GraphNodeDependency, d2 GraphNodeDependency, t *testing.T) {

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
	d := make([]GraphNodeDependency, 0)
	s := EmptyStructPath()

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
		tag         reflect.StructTag
		defaultName string
		name        string
	}{
		{
			tag:         "inj:\"\" someothertag:\"\"",
			defaultName: "default",
			name:        "default",
		},
	}

	for _, inp := range inps {
		d := parseStructTag(inp.tag, inp.defaultName)

		if g, e := d.Name, inp.name; g != e {
			t.Errorf("%s: got name %s, expected %s", g, e)
		}
	}
}
