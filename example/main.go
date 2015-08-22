/*
Package inj/example is a demonstration the inj package. It takes the form of a very simple web application
that is modelled by a struct. Depdendencies for the application are defined on struct fields using tags,
and provided in the main function. On startup, the application will sort out and check its dependencies,
then start a simple HTTP server with two endpoints.

The functions for the endpoints use a rudimentary HTTP middleware wrapped so they can make use of inj.Inject,
and thus have no-standard argument requirements.

Hopefully the code is reasonably self-explanatory.
*/
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/yourheropaul/inj"
)

type Application struct {
	Config Configurer `inj:""`
	Log    Logger     `inj:""`
	Exit   ExitChan   `inj:""`
}

type Configurer interface {
	Port() int
}

type Config struct{}

func (c *Config) Port() int { return 8080 }

type Logger func(string, ...interface{}) (int, error)

type ExitChan chan interface{}

type Responder func(w http.ResponseWriter, r *http.Request)

func WriteResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {

	app := Application{}

	inj.Provide(
		&app,
		make(ExitChan),
		&Config{},
		fmt.Printf,
		WriteResponse,
	)

	if valid, messages := inj.Assert(); !valid {
		fmt.Println(messages)
		os.Exit(1)
	}

	app.run()
}

func (a Application) run() {

	http.HandleFunc("/", middleware(a.handler))
	http.HandleFunc("/shutdown", middleware(a.shutdown))

	a.Log("Running server on port %d...\n", a.Config.Port())
	go http.ListenAndServe(fmt.Sprintf(":%d", a.Config.Port()), nil)

	<-a.Exit

	a.Log("Received data on exit channel.\n")
}

func (a Application) handler(w http.ResponseWriter, r *http.Request, responder Responder) {
	responder(w, r)
}

func (a Application) shutdown(w http.ResponseWriter, e ExitChan) {
	w.Write([]byte("Shutting down!"))
	a.Log("Shutting down...\n")
	e <- 0
}

func middleware(fn interface{}) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		inj.Inject(fn, w, r)
	}
}
