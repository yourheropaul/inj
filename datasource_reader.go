package inj

// A datasource reader provides an object from a datasource, identified by a
// given string key. Refer to the documentation for the Datasource interface for more information.
type DatasourceReader interface {
	Read(string) (interface{}, error)
}
