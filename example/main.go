/*
Package inj/example is a demonstration the inj package. It takes the form of a very simple web application
that is modelled by a struct. Depdendencies for the application are defined on struct fields using tags,
and provided in the main function. On startup, the application will provide and check its dependencies,
then start a simple HTTP server with two endpoints.

The functions for the endpoints use a rudimentary HTTP middleware wrapped so they can make use of inj.Inject,
and thus have non-standard argument requirements. That will make sense when you see it in action.

Hopefully the code is reasonably self-explanatory. To experience the full utility of inj, compare main.go to
main_test.go.
*/
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/yourheropaul/inj"
)

/*
This is the struct that describes the application. The struct has one function (`run()`, which you'll see below)
and three fields that represent things it will need for the demo. None of the fields will be populates explicitly:
rather, they will be provided using dependency injection.

As you can probably see, depdendency requirements are indicated using an `inj:""` struct tag. Fields with this tag
will be assigned, if possible, during the provision call; fields without the tag will be ignored.

There are currently no values to pass in the struct tag. Future versions may make more use of them, but for the
moment all tags must pass strings in double quotes to satisfy Go's reflect.StructTag parsing.

It should be noted that any field can be associatd with the tag. The field doesn't have to be a special type, and
there are no limits (aside from the constraints of the language itself) as to what they can be.

`Application`'s functions (there are three) are at the end of this file.
*/
type Application struct {
	Config Configurer `inj:""`
	Log    Logger     `inj:""`
	Exit   ExitChan   `inj:""`
}

/////////////////////////////////////////////
// Application types
/////////////////////////////////////////////

/*
For the purposes of this app, the `Configurer` type is very simple: it's an interface that stipulates exactly one
function. The return value of the the implementing object will used as the port number in the HTTP server.
*/
type Configurer interface {
	Port() int
}

/* A logger is a simple function that happens to conform to `fmt.Printf`'s signature. */
type Logger func(string, ...interface{}) (int, error)

/* This is a basic implementation of the Configurer interface. */
type Config struct{}

func (c *Config) Port() int { return 8080 }

/* The final struct type, ExitChan, is just a channel that accepts blank interfaces */
type ExitChan chan interface{}

/////////////////////////////////////////////
// Other types
/////////////////////////////////////////////

/*
One of the application behaviours is write to an HTTP client. For code readability, this is a type that can do
that.
*/
type Responder func(w http.ResponseWriter, r *http.Request)

/* ... and here's a trivial implementation of that type. */
func WriteResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

/////////////////////////////////////////////
// Entrypoint
/////////////////////////////////////////////

func main() {

	app := Application{}

	/*
		This is the first of three `inj` API calls. (The others are `inj.Assert()` and `inj.Inject`, more on those
		later). `inj.Provide()` is a variadic function that takes any number of `interface{}`s and inserts then into
		a list of globally-available objects call the application graph.

		After each call to `inj.Provide()`, the graph is 'connected', meaning that any structs with depdendencies (that
		is, structs containing fields with `inj:""` tags) have their dependent fields assigned to nodes on the graph.
		This how how we'll assign values to the three fields of the `Application` struct.

		You can call `inj.Provide()` as many times as you like, which means that it's possible to build a cascade of
		dependencies from imported packages using the `init()` function. Be careful with that, because `inj.Provide()`
		is where all the heavy lifting is done, and it's not fast enough to execute for each request or operation.
		Generally it should be called at startup, like we're doing here.
	*/
	inj.Provide(

		// A pointer to our `Application` is necessary to fulfil its depdendencies. (It must be a pointer so the values
		// can be assigned).
		&app,

		// Objects are stored in the graph, even if they're stored anywhere else. Here we're creating a channel.
		make(ExitChan),

		// The same with the `Config` struct; just create a new instance for the graph and move on.
		&Config{},

		// As noted above, `Logger` happens to match the signature of `fmt.Printf`, so we use that.
		fmt.Printf,

		// This is that 'Hello, I love %s' writer. `Application` doesn't use it, but other graph-dependent functions do.
		WriteResponse,
	)

	/*
		Once you've made all your calls to `inj.Provide()`, it's time to make sure that all the dependencies were met. In
		`inj`, that's fairly simple: use the `inj.Assert()` function. It returns a boolean and a slice of strings. If the
		boolean is false, then the slice will be populated with all the errors that occured during struct injection.
	*/
	if valid, messages := inj.Assert(); !valid {
		fmt.Println(messages)
		os.Exit(1)
	}

	/*
		If we get to this point, we know that `app` has had all its dependencies met, so we can call `app.run()`.
	*/
	app.run()
}

/////////////////////////////////////////////
// A coda: middleware and `inj.Inject()`
/////////////////////////////////////////////

