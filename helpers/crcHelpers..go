package helpers

import (
	"encoding/binary"
	"rloop/Go-Ground-Station-1/constants"
)

func IsCrc16Check(dataCrc []byte, data []byte, dataLength uint32) (bool, error) {
	var result bool
	var err error
	var u16CRCVal uint16
	var u16CRCCalculatedVal uint16

	u16CRCVal = binary.LittleEndian.Uint16(dataCrc)
	u16CRCCalculatedVal = Crc16Val(data, dataLength)

	result = u16CRCVal == u16CRCCalculatedVal

	return result, err
}

func Crc16Val(data []byte, u32Length uint32) uint16 {
	var u16CRC uint16
	var u32Counter uint32
	u16CRC = 0

	for u32Counter = 0; u32Counter < u32Length; u32Counter++ {
		idx := ((u16CRC >> 8) ^ uint16(data[u32Counter])) & 0xFF
		u16CRC = constants.U16SWCRC_CRC_TABLE[idx] ^ (u16CRC << 8)
	}
	return u16CRC
}

func Crc16Bytes(data []byte, u32Length uint32) ([]byte, error) {
	return ParseValueToBytes(Crc16Val(data, u32Length))
}
