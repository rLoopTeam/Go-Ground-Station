package gstypes

import "strconv"

func (data *DataStoreUnit) AsString() string{
	var str string
	switch data.ValueIndex{
	case 1:
		str = strconv.FormatInt(int64(data.Int8Value),10)
	case 2:
		str = strconv.FormatInt(int64(data.Int16Value),10)
	case 3:
		str = strconv.FormatInt(int64(data.Int32Value),10)
	case 4:
		str = strconv.FormatInt(int64(data.Int64Value),10)
	case 5:
		str = strconv.FormatUint(uint64(data.Uint8Value),10)
	case 6:
		str = strconv.FormatUint(uint64(data.Uint16Value),10)
	case 7:
		str = strconv.FormatUint(uint64(data.Uint32Value),10)
	case 8:
		str = strconv.FormatUint(data.Uint64Value,10)
	case 9:
		str = strconv.FormatFloat(float64(data.FloatValue),'f',10,32)
	case 10:
		str = strconv.FormatFloat(data.Float64Value,'f',10,64)
	}
	return str
}
