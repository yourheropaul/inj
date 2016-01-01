package inj

// A Datasource is any external interface that provides an interface given a
// string key. It's split into two fundamental components: the DatasourceReader
// and the DatasourceWriter. For the purposes of automatic configuration via
// dependency injection, a DatasourceReader is all that's required.
//
// Datasource paths are supplied by values in structfield tags (which are empty
// for normal inj operation). Here's an example of a struct that will try to
// fulfil its dependencies from a datasource:
//
//  type MyStruct stuct {
//      SomeIntVal int `inj:"some.path.in.a.datasource"`
//  }
//
// When trying to connect the dependency on an instance of the struct above, inj will
// poll any available DatasourceReaders in the graph by calling their Read() function
// with the string argument "some.path.in.a.datasource". If the function doesn't return
// an error, then the resultant value will be used for the dependency injection.
//
// A struct tag may contain multiple datasource paths, separated by commas. The paths
// will be polled in order of their appearance in the code, and the value from the first
// DatasourceReader that doesn't return an error on its Read() function will be used to
// meet the dependency.
//
// DatasourceWriters function in a similar way: when a value is set on an instance of a
// struct via inj's dependency injection, any associated DatasourceWriters' Write() functions
// are called for each datasource path.
//
type Datasource interface {
	DatasourceReader
	DatasourceWriter
}
