package inj

import "fmt"

// Add any number of Datasources, DatasourceReaders or DatasourceWriters
// to the graph. Returns an error if any of the supplied arguments aren't
// one of the accepted types.
//
// Once added, the datasources will be active immediately, and the graph
// will automatically re-Provide itself, so that any depdendencies that
// can only be met by an external datasource will be wired up automatically.
//
func (g *graph) AddDatasource(ds ...interface{}) error {

	for i, d := range ds {
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
			return fmt.Errorf("Supplied argument %d isn't a DatasourceReader or a DatasourceWriter", i)
		}
	}

	g.Provide()

	return nil
}
