package utils

import "reflect"

func SliceContains(sliceType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(sliceType)

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
