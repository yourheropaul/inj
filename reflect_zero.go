package inj

import "reflect"

func zero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	}
}
