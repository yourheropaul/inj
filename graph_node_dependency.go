package inj

import (
	"reflect"
	"strings"
)

type graphNodeDependency struct {
	DatasourcePaths []string
	Path            structPath
	Type            reflect.Type
}

func findDependencies(t reflect.Type, deps *[]graphNodeDependency, path *structPath) error {

	for i := 0; i < t.NumField(); i++ {

		f := t.Field(i)

		// Ignore unpexported fields, regardless of tags
		if f.PkgPath != "" {
			continue
		}

		// Get all tags
		tag := f.Tag

		// Generate a struct path branch
		branch := path.Branch(f.Name)

		// Ignore tags that don't have injection deps
		if !strings.Contains(string(tag), "inj:") {

			if f.Type.Kind() == reflect.Struct {

				// Recurse
				findDependencies(f.Type, deps, &branch)
			}

			continue
		}

		// Assemble everything we know about the dependency
		dep := parseStructTag(tag)

		// Add the path in the struct
		dep.Path = branch

		// We also know the type
		dep.Type = f.Type

		// Add the dependency
		*deps = append(*deps, dep)
	}

	return nil
}

func parseStructTag(t reflect.StructTag) (d graphNodeDependency) {

	parts := strings.Split(t.Get("inj"), ",")

	if len(parts) > 0 && len(parts[0]) > 0 {
		d.DatasourcePaths = parts
	}

	return
}
