package inj

import (
	"fmt"
	"reflect"
	"testing"
)

//////////////////////////////////////////////
// Types
//////////////////////////////////////////////

type validConnectTester struct {
	Child1 *connectTesterChild1 `inj:""`
	Child2 *connectTesterChild2 `inj:""`
}

type invalidConnectTester struct {
	Child1 *connectTesterChild1 `inj:""`
	child2 *connectTesterChild2 `inj:""`
}

type connectTesterChild1 struct {
	Value string
}

type connectTesterChild2 struct {
	Value string
}

func newChildren() (*connectTesterChild1, *connectTesterChild2) {
	return &connectTesterChild1{"1"}, &connectTesterChild2{"2"}
}

//////////////////////////////////////////////
// Unit tests
//////////////////////////////////////////////

// Connected objects should be in the graph
func Test_ConnectHappyPath(t *testing.T) {

	g, p := NewGraph(), validConnectTester{}
	c1, c2 := newChildren()

	g.Provide(&p, c1, c2)

	// Provide calls connect, but call it again
	// explicitly
	g.connect()

	// Basic tests against injection failure
	if p.Child1 == nil {
		t.Errorf("Child1 is nil")
	}

	if p.Child2 == nil {
		t.Errorf("Child2 is nil")
	}

	// The main test is of the graph itself
	c1_found, c2_found := false, false

	for _, n := range g.nodes {

		if n.Object == c1 {
			c1_found = true
		}

		if n.Object == c2 {
			c2_found = true
		}
	}

	if !c1_found {
		t.Errorf("Didn't find c1")
	}

	if !c2_found {
		t.Errorf("Didn't find c2")
	}
}

// Unmet dependencies should be counted, and their
// errors stored
func Test_ConnectDepCount(t *testing.T) {

	g, p := NewGraph(), validConnectTester{}

	g.Provide(&p)

	g.connect()

	if g, e := g.UnmetDependencies, 2; g != e {
		t.Errorf("Got %d unmet deps, expected %d", g, e)
	}

	if g, e := len(g.Errors), 2; g != e {
		t.Errorf("Got %d unmet dep errors, expected %d", g, e)
	}
}

// Values should actually be assigned
func Test_ConnectAssignmentHappyPath(t *testing.T) {

	g, p := NewGraph(), &validConnectTester{}
	c1, c2 := newChildren()
	v := reflect.ValueOf(p)

	gnds := []graphNodeDependency{
		graphNodeDependency{
			Path: structPath(".Child1"),
			Type: reflect.TypeOf(c1),
		},
		graphNodeDependency{
			Path: structPath(".Child2"),
			Type: reflect.TypeOf(c2),
		},
	}

	g.Provide(c1, c2)

	for _, gnd := range gnds {
		if err := g.assignValueToNode(v, gnd); err != nil {
			t.Errorf("assignValueToNode: %s", err.Error())
		}
	}

	// Basic tests against injection failure
	if p.Child1 == nil {
		t.Errorf("Child1 is nil")
	}

	if p.Child2 == nil {
		t.Errorf("Child2 is nil")
	}
}

// Values should not be assigned to nodes which already
// have values
func Test_ConnectAssignmentNeutralPath(t *testing.T) {

	g, p := NewGraph(), &validConnectTester{}
	c1, c2 := newChildren()
	c3, c4 := newChildren()
	v := reflect.ValueOf(p)

	// Manually assign the deps
	p.Child1 = c1
	p.Child2 = c2

	gnds := []graphNodeDependency{
		graphNodeDependency{
			Path: structPath(".Child1"),
			Type: reflect.TypeOf(c1),
		},
		graphNodeDependency{
			Path: structPath(".Child2"),
			Type: reflect.TypeOf(c2),
		},
	}

	// Assign everything
	g.Provide(c3, c4)

	// Run through and re-assign (shouldn't error)
	for _, gnd := range gnds {
		if err := g.assignValueToNode(v, gnd); err != nil {
			t.Errorf("assignValueToNode: %s", err.Error())
		}
	}

	if p.Child1 != c1 {
		t.Errorf("Child1 isn't the original c1")
	}

	if p.Child2 != c2 {
		t.Errorf("Child1 isn't the original c1")
	}
}

// Internal variabled should cause the assignment to fail
func Test_ConnectSadPath1(t *testing.T) {

	g, p := NewGraph(), &invalidConnectTester{}
	c1, c2 := newChildren()
	v := reflect.ValueOf(p)

	gnds := []graphNodeDependency{
		graphNodeDependency{
			Path: structPath(".child2"),
			Type: reflect.TypeOf(c2),
		},
	}

	g.Provide(c1, c2)

	// Run through and assign (should error)
	for _, gnd := range gnds {
		if err := g.assignValueToNode(v, gnd); err == nil {
			t.Errorf("assignValueToNode: didn't error")
		}
	}
}

// Unmet dependencies should cause an error
func Test_ConnectSadPath2(t *testing.T) {

	g, p := NewGraph(), &validConnectTester{}
	c1, c2 := newChildren()
	v := reflect.ValueOf(p)

	gnds := []graphNodeDependency{
		graphNodeDependency{
			Path: structPath(".Child1"),
			Type: reflect.TypeOf(c1),
		},
		graphNodeDependency{
			Path: structPath(".Child2"),
			Type: reflect.TypeOf(c2),
		},
	}

	// Run through and assign (should error)
	for _, gnd := range gnds {
		if err := g.assignValueToNode(v, gnd); err == nil {
			t.Errorf("assignValueToNode: didn't error")
		}
	}
}

// Should find a reflect value for a path
func Test_ConnectFindFieldValue(t *testing.T) {

	g, p := NewGraph(), &validConnectTester{}
	v := reflect.ValueOf(p)

	var descs = []struct {
		fieldName string
		path      structPath
	}{
		{
			fieldName: "*inj.connectTesterChild1",
			path:      structPath(".Child1"),
		},
		{
			fieldName: "*inj.connectTesterChild2",
			path:      structPath(".Child2"),
		},
	}

	for _, d := range descs {
		rv, e := g.findFieldValue(v, d.path, &[]reflect.Value{})

		if e != nil {
			t.Errorf("findFieldValue: %s", e.Error())
		}

		if g, e := rv.Type().String(), d.fieldName; g != e {
			t.Errorf("fieldname: got %s, expected %s", g, e)
		}
	}
}

// Error should be returned if input reflect value isn't
// a struct
func Test_ConnectF(t *testing.T) {

	g := NewGraph()

	_, e := g.findFieldValue(reflect.ValueOf("123"), ".Child1", &[]reflect.Value{})

	if e == nil {
		fmt.Errorf("Didn't error when type wasn't struct")
	}
}

// Incorrect paths should cause an error to be returned
func Test_ConnectG(t *testing.T) {

	g, p := NewGraph(), &validConnectTester{}

	_, e := g.findFieldValue(reflect.ValueOf(p), ".This.Doesnt.Exist", &[]reflect.Value{})

	if e == nil {
		fmt.Errorf("Didn't error when path was wrong")
	}
}
