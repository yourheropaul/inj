Are you troubled by dependency configuration issues in the middle of the night? Do you experience feelings of dread at the code boundaries of your application or during functional testing? Have you or your team ever used a global variable, exported package-level object or a build constraint hack? If the answer is *yes*, then don't wait another minute. Pick up your terminal and `go get github.com/yourheropaul/inj`.

Remember: inject yourself before you wreck yourself.

### What *is* this thing?

`inj` provides reflection-based dependency injection for Go structs and functions. Some parts of it will be familiar to anyone who's ever used [facebookgo/inject](https://github.com/facebookgo/inject); others bear a passing similarity to [depdendency injection in Angular.js](https://docs.angularjs.org/guide/di).  It's really designed for medium to large applications, and it's especially useful if your project is is BDD/TDD-orientated.

### Depdendency injection is boring and hard to maintain. I want something that works out of the box.

`inj` may well be for you. A good, idiomatic Go package is a simple, clearly laid out, well-tested Go package – and `inj` is, or aspires to be, all of those things. It works well and it works intuitively. There's only one, optional  configuration option (whether or not to use a global variable inside the package) and there's very little to think about once you're set up.

### Enough talk, let's see some code.

There are two ways to use `inj`. The first is struct-oriented, and not dissimilar to [facebookgo/inject](https://github.com/facebookgo/inject). In essence, you will create a 'graph' of objects which will automatically be wired together behind the scenes. usually the graph would be created at runtime, and you wouldn't change it during the execution of the application.

Suppose you're building an application reads config from a file, connects to a database of some kind, and the does some business logic until it's told to stop. You might model those behaviours with a struct that describes your application, like this:


```
type ExitChan chan interface{}

type Application struct {
    Database DatabaseInterface
    Config   ConfigInterface
    Exit     ExitChan
}
```

(Here, `DatabaseInterface` and `ConfigInterface` are implemented by various structs defined elsewhere, and their concrete types may vary dramatically between instances of the application. The `ExitChan` type is a channel on which the application will wait for data, and exit when it recieves any.)

Let's say you're creating an instance of `Application` in your `main()` function. How do you link the components together?

One way of solving the problem is assigning concrete objects directly to the `Application` object. For example:

```
func main() {
    app := Application{
        Database: DatabaseStruct{},
        Database: ConfigStruct{},
        Database: SomeChannel,
    }
}
```

That's fine; it'll work perfectly. But how about if another struct type wants to share the config? 


