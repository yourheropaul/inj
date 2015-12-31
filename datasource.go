package inj

// A Datasource is any external interface that provides an interface given a
// string key. It's split into two fundamental components: the DatasourceReader
// and the DatasourceWriter. For the purposes of automatic configuration via
// dependency injection, a DatasourceReader is all that's required.
//
type Datasource interface {
	DatasourceReader
	DatasourceWriter
}
