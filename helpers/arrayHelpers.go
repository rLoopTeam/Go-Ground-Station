package helpers

import (
	"rloop/Go-Ground-Station/gstypes"
)

func Push(array *[]gstypes.DataStoreElement, idx int, element gstypes.DataStoreElement) {
	arrSize := cap(*array)
	if idx >= arrSize {
		newArr := make([]gstypes.DataStoreElement, arrSize)
		*array = append(*array, newArr...)
	}
	(*array)[idx] = element
}

func AppendVariadic(arrays ...[]byte) []byte {
	newArr := make([]byte, 1)
	if arrays == nil {
		return newArr
	}

	for _, arg := range arrays {
		newArr = append(newArr, arg...)
	}

	return newArr[1:]
}
