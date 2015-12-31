package inj

import "fmt"

// Add a Datasource, DatasourceReader or DatasourceWriter to
// the graph. Returns an error if the supplied argument isn't
// one of the accepted types.
func (g *Graph) AddDatasource(d interface{}) error {

	found := false

	if v, ok := d.(DatasourceReader); ok {
		found = true
		g.datasourceReaders = append(g.datasourceReaders, v)
	}

	if v, ok := d.(DatasourceWriter); ok {
		found = true
		g.datasourceWriters = append(g.datasourceWriters, v)
	}

	if !found {
		return fmt.Errorf("Supplied argument isn't a DatasourceReader or a DatasourceWriter")
	}

	return nil
}
