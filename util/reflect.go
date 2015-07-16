package util

import "reflect"

// InterfaceIsZero determines whether or not an interface is its type's zero value.
func InterfaceIsZero(targetInterface interface{}) bool {
	return reflect.ValueOf(targetInterface) == reflect.Zero(reflect.TypeOf(targetInterface))
}

// ValueHasField returns whether a value has a given field.
func ValueHasField(value reflect.Value, fieldName string) bool {
	return value.FieldByName(fieldName).Kind() != reflect.Invalid
}
