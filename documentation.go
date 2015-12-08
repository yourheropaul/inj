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
// The easiest way to understand the inj package is to refer to the
// example application: https://github.com/yourheropaul/inj/tree/master/example.
//
// For more general information, see the // Wikipedia article for a technical breakdown of dependency injection:
// https://en.wikipedia.org/wiki/Dependency_injection.
//
package inj
