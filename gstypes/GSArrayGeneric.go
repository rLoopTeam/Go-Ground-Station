package gstypes

type GSArrayGeneric struct {
	Count int
	Capacity int
	Data []interface{}
}

func (arr *GSArrayGeneric) Push (element interface{})  {
	if arr.Count >= arr.Capacity {
		newArr := make([]interface{},arr.Capacity*2)
		copy(newArr,arr.Data)
		arr.Data = newArr
		arr.Capacity *= 2
	}
	arr.Data[arr.Count] = element
	arr.Count++
}

func (arr *GSArrayGeneric) Empty() bool{
	return arr.Count < 1
}

func (arr *GSArrayGeneric) Length() int{
	return arr.Count
}

func (arr *GSArrayGeneric) Get(idx int) interface{}{
	var element interface{}
	if idx < arr.Count{
		element = arr.Data[idx]
	}
	return element
}

func (arr *GSArrayGeneric) CurrentIndex() int{
	return arr.Count-1
}

func (arr *GSArrayGeneric) Remove(idx int) bool{
	if idx < arr.Count{
		arr.Count--
		copy(arr.Data[idx:], arr.Data[idx+1:])
		arr.Data[arr.Count] = nil
		//arr.Data = arr.Data[:arr.Count]
		//arr.Data = append(arr.Data[:idx], arr.Data[idx+1:]...)
		return true
	}else{
		return false
	}
}

func (arr *GSArrayGeneric) GetSlice() []interface{}{
	return arr.Data[:arr.Count]
}

func NewGSArray() GSArrayGeneric {
	array := GSArrayGeneric{
		Count: 0,
		Capacity: 1,
		Data: make([]interface{},1)}
	return array
}