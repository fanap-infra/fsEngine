package utils

import (
	"reflect"
)

func HasArray(a []interface{}, e interface{}) bool {
	for _, b := range a {
		if b == e {
			return true
		}
	}
	return false
}

// ToDo: correct this method
func ItemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)
	if arr.Kind() != reflect.Slice {
		panic("Invalid data-type")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}
