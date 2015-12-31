package inj

// A datasource reader provides an object from a datasource, identified by a
// given string key.
type DatasourceReader interface {
	Read(string) (interface{}, error)
}
