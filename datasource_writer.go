package inj

// A datasource reader sends an object to a datasource, identified by a
// given string key. Refer to the documentation for the Datasource interface for more information.
type DatasourceWriter interface {
	Write(string, interface{}) error
}
