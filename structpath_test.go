package inj

import "testing"

// smptystructPath should return an empty struct path
func Test_structPathA(t *testing.T) {

	if string(emptyStructPath()) != "" {
		t.Errorf("smptystructPath() didn't return empty string")
	}
}

// Branching should work intuitively
func Test_structPathB(t *testing.T) {

	inputs := []struct {
		start structPath
		add   string
		end   structPath
	}{
		{"", "Node", ".Node"},
		{".Node", "Node", ".Node.Node"},
		{".some thing", "else", ".some thing.else"},
	}

	for i, input := range inputs {
		sp := structPath(input.start)
		sp = sp.Branch(input.add)

		if g, e := sp, input.end; g != e {
			t.Errorf("[%d] Got %s, expected %s", i, g, e)
		}
	}
}

// Shifting should work intuitively
func Test_structPathC(t *testing.T) {

	inputs := []struct {
		start structPath
		str   string
		sp    structPath
	}{
		{"", "", ""},
		{".Node", "Node", ""},
		{".Node.NodeTwo", "Node", ".NodeTwo"},
		{".Parent.Node.NodeTwo", "Parent", ".Node.NodeTwo"},
	}

	for i, input := range inputs {
		sp := structPath(input.start)
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
func Test_structPathD(t *testing.T) {
	sp := structPath("")

	if !sp.Empty() {
		t.Errorf("structPath wasn't empty when it should be")
	}
}
