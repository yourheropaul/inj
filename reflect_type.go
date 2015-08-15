package inj

import "reflect"

func getReflectionTypes(input interface{}) (master_type, specific_type reflect.Type) {

	// Grab the master type
	master_type = reflect.TypeOf(input)

	// We need the specific type
	if master_type.Kind() == reflect.Ptr {
		specific_type = master_type.Elem()
	} else {
		specific_type = master_type
	}

	return
}
