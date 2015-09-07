package main

/*
The test component of the inj/example application demonstrates the use of mock objects by
constructing a special graph of test objects that conforms to the application requirements,
and then runs the application.

Run this package's tests in verbose mode (`go test -test.v`) to get the full experience.
*/

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourheropaul/inj"
)

const (
	MOCK_HTTP_PORT = 8080
)

// This a mock implementation of the `Configurer` interface. The only real difference is that
// it reads its port from a const, which is also used later in the test suite.
type MockConfig struct{}

func (c *MockConfig) Port() int { return MOCK_HTTP_PORT }

// This fills in for the logging function. The test graph includes a pointer to a testing.T
// object, which is used by the mock logger. (These messages will only be visible if go test
// is running in verbose mode. See http://golang.org/pkg/testing/#T.Logf)
func MockLogger(s string, a ...interface{}) (int, error) {

	inj.Inject(func(t *testing.T) {
		t.Logf(s, a...)
	})

	return 0, nil
}

// We're also going to need a mock `Responder` (the function that the root HTTP hander
// invokes in the application). The application always responds with "Hi there, I love X",
// but the test can be simpler: it can just echo out whatever was added to the URL.
func MockWriteResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", r.URL.Path[1:])
}

/*
The first test assembles the components in an isolated graph so that we don't
contaminate the global graph during the scope of the test. All it really does
is `Provide()`s the required objects and asserts the intergrity of the graph.
*/
func Test_GraphInitialisation(t *testing.T) {

	g := inj.NewGraph()

	g.Provide(

		// An instance of the application object. There's no need to keep a
		// pointer to it, since we won't be using it
		&Application{},

		// The exit channel is the same as the non-test version
		make(ExitChan),

		// Our mock config will fill the `Configurer` dependency
		&MockConfig{},

		// Use the mock `Logger`
		MockLogger,

		// The mock `Responder`
		MockWriteResponse,

		// Finally, add a pointer to a `testing.T`. It's not used in this test,
		// but it's worth including for completeness.
		t,
	)

	// Make sure everything was loaded properly
	if valid, messages := g.Assert(); !valid {
		t.Fatalf("g.Assert() failed: %s", messages)
	}
}

/*
The second test ensures that our mock `Responder` does its job. This is accomplished
uing the httptest package from the standard library.
*/
func Test_MockResponder(t *testing.T) {

	// Wrap the mock function in the middleware, just like the app does
	handler := middleware(MockWriteResponse)

	// A bunch of strings with which to test the system. In a real application, these
	// would likely be randomly-generated.
	tests := []string{
		"foo",
		"basket",
		"ramakin",
		"logical fallacy",
		"existential terror",
		"€*¢",
	}

	// Test each string using httptest.Recorder. Again, there are many better ways to
	// do this; what follows is an aside for the example text suite.
	for _, input := range tests {
		req, err := http.NewRequest("GET", "http://localhost/"+input, nil)

		if err != nil {
			t.Fatalf("http.NewRequest failed: %s", err)
		}

		w := httptest.NewRecorder()
		handler(w, req)

		t.Logf("Testing %s...", input)

		if g, e := w.Body.String(), input; g != e {
			t.Fatalf("Got %s, expected %s")
		}
	}
}

/*
This is the main test: add all the dependencies to the global graph and run the
application. We can use the http package from the standard library to make sure
the responses from the server are correct, and shut down the server properly when
we're done.
*/
func Test_Application(t *testing.T) {

	app := Application{}
	exit := make(ExitChan)

	// Set up the graph, much the same was as in main.go – and exactly as in the test
	// above, only this time globally
	inj.Provide(
		&app,              // This time it is a pointer
		exit,              // `ExitChan`, this time shared
		&MockConfig{},     // Use the mock config
		MockLogger,        // `Logger`, as above
		MockWriteResponse, // Our specialised response writer
		t,                 // Pointer to `testing.T`, which is now going to be used
	)

	// Make sure everything was loaded properly again
	if valid, messages := inj.Assert(); !valid {
		t.Fatalf("inj.Assert() failed: %s", messages)
	}

	// Run the application in its own goroutine. It should stay running until it's explicitly
	// shut down
	go app.run()

	// Here we make some requests to the application, and check the responses. Again, there are
	// better ways to choose test strings
	tests := []string{
		"fulfilment",
		"doubt",
		"futility",
		"inetivatable narrative hijacking",
	}

	for _, input := range tests {

		res, err := http.Get(fmt.Sprintf("http://localhost:%d/%s", MOCK_HTTP_PORT, input))

		if err != nil {
			t.Fatalf("http.Get(): %s", err)
		}

		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()

		if err != nil {
			log.Fatal(err)
			t.Fatalf("ioutil.ReadAll %s", err)
		}

		t.Logf("Testing %s...", input)

		if g, e := res.StatusCode, 200; g != e {
			t.Errorf("Expected status code %d, got %d", g, e)
		}

		if g, e := string(body), input; g != e {
			t.Fatalf("Got %s, expected %s")
		}
	}

	// Kill the app
	exit <- struct{}{}
}

/*
That's pretty much it for the example application test suite. To recap, the suite:

- defines a mock `Configurer`, `Logger` and `Responder`
- tests an isolated graph with all dependencies
- explicitly tests the mock `Responder`
- runs the application (with all the mock dependencies in the global graph)
- checks that the real HTTP responses from the application are as expected.

Once again, this is an example. One would expect any real world application to be far more
complex, and to use a more sophisticated structure. In any case, the `inj` procedure is
always pretty much the same: your application pulls its dependencies from the graph, either
directly (using `inj.Provide`) or through a functional callback (ie. `inj.Inject`) – or both.
In the production code, the dependencies are assembled somewhere near the `main()` function; in
the tests, they can be provided specially, thus allowing discreet segments of the code to be
tested in the most realistic way possible.
*/
