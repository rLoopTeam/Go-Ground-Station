package helpers

import (
	"github.com/howeyc/crc16"
	"fmt"
	"go/types"
)

func isCrc16Check(payloadlength int, payload []byte, payloadCrc []byte) (bool,error){

	result := false
	var err error = nil
	var payloadCrcValue uint16
	payloadCrcStoreUnit,err := ParseByteToValue(types.Uint16,payloadCrc)
	payloadCrcValue = payloadCrcStoreUnit.Uint16Value
	controlCrc := crc16.ChecksumCCITTFalse(payload)

	result = controlCrc == payloadCrcValue

	fmt.Printf("Payload crc bytes: %v\n",payloadCrc)
	fmt.Printf("Payload crc: %v\n",payloadCrcValue)
	fmt.Printf("Control crc: %v\n",controlCrc)

	/*
	var crc uint16 = 0x00
	var j uint16
	var c uint8

	for i := 0; i < payloadlength; i++ {
		c = payload[i];
		if c > 255 {
			err = errors.New("out of range")
			break
			}

		j = ((crc >> 8) ^ c) & 0xFF
		crc = constants.U16SWCRC_CRC_TABLE[j] ^ (crc << 8)
	}

	if(err == nil){
		result = crc == crcIntValue
	}
	*/
	return result, err
}