/*
This application uses `http.HandleFunc` to, well, handle HTTP requests using a function. That's because I needed to find a
vaguely realistic way of demonstrating the third and last API function, `inj.Inject()`. As it turns out, it's quite a nice
use of the `inj` API, but I realise you might not use this in your own applications.

`inj.Inject()` is a variadic function that requires a function (any kind of function) as its first argument, and then zero
or more other objects as it variadic component. Under the hood, it examines the arguments of the function, tries to find
values for each one from the graph (or, if they're not available in the graph, the variadic arguments to inj.Inject()`
itself) and then calls the function.

This can be a little bit confusing to think about to begin with, and HTTP middleware (which is typically a function that
returns a function) doesn't necessarily clear things up. It might be easier to visualise it like this:

	inj.Provide("This string value is in the graph")
	inj.Inject(
		func(s string, i int) {
			fmt.Sprintf("String value: %s, int value: %d", s, i)
		},
		10
	)

The code above inserts a string value into the graph and then `inj.Inject()` to call a function. The function is called
immediately, and prints some text to stdout: "String value: This string value is in the graph, int value: 10".

The values come from the graph, and the variadic arguments to `inj.Inject()`. They're specified by the arguments to function
itself. When `inj.Inject()` is called, it examines the each argument in turn and checks if it has a value for it. Because
there's a string in the graph, the string argument is injectable; because there's an integer in the call to `inj.Inject()`
the integer agument is injectable. The order of the arguments doesn't matter. If an argument is not available, then a runtime
panic will be thrown.

It's not just for callback functions. Closures in Go mean that you can chuck in a call to `inj.Inject()` anywhere, and pull
objects from the graph. Imagine you have an instance of a database accessor in your graph:

	inj.Provide(databaseStruct{})

Now say you have some business logic function that needs access to the database accessor, but it has no way to directly find it.
You can use a simple nested closure with `inj.Inject()` to assemble all the values you need:

	func (b *SomeBusinessLogicController) DoSomethingImportant(ctx *business.Context, p *ptrToThing) {
		inj.Inject(func(db DatatabaseInterface) {
			// Now `ctx`, `p` AND `db` are all in scope, and you didn't have to change the arguments to your function
		})
	}

I don't want to brag, but that's pretty nifty, right?

Getting back to the demo, the  `middleware()` function returns another function that makes use of `inj.Inject()`.
Everything in the graph will be available to the the `fn` function, as well as the `http.ResponseWriter` and `*http.Request`
arguments that are passed by the http package. You'll see implementations for the handle functions at the end of this file,
and you'll notice that they can specify any arguments they like, not just the standard ones you'd expect. This is the really
nifty part of the `inj` package: callback functions can decide which variables they want dynamically.
*/

func middleware(fn interface{}) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		inj.Inject(fn, w, r)
	}
}

/////////////////////////////////////////////
// Application functions
/////////////////////////////////////////////

/* This is the function that's called at the end of `main()` */
func (a Application) run() {

	// The handler functions use the middleware as outlined above.
	http.HandleFunc("/", middleware(a.handler))
	http.HandleFunc("/shutdown", middleware(a.shutdown))

	// The `Config` type (which implements `Configurer`) will have been automatically injected as part of the initial
	// `inj.Provide()` call in `main()`.
	port := a.Config.Port()

	// The `Logger` will also have been automatically injected.
	a.Log("Running server on port %d...\n", port)

	// This has nothing to do with `inj`; it's a standard library call
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	// Remember the `ExitChan` type that was created and provided in `main()`? This is what it's for.
	<-a.Exit

	a.Log("Received data on exit channel.\n")
}

/*
This is the "/" handler, which is wrapped using `middleware()`, as outlined above. Its purpose in the application is to
route requests through the `WriteResponse` ("Hello! I love %s") function defined near the top of the file. That function
is a `Responder` type, and it happens to be in the graph because it was part of the `inj.Provide()` call in `main()`.

This callback is wrapped by `inj.Inject()`, so it can specify whatever arguments it wants. Here, we're requesting the
`w http.ResponseWriter` and `r *http.Request` variables that the `http` package sends, and also a `Responder` from the
graph.
*/
func (a Application) handler(w http.ResponseWriter, r *http.Request, responder Responder) {
	responder(w, r)
}

/*
Here's another handler, registed on "/shutdown". It calls for the `w http.ResponseWriter` (but not the `r *http.Request`),
and also an `ExitChan`.

(Yes, I know the ExitChan is already in the `a` variable. This is a demo application. It doesn't have to make complete,
real-world sense. The function illustrates that you can omit normally-required functions and add others, and that's all
it needs to do).
*/
func (a Application) shutdown(w http.ResponseWriter, e ExitChan) {
	w.Write([]byte("Shutting down!"))
	a.Log("Shutting down...\n")
	e <- 0
}

/*
That's pretty much it! To recap:

- There are three API functions in the `inj` package: `inj.Provide()`, `inj.Assert()` and `inj.Inject()`.
- Use `inj.Provide()` any number of times to initially wire up your graph
- Then use `inj.Assert()` to make sure all the dependencies were met.
- If you want to use dynamic function values, use `inj.Inject()`, which can also accept local arguments.

Do check out main_test.go, which demonstrates how to use `inj` for testing.
*/
