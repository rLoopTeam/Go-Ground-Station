package helpers

import (
	"encoding/binary"
	"rloop/Go-Ground-Station/gstypes"
	"go/types"
	"bytes"
)

func ParseByteToValue(valueType types.BasicKind,slice []byte)(gstypes.DataStoreUnit, error) {
	var err error = nil
	var value gstypes.DataStoreUnit
	var endianness binary.ByteOrder = binary.LittleEndian

	switch valueType {
	case types.Int8:
		var val int8
		buf := bytes.NewReader(slice)
		err = binary.Read(buf, endianness, &val)
		if err == nil{
			value.Int8Value = val
			value.ValueIndex = 1
		}
	case types.Int16:
		var val int16
		buf := bytes.NewReader(slice)
		err = binary.Read(buf, endianness, &val)
		if err == nil{
			value.Int16Value = val
			value.ValueIndex = 2
		}
	case types.Int32:
		var val int32
		buf := bytes.NewReader(slice)
		err = binary.Read(buf, endianness, &val)
		if err == nil{
			value.Int32Value = val
			value.ValueIndex = 3
		}
	case types.Int64:
		var val int64
		buf := bytes.NewReader(slice)
		err = binary.Read(buf, endianness, &val)
		if err == nil{
			value.Int64Value = val
			value.ValueIndex = 4
		}
	case types.Uint8:
		var val uint8
		buf := bytes.NewReader(slice)
		err = binary.Read(buf, endianness, &val)
		if err == nil{
			value.Uint8Value = val
			value.ValueIndex = 5
		}
	case types.Uint16:
		var val uint16
		buf := bytes.NewReader(slice)
		err = binary.Read(buf, endianness, &val)
		if err == nil{
			value.Uint16Value = val
			value.ValueIndex = 6
		}
	case types.Uint32:
		var val uint32
		buf := bytes.NewReader(slice)
		err = binary.Read(buf, endianness, &val)
		if err == nil{
			value.Uint32Value = val
			value.ValueIndex = 7
		}
	case types.Uint64:
		var val uint64
		buf := bytes.NewReader(slice)
		err = binary.Read(buf, endianness, &val)
		if err == nil{
			value.Uint64Value = val
			value.ValueIndex = 8
		}
	case types.Float32:
		var val float32
		buf := bytes.NewReader(slice)
		err = binary.Read(buf, endianness, &val)
		if err == nil{
			value.FloatValue = val
			value.ValueIndex = 9
		}
	case types.Float64:
		var val float64
		buf := bytes.NewReader(slice)
		err = binary.Read(buf, endianness, &val)
		if err == nil{
			value.Float64Value = val
			value.ValueIndex = 10
		}
	}
	return value,err
}

func ParseValueToBytes(value interface{}) ([]byte, error){
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, value)
	return buf.Bytes(),err
}