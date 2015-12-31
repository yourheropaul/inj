package inj

// A datasource reader sends an object to a datasource, identified by a
// given string key.
type DatasourceWriter interface {
	Write(string, interface{}) error
}
