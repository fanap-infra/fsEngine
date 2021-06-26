package utils

import (
	"reflect"

	"github.com/fanap-infra/log"
)

func HasArray(a []interface{}, e interface{}) bool {
	for _, b := range a {
		if b == e {
			return true
		}
	}
	return false
}

func ItemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)
	if arr.Kind() != reflect.Slice {
		log.Warnv("Invalid data-type array", "kind", arr.Kind())
	}
	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}
