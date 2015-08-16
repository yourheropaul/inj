package inj

const (
	HELLO_SAYER_MESSAGE   = "Hello!"
	GOODBYE_SAYER_MESSAGE = "Bye!"
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

// Sample concrete type which requires two interfaces,
// the func type, the channel type and a string
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

// A nested type that contains depdendencies
type NestedType struct {
	Hello   InterfaceOne `inj:""`
	Goodbye InterfaceTwo `inj:""`
}

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
