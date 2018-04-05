package parsing

import (
	"rloop/Go-Ground-Station/constants"
	"rloop/Go-Ground-Station/gstypes"
	"encoding/binary"
	"rloop/Go-Ground-Station/helpers"
	"fmt"
	"time"
	"errors"
)

func ParsePacket(nodePort int, packet []byte, packetStoreChannel chan <- gstypes.PacketStoreElement, errcount *int) {
	defer func() {
		if r := recover(); r != nil {
			*errcount++
			fmt.Println("Problem with parsing packet in: ", r)
			fmt.Printf("errcount on nodeport %d: %d\n", nodePort,*errcount)
		}
	}()
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
	packetStoreElement, err := parsePayload(definition, payloadLengthInt,nodePort, payload)
	if err == nil {
		packetStoreChannel <- packetStoreElement
	}
	//testCrc,_ := helpers.isCrcCheck(payloadLengthInt,payload,packet[lastPayloadByteIndex:])
	//fmt.Printf("CRC check: %t\n", testCrc)
	/*
	if result, err := helpers.isCrcCheck(payloadLengthInt,payload,packet[lastPayloadByteIndex:]); result{
		parsePayload(definition, payloadLengthInt, payload)
	}else if err != nil{
		//log the error or faulty crc
	}
	*/
}

func parsePayload(definition gstypes.PacketDefinition, payloadLength int, port int,payload []byte) (gstypes.PacketStoreElement, error) {
	packetStoreElement := gstypes.PacketStoreElement{}
	var parseError error = nil
	//Declare variable that will be forwarded to the consumers
	var dataStoreElement gstypes.DataStoreElement
	var currentParameter gstypes.Param
	//need the nodenames to distinguish the hosts with similar functions
	nodeName := constants.Hosts[port].Name
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


