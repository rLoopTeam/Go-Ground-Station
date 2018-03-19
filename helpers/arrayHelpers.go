package helpers

import (
	"Go-Ground-Station/gstypes"
)

func Push (array *[]gstypes.DataStoreElement, idx int, element gstypes.DataStoreElement){
	arrSize := cap(*array)
	if idx >= arrSize {
		newArr := make([]gstypes.DataStoreElement,arrSize)
		*array = append(*array, newArr...)
	}
	(*array)[idx] = element
}
