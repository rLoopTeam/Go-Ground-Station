package dataStore

import (
	"Go-Ground-Station/gstypes"
	"strings"
	"fmt"
	"time"
	"sync"
	"Go-Ground-Station/gsgrpc"
	"runtime"
)

type DataStoreManager struct {
	packetChannel <-chan gstypes.PacketStoreElement
	receiversChannelHolder gsgrpc.ChannelsHolder
	rtDataStore RealTimeDataStore
	ticker *time.Ticker
	receiversCoordinator gstypes.ReceiversCoordination
	packetStoreCount int64

	rtDataStoreMutex *sync.Mutex
	rtData map[string]gstypes.RealTimeStreamElement
}

func (manager *DataStoreManager) Start (){
	go manager.run()
	//go manager.checker()
}

func (manager *DataStoreManager) run (){
	for {
		element := <- manager.packetChannel
		manager.StorePacket(element)
		//check if grpc wants to push a subscriber to the map
		select {
			case <- manager.receiversCoordinator.Call:
				manager.receiversCoordinator.Ack <- true
				<- manager.receiversCoordinator.Done
			default:
		}
		//this call is necessary so that the goroutine doesn't use too many cpu time at once
		runtime.Gosched()
	}
}

func (manager *DataStoreManager) checker (){
	//TODO: refactor function to databundle
	//check all RxTimes on data and set to 0 when RX greater than 4 seconds
	for t := range manager.ticker.C {
		manager.rtDataStoreMutex.Lock()
		currentTime := t.Unix()
		fmt.Println("checking...\n")
		for _, rtStreamElement := range manager.rtData{
			recordedTime := rtStreamElement.Data.RxTime
			if (currentTime - recordedTime) > 4{
				rtStreamElement.Data.Int64Value = 0
				rtStreamElement.Data.Uint64Value = 0
				rtStreamElement.Data.Float64Value = 0
				rtStreamElement.Data.RxTime = time.Now().Unix()
				manager.setParameter(rtStreamElement)
			}
		}
		manager.rtDataStoreMutex.Unlock()
	}
}

func (manager *DataStoreManager) StorePacket(packet gstypes.PacketStoreElement){
	rxTime := packet.RxTime
	packetName := packet.PacketName
	parameters := packet.Parameters
	paramCount := len(parameters)
	prefix := packet.ParameterPrefix
	dataBundle := gstypes.RealTimeDataBundle{}
	dataBundle.Data = make([]gstypes.RealTimeStreamElement,paramCount)

	for idx := 0; idx < paramCount; idx++ {
		parameter := parameters[idx]
		fullyFormattedName := cleanJoin(prefix,parameter.ParameterName)
		parameter.ParameterName = fullyFormattedName
		elem := manager.AddValue(packetName,rxTime,parameter)
		dataBundle.Data[idx] = elem
	}
	manager.sendDataBundle(dataBundle)
	manager.packetStoreCount++
	fmt.Printf("stored packet count: %d\n",manager.packetStoreCount)
}

func (manager *DataStoreManager) AddValue(packetName string,rxTime int64, staticElement gstypes.DataStoreElement) gstypes.RealTimeStreamElement{
	manager.rtDataStoreMutex.Lock()
	unit := gstypes.RealTimeDataStoreUnit{}
	unit.RxTime = rxTime
	switch staticElement.Data.ValueIndex{
		case 1:
			unit.Int64Value = int64(staticElement.Data.Int8Value)
			unit.ValueIndex = 1
		case 2:
			unit.Int64Value = int64(staticElement.Data.Int16Value)
			unit.ValueIndex = 1
		case 3:
			unit.Int64Value = int64(staticElement.Data.Int32Value)
			unit.ValueIndex = 1
		case 4:
			unit.Int64Value = staticElement.Data.Int64Value
			unit.ValueIndex = 1
		case 5:
			unit.Uint64Value = uint64(staticElement.Data.Uint8Value)
			unit.ValueIndex = 2
		case 6:
			unit.Uint64Value = uint64(staticElement.Data.Uint16Value)
			unit.ValueIndex = 2
		case 7:
			unit.Uint64Value = uint64(staticElement.Data.Uint32Value)
			unit.ValueIndex = 2
		case 8:
			unit.Uint64Value = staticElement.Data.Uint64Value
			unit.ValueIndex = 2
		case 9:
			unit.Float64Value = float64(staticElement.Data.FloatValue)
			unit.ValueIndex = 3
		case 10:
			unit.Float64Value = staticElement.Data.Float64Value
			unit.ValueIndex = 3
	}
	unit.Units = staticElement.Units

	rtDatastoreElement := gstypes.RealTimeStreamElement{}
	rtDatastoreElement.ParameterName = staticElement.ParameterName
	rtDatastoreElement.PacketName = packetName
	rtDatastoreElement.Data = unit

	manager.setParameter(rtDatastoreElement)
	manager.rtDataStoreMutex.Unlock()
	return rtDatastoreElement
}

func (manager *DataStoreManager) setParameter(element gstypes.RealTimeStreamElement){
	//update the map
	manager.rtData[element.ParameterName] = element
}

func (manager *DataStoreManager) sendDataBundle(dataBundle gstypes.RealTimeDataBundle){
	for channel := range manager.receiversChannelHolder.Receivers{
		select {
		case *channel <- dataBundle:
		default: fmt.Printf("streamerchannel is full \n")
		}
	}
}

func cleanJoin(prefix string, name string) string{
	var fullyFormattedName string

	prefix = strings.TrimSpace(prefix)
	name = strings.TrimSpace(name)

	if prefix != "" {
		s := []string{prefix,name}
		fullyFormattedName = strings.Join(s," ")
	}else {
		fullyFormattedName = name
	}
	return fullyFormattedName
}

func joiner(prefix string, name string) string{
	var fullyFormattedName string
	var formattedPrefix string
	var formattedName string

	prefix = strings.TrimSpace(prefix)
	name = strings.TrimSpace(name)

	if len(prefix) > 0 {
		explodedPrefix := strings.Split(prefix," ")
		formattedPrefix = formatter(explodedPrefix)
	}

	if len(name) > 0 {
		explodedName := strings.Split(name, " ")
		formattedName = formatter(explodedName)
	}

	s := []string{strings.TrimSpace(formattedPrefix),strings.TrimSpace(formattedName)}
	fullyFormattedName = strings.Join(s,"")
	return fullyFormattedName
}

func formatter(explodedString []string) string {
	formattedString := ""
	arrayLength := len(explodedString)
	for i := 0; i < arrayLength; i++ {
		str := explodedString[i]
		firstLetter := string(str[0])
		capitalized := strings.ToUpper(firstLetter)
		arr := []string{capitalized,string(str[1:])}
		formatted := strings.Join(arr,"")
		formattedString += formatted
	}
	return formattedString
}

func New (channelsHolder gsgrpc.ChannelsHolder, packetStoreChannel <- chan gstypes.PacketStoreElement, receiversCoordinator gstypes.ReceiversCoordination) *DataStoreManager{
	storeManager := &DataStoreManager{
		packetStoreCount:0,
		receiversChannelHolder: channelsHolder,
		packetChannel:packetStoreChannel,
		rtData: map[string]gstypes.RealTimeStreamElement{},
		rtDataStoreMutex: &sync.Mutex{},
		ticker: time.NewTicker(time.Second*2),
		receiversCoordinator:receiversCoordinator}
	return storeManager
}