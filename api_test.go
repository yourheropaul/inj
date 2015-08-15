package inj

///////////////////////////////////////////////////
// Types for the feature test
///////////////////////////////////////////////////

type InterfaceOne interface {
	SayHello() string
}

type InterfaceTwo interface {
	SayGoodbye() string
}

type FuncType func(string) string
type ChanType chan interface{}

// Sample concrete type which requires two interfaces
type ConcreteType struct {
	Hello    InterfaceOne `inj:""`
	Goodbye  InterfaceTwo `inj:",private"`
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
	Goodbye InterfaceTwo `inj:"`
}

// Implementation of a hello-sayer
type helloSayer struct{}

func (g *helloSayer) SayHello() string { return "hello!" }

// Implementation of a goodbye-sayer
type goodbyeSayer struct{}

func (g *goodbyeSayer) SayGoodbye() string { return "bye!" }

// Implementation of a FuncType
func funcInstance(s string) string {
	return s
}
