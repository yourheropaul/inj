package inj

import "testing"

// EmptyStructPath should return an empty struct path
func Test_StructPathA(t *testing.T) {

	if string(EmptyStructPath()) != "" {
		t.Errorf("EmptyStructPath() didn't return empty string")
	}
}

// Branching should work intuitively
func Test_StructPathB(t *testing.T) {

	inputs := []struct {
		start StructPath
		add   string
		end   StructPath
	}{
		{"", "Node", ".Node"},
		{".Node", "Node", ".Node.Node"},
		{".some thing", "else", ".some thing.else"},
	}

	for i, input := range inputs {
		sp := StructPath(input.start)
		sp = sp.Branch(input.add)

		if g, e := sp, input.end; g != e {
			t.Errorf("[%d] Got %s, expected %s", i, g, e)
		}
	}
}

// Shifting should work intuitively
func Test_StructPathC(t *testing.T) {

	inputs := []struct {
		start StructPath
		str   string
		sp    StructPath
	}{
		{"", "", ""},
		{".Node", "Node", ""},
		{".Node.NodeTwo", "Node", ".NodeTwo"},
		{".Parent.Node.NodeTwo", "Parent", ".Node.NodeTwo"},
	}

	for i, input := range inputs {
		sp := StructPath(input.start)
		s, sp2 := sp.Shift()

		if g, e := s, input.str; g != e {
			t.Errorf("[%d] Got string %s, expected %s", i, g, e)
		}

		if g, e := sp2, input.sp; g != e {
			t.Errorf("[%d] Got path %s, expected %s", i, g, e)
		}
	}
}

// Emptying a structpath should cause the Empty() func to
// return true
func Test_StructPathD(t *testing.T) {
	sp := StructPath("")

	if !sp.Empty() {
		t.Errorf("StructPath wasn't empty when it should be")
	}
}
