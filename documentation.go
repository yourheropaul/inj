// Package inj provides reflection-based dependency injection for structs and functions. Some parts of it will be familiar to
// anyone who's ever used https://github.com/facebookgo/inject; others bear a passing similarity to depdendency injection in
// Angular.js. It's designed for medium to large applications, but it works just fine for small apps too.  It's especially
// useful if your project is is BDD/TDD-orientated.
//
// The essential premise of the package is that of object graphs. An object graph is just an index of objects organised by
// their explicitly nominated relationships. Depending on the application, a graph can be very simple - struct A depends on
// struct B - or extremely complicated. Package inj aims to simplify the process of creating and maintaining a graph of any
// complexity using out-of-the-box features of the Go language, and to povide easy access to the objects in the graph.
//
// A simple and unrealistically trivial example is this:
//
//  package main
//
//  import (
//      "fmt"
//      "github.com/yourheropaul/inj"
//  )
//
//  type ServerConfig struct {
//      Port int    `inj:""`
//      Host string `inj:""`
//  }
//
//  func main() {
//      config := ServerConfig{}
//      inj.Provide(&config, 6060, "localhost")
//
//      // The struct fields have now been set by the graph, and
//      // this will print "localhost:6060"
//      fmt.Printf("%s:%d",config.Host,config.Port)
//  }
//
// To understand what's happening there, you have to perform a bit of counter-intuitive reasoning. Inj has a global graph
// (which is optional, and can be disabled with the noglobals build tag) that is accessed through the three main API functions –
// inj.Provide(), inj.Assert() and inj.Inject(). The first and most fundamental of those functions - inj.Provide() - inserts anything
// passed to it into the graph, and then tries to wire up any dependency requirements.
//
// Before we get on to what dependency requirements are, it's important to know how a graph works. It is, in its simplest form, a map
// of types to values. In the example above, the graph - after inj.Provide() is called - will have three entries:
//
//  [ServerConfig] => (a pointer to the config variable)
//  [int] => 6060
//  [string] => "localhost"
//
// There can only be one entry for each type in a graph, but since Go allows new types to be created arbitrarly, that's not very much of
// a problem for inj. For example, if I wanted a second string-like type to be stored in the graph, I could simply type one:
//
//  package main
//
//  import (
//      "github.com/yourheropaul/inj"
//  )
//
//  type MyString string
//
//  func main() {
//      var basicString string = "this is a 'normal' string"
//      var myString MyString = "this is a typed string"
//      inj.Provide(basicString,myString)
//  }
//
// In that example, the graph would now contain two separate, stringer entities:
//
//  [string] => "this is a 'normal' string"
//  [MyString] => "this is a typed string"
//
// Back to depdendency requirements. The ServerConfig struct above has two, indicated by the inj:"" struct field tags (advanced usage of the
// package makes use of values in the tags, but we can ignore that for now). At the end of the inj.Provide() call, the graph is wired up –
// which essentially means finding values for all of the dependency requirements by type. The Port field of the ServerConfig struct requires
// and int, and the graph has one, so it's assigned; the Host field requires as string, and that can be assigned from the graph too.
//
// Obviously these examples are trivial in the extreme, and you'd probably never use the inj package in that way. The easiest way to understand
// the package for real-world applications is to refer to the example application: https://github.com/yourheropaul/inj/tree/master/example.
//
// For more general information, see the Wikipedia article for a technical breakdown of dependency injection:
// https://en.wikipedia.org/wiki/Dependency_injection.
//
package inj
