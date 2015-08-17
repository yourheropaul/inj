package inj

import (
	"reflect"
	"strings"
)

type GraphNodeDependency struct {
	Name string
	Path StructPath
	Type reflect.Type
}

func findDependencies(t reflect.Type, deps *[]GraphNodeDependency, path *StructPath) error {

	for i := 0; i < t.NumField(); i++ {

		f := t.Field(i)

		// Ignore unpexported fields, regardless of tags
		if f.Anonymous {
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
		dep := parseStructTag(tag, identifier(f.Type))

		// Add the path in the struct
		dep.Path = branch

		// We also know the type
		dep.Type = f.Type

		// Add the dependency
		*deps = append(*deps, dep)
	}

	return nil
}

func parseStructTag(t reflect.StructTag, defaultName string) (d GraphNodeDependency) {

	d.Name = defaultName

	parts := strings.Split(t.Get("inj"), ",")

	if len(parts) > 0 && len(parts[0]) > 0 {
		d.Name = parts[0]
	}

	return
}
