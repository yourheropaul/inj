[![Build Status](https://travis-ci.org/yourheropaul/inj.svg?branch=master)](https://travis-ci.org/yourheropaul/inj) [![Coverage Status](https://coveralls.io/repos/yourheropaul/inj/badge.svg?branch=master&service=github)](https://coveralls.io/github/yourheropaul/inj?branch=master) [![GoDoc](https://godoc.org/github.com/yourheropaul/inj?status.svg)](https://godoc.org/github.com/yourheropaul/inj) [![Example](https://img.shields.io/badge/code-example-2a988f.svg)](https://github.com/yourheropaul/inj/tree/master/example) [![YHP](https://img.shields.io/badge/author-YHP-f08332.svg)](http://yhp.io) [![Slashie](https://img.shields.io/badge/ridiculously-goodlooking-736caf.svg)](http://yhp.io)
======

Are you troubled by dependency configuration issues in the middle of the night? Do you experience feelings of dread at the code boundaries of your application or during functional testing? Have you or your team ever used a global variable, exported package-level object or a build constraint hack? If the answer is *yes*, then don't wait another minute. Pick up your terminal and `go get github.com/yourheropaul/inj` today.

Inject yourself before you wreck yourself.

### What *is* this thing?

`inj` provides reflection-based dependency injection for Go structs and functions. Some parts of it will be familiar to anyone who's ever used [facebookgo/inject](https://github.com/facebookgo/inject); others bear a passing similarity to [dependency injection in Angular.js](https://docs.angularjs.org/guide/di).  It's designed for medium to large applications, but it works just fine for small apps too. It's especially useful if your project is is BDD/TDD-orientated.

### Dependency injection is boring and hard to maintain. I want something that works out of the box.

`inj` may well be for you. A good, idiomatic Go package is a simple, clearly laid out, well-tested Go package – and `inj` is, or aspires to be, all of those things. It works well and it works intuitively. There's only one, optional configuration option (whether or not to use a global variable inside the package) and there's very little to think about once you're set up.

### How do I use it?

A simple and unrealistically trivial example is this:

```
package main

import (
	"fmt"

	"github.com/yourheropaul/inj"
)

type ServerConfig struct {
	Port int    `inj:""`
	Host string `inj:""`
}

func main() {
	config := ServerConfig{}
	inj.Provide(&config, 6060, "localhost")

	// The struct fields have now been set by the graph, and
	// this will print "localhost:6060"
	fmt.Printf("%s:%d", config.Host, config.Port)
}
```
There's a full explanation for this basic example in the [Godoc](https://godoc.org/github.com/yourheropaul/inj). 
Obviously this example is trivial in the extreme, and you'd probably never use the the package in that way. The easiest way to understand
 `inj` for real-world applications is to refer to the [example application](https://github.com/yourheropaul/inj/tree/master/example) in this repository. The API is small, and everything in the core API is demonstrated there. 
### Dependency injection is great and everything, but I really want to be able to pull data directly from external services, not just the object graph. 
 
You mean you want to read from a JSON or TOML config file, and inject the values into Go objects directly? Maybe you'd like to pull values from a DynamoDB instance and insert them into Go struct instances with almost zero code overhead?

That's what's `inj` is designed for! And what's more, intrepid programmer [Adrian Duke](http://adeduke.com/) has already done the leg work for you in his fantastic [configr](https://github.com/adrianduke/configr) package – see his readme for brief instructions.

### I want absolutely, positively no globals in my application. None. Can I do that with this package?

Of course, and it couldn't be easier! Just compile your application with the tag `noglobals`, and the package-level API functions (including the one package-level variable they use) won't be included. You can create a new graph for your application by calling `inj.NewGraph()`, which has the same functional interface as the package API.

### This whole thing sounds too useful to be true

I appreciate your skepticism, so let's gather some data. There are two things you need to be aware of when using `inj`.

The first is that the application graph is one dimensional and indexed by type. That means you can't call `inj.Provide(someIntValue,someOtherIntValue)` and expect both integers to be in the graph – the second will override the first. Other dependency injection approaches allow for named  and private depdendencies, but that has been sacrified here so that `inj.Inject()` could exist in a consistent way. When there's a choice, `inj` errs on the side of simplicity, consistency and idiomatic implementation over complexity and magic.

The second consideration is execution speed. Obviously, calling `inj.Inject(fn)` is slower than calling `fn()` directly. In Go 1.4, with a medium-sized graph, it takes about 350 times longer to execute the call; in Go 1.5 rc1, it's about 240 times. If those numbers seem high, it's because they are. The impact on an application is measurable, but for most purposes negligible. 

If the average execution time of a pure Go function is around 4 nanoseconds (as it is in my tests) then the execution time of `inj.Inject()` will be somewhere between 900 and 1,400 nanoseconds. Or in more useful terms, 0.0014 milliseconds (which is 1.4e-6 seconds). If your application is built for speed, then you will need to be judicious in your use of `inj.Inect()`. Even if speed isn't a concern, it's generally not a good idea to nest injection calls, or put them in loops.

Finally, `inj.Provide()` is fairly slow, but it's designed to executed at runtime only. There are benchmark tests in the package if you want to see how it performs on your system.

### But how do I use it?

Seriously? I just explained that a minute ago. Maybe look at the [example application](https://github.com/yourheropaul/inj/tree/master/example) or the [Godoc](https://godoc.org/github.com/yourheropaul/inj).
