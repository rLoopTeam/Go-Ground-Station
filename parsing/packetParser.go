package parsing

import (
	"rloop/Go-Ground-Station/constants"
	"rloop/Go-Ground-Station/gstypes"
	"rloop/Go-Ground-Station/helpers"
	"encoding/binary"
	"time"
	"errors"
)

func ParsePacket(nodePort int, nodeName string, packet []byte, errcount *int) (gstypes.PacketStoreElement,error) {
	//var packetStoreElement gstypes.PacketStoreElement
	//var err error
	packetLength := len(packet)
	//get the index of the last byte of the payload
	lastPayloadByteIndex := packetLength - 2
	//get the bytes defining the packet type
	packetTypeBytes := packet[4:6]
	//get the bytes defining the payload length
	payloadLengthBytes := packet[6:8]
	//extract the payload bytes
	payload := packet[8:lastPayloadByteIndex]
	//convert the bytes to int
	packetType := binary.LittleEndian.Uint16(packetTypeBytes)
	//convert the bytes into int
	payloadLength := binary.LittleEndian.Uint16(payloadLengthBytes)
	payloadLengthInt := int(payloadLength)
	//retrieve the definition of the packet
	definition := constants.PacketDefinitions[packetType]
	//make crc check and if correct proceed to parsing
	//TODO: fix crc
	/*
	crcCheck, err := helpers.IsCrc16Check(payloadLengthInt,payload,packet[lastPayloadByteIndex:])
	if crcCheck {
		packetStoreElement, err = ParsePayload(definition, payloadLengthInt, nodePort, payload)
	}else if err == nil{
		err = errors.New("CRC check failed")
	}
	*/
	return ParsePayload(definition, payloadLengthInt, nodePort, nodeName, payload)
}

func ParsePayload(definition gstypes.PacketDefinition, payloadLength int, port int, nodeName string, payload []byte) (gstypes.PacketStoreElement, error) {
	packetStoreElement := gstypes.PacketStoreElement{}
	var parseError error = nil
	//Declare variable that will be forwarded to the consumers
	var dataStoreElement gstypes.DataStoreElement
	var currentParameter gstypes.Param
	//Retrieve the metadata for the particular host, in some cases the same packet types are used by different hosts
	nodeMetaData := definition.MetaData[nodeName]
	//retrieve all the parameters (array of parameter objects)
	parameters := definition.Parameters
	isInParameterLoop := false
	//used in packets where a loop is present
	parameterLoopOffset := 0
	currentParameterIndex := 0
	//countParamaters := len(parameters)
	//index used in datastore array, where processed parameters are pushed
	dataStoreElementIdx := 0
	//datastoreElementArray := make([]gstypes.DataStoreElement,countParamaters)

	packetStoreElement.Port = port
	packetStoreElement.PacketType = definition.PacketType
	packetStoreElement.PacketName = nodeMetaData.Name
	packetStoreElement.ParameterPrefix = nodeMetaData.ParameterPrefix
	packetStoreElement.RxTime = time.Now().Unix()

	TestArr := gstypes.NewGSarrayDSE()

	//count he number of params in the loop
	for _, param := range parameters {
		if param.BeginLoop{
			isInParameterLoop = true
		}

		if isInParameterLoop{
			//increment the count of parameters in the loop
			parameterLoopOffset++
		}

		if param.EndLoop {
			isInParameterLoop = false
		}
	}

	for currentPayloadByte := 0; currentPayloadByte < payloadLength; currentPayloadByte += currentParameter.Size{
		currentParameter = parameters[currentParameterIndex]
		lastPayloadByte := currentPayloadByte + currentParameter.Size

		if currentParameter.BeginLoop{
			isInParameterLoop = true
		}
		//slice the bytes corresponding to the current parameter
		payloadSlice := payload[currentPayloadByte:lastPayloadByte]
		//parse the bytes to their respective datatypes
		dataUnit, err := helpers.ParseByteToValue(currentParameter.Type,payloadSlice)

		if err == nil {
			dataStoreElement.ParameterName = currentParameter.Name
			dataStoreElement.Data = dataUnit
			dataStoreElement.Units = currentParameter.Units
			//fmt.Printf("storing parameter in array with index: %d\n", dataStoreElementIdx)
			//helpers.Push(&datastoreElementArray,dataStoreElementIdx,dataStoreElement)
			//fmt.Println("packet stored in array succesful")

			TestArr.Push(dataStoreElement)
		}else{
			errMessage := "Parser Error: "  + err.Error()
			parseError = errors.New(errMessage)
			break
		}
		//if the current parameter is the last parameter in the loop
		//we revert back to the first parameter in the loop
		//else we increment the index to the next parameter in the loop
		if isInParameterLoop && currentParameter.EndLoop{
			isInParameterLoop = false
			currentParameterIndex -= parameterLoopOffset
		}else{
			currentParameterIndex++
		}
		dataStoreElementIdx++
	}

	packetStoreElement.Parameters = TestArr.GetSlice()
	//packetStoreElement.Parameters = datastoreElementArray[:dataStoreElementIdx]
	return packetStoreElement, parseError
}


