package gstypes

type GSArrayDSE struct{
	array GSArrayGeneric
}

func (arr *GSArrayDSE) Push (element DataStoreElement)  {
	arr.array.Push(element)
}

func (arr *GSArrayDSE) Empty() bool{
	return arr.array.Empty()
}

func (arr *GSArrayDSE) Length() int{
	return arr.array.Length()
}

func (arr *GSArrayDSE) Get(idx int) DataStoreElement{
	return arr.array.Get(idx).(DataStoreElement)
}

func (arr *GSArrayDSE) CurrentIndex() int{
	return arr.array.CurrentIndex()
}

func (arr *GSArrayDSE) Remove(idx int) bool{
	return arr.array.Remove(idx)
}

func (arr *GSArrayDSE) GetSlice() []DataStoreElement{
	slice := arr.array.GetSlice()
	dataStoreElementArray := make([]DataStoreElement,arr.array.Count)
	for idx, element := range slice{
		dataStoreElementArray[idx] = element.(DataStoreElement)
	}
	return dataStoreElementArray
}

func NewGSarrayDSE() GSArrayDSE{
	arr := GSArrayDSE{
		array:NewGSArray()}
	return arr
}