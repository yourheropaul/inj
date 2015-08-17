package inj

import "testing"

const (
	HELLO_SAYER_MESSAGE   = "Hello!"
	GOODBYE_SAYER_MESSAGE = "Bye!"
	DEFAULT_STRING        = "this is a string"
)

///////////////////////////////////////////////////
// Types for the unit and feature tests
///////////////////////////////////////////////////

type InterfaceOne interface {
	SayHello() string
}

type InterfaceTwo interface {
	SayGoodbye() string
}

type FuncType func(string) string
type ChanType chan interface{}

///////////////////////////////////////////////////
// Sample concrete type which requires two interfaces,
// the func type, the channel type and a string
///////////////////////////////////////////////////

type ConcreteType struct {
	Hello    InterfaceOne `inj:""`
	Goodbye  InterfaceTwo `inj:""`
	Stringer FuncType     `inj:""`
	Channel  ChanType     `inj:""`
	String   string       `inj:""`

	// This is nested
	Nested NestedType

	// These are not included in the injection
	Something     string `in:`
	SomethingElse int
}

// A nested type that contains dependencies
type NestedType struct {
	Hello   InterfaceOne `inj:""`
	Goodbye InterfaceTwo `inj:""`
}

// Channel instance
var ichannel = make(ChanType)

///////////////////////////////////////////////////
// Implementation of a hello-sayer
///////////////////////////////////////////////////

type helloSayer struct{}

func (g *helloSayer) SayHello() string { return HELLO_SAYER_MESSAGE }

///////////////////////////////////////////////////
// Implementation of a goodbye-sayer
///////////////////////////////////////////////////

type goodbyeSayer struct{}

func (g *goodbyeSayer) SayGoodbye() string { return GOODBYE_SAYER_MESSAGE }

///////////////////////////////////////////////////
// Implementation of a FuncType
///////////////////////////////////////////////////

func funcInstance(s string) string {
	return s
}

//////////////////////////////////////////
// Assertion for concrete type
//////////////////////////////////////////

// Once the dependencies have been injected, all the dependent
// members should be non-nil and functional.
func assertConcreteValue(c ConcreteType, t *testing.T) {

	if c.Hello == nil {
		t.Errorf("c.Hello is nil")
	}

	if c.Goodbye == nil {
		t.Errorf("c.Goodbye is nil")
	}

	if c.Stringer == nil {
		t.Errorf("c.Stringer is nil")
	}

	if c.Channel == nil {
		t.Errorf("c.Channel is nil")
	}

	if c.String == "" {
		t.Errorf("c.String is nil")
	}

	if c.Nested.Hello == nil {
		t.Errorf("c.Hello is nil")
	}

	if c.Nested.Goodbye == nil {
		t.Errorf("c.Goodbye is nil")
	}

	if g, e := c.Hello.SayHello(), HELLO_SAYER_MESSAGE; g != e {
		t.Errorf("i2.SayHello(): got %s, expected %s", g, e)
	}

	if g, e := c.Goodbye.SayGoodbye(), GOODBYE_SAYER_MESSAGE; g != e {
		t.Errorf("i2.SayHello(): got %s, expected %s", g, e)
	}

	// test the function
	if g, e := c.Stringer(DEFAULT_STRING), DEFAULT_STRING; g != e {
		t.Errorf("Test Stringer: got %s, expected %s", g, e)
	}
}
