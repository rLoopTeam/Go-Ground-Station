package datastore

import (
	"fmt"
	"rloop/Go-Ground-Station/gsgrpc"
	"rloop/Go-Ground-Station/gstypes"
	"runtime"
	"strings"
	"sync"
	"time"
)

type DataStoreManager struct {
	IsRunning              bool
	doRun                  bool
	signalChannel          chan bool
	checkerSignalChannel   chan bool
	packetChannel          <-chan gstypes.PacketStoreElement
	receiversChannelHolder *gsgrpc.ChannelsHolder
	ticker                 *time.Ticker
	packetStoreCount       int64

	rtDataStoreMutex *sync.Mutex
	rtData           map[string]gstypes.RealTimeStreamElement
}

func (manager *DataStoreManager) Start() {
	manager.doRun = true
	if !manager.IsRunning {
		fmt.Println("go run manager run")
		go manager.run()
		fmt.Println("go run checker")
		go manager.checker()
	}
}

func (manager *DataStoreManager) Stop() {
	manager.doRun = false
	if manager.IsRunning {
		manager.signalChannel <- true
		manager.checkerSignalChannel <- true
	}
}

func (manager *DataStoreManager) run() {
	manager.IsRunning = true
MainLoop:
	for {
		select {
		case element := <-manager.packetChannel:
			manager.ProcessNewPacket(element)
		case <-manager.signalChannel:
			break MainLoop
		}
		//this call is necessary so that the goroutine doesn't use too many cpu time at once
		runtime.Gosched()
	}
	fmt.Println("manager isrunning = false")
	manager.IsRunning = false
}

func (manager *DataStoreManager) checker() {
	fmt.Println("Checker started")
	//check all RxTimes on data and set to 0 when RX greater than 4 seconds
CheckerLoop:
	for {
		select {
		case t := <-manager.ticker.C:
			manager.checkDatastore(t)
		case <-manager.checkerSignalChannel:
			break CheckerLoop
		}
	}
}

func (manager *DataStoreManager) ProcessNewPacket(packet gstypes.PacketStoreElement) {
	rxTime := packet.RxTime
	packetName := packet.PacketName
	parameters := packet.Parameters
	paramCount := len(parameters)
	prefix := packet.ParameterPrefix
	dataBundle := gstypes.RealTimeDataBundle{}
	dataBundle.Data = make([]gstypes.RealTimeStreamElement, paramCount)

	for idx := 0; idx < paramCount; idx++ {
		parameter := parameters[idx]
		fullyFormattedName := cleanJoin(prefix, parameter.ParameterName)
		parameter.ParameterName = fullyFormattedName
		elem := manager.MakeRTSElement(packetName, rxTime, parameter)
		dataBundle.Data[idx] = elem
	}
	manager.rtDataStoreMutex.Lock()
	manager.saveToDataStore(dataBundle)
	manager.rtDataStoreMutex.Unlock()
	manager.packetStoreCount++
	//fmt.Printf("stored packet count: %d\n", manager.packetStoreCount)
	//fmt.Printf("latest datastore state: \n %v \n", manager.rtData)
}

func (manager *DataStoreManager) MakeRTSElement(packetName string, rxTime int64, staticElement gstypes.DataStoreElement) gstypes.RealTimeStreamElement {
	unit := gstypes.RealTimeDataStoreUnit{}
	switch staticElement.Data.ValueIndex {
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
	unit.RxTime = rxTime
	unit.Units = staticElement.Units
	unit.IsStale = false

	rtDatastoreElement := gstypes.RealTimeStreamElement{}
	rtDatastoreElement.ParameterName = staticElement.ParameterName
	rtDatastoreElement.PacketName = packetName
	rtDatastoreElement.Data = unit

	return rtDatastoreElement
}

func (manager *DataStoreManager) checkDatastore(currTime time.Time) {
	manager.rtDataStoreMutex.Lock()
	//this variable will be used to slice the array the right size, with the number of zeroed parameters
	count := 0
	//the current length or amount of parameters in the datastore
	paramLen := len(manager.rtData)
	//make a new array that will be populated with the new values, enough to fit all current parameters
	data := make([]gstypes.RealTimeStreamElement, paramLen)
	//used to calculate the time difference and to set the new time
	//of when the parameters were updated last, only for parameters that will be zeroed
	currentTime := currTime.Unix()
	fmt.Println("checking...")
	for _, rtStreamElement := range manager.rtData {
		recordedTime := rtStreamElement.Data.RxTime
		if (currentTime - recordedTime) > 4 {
			rtStreamElement.Data.IsStale = true
			rtStreamElement.Data.Int64Value = 0
			rtStreamElement.Data.Uint64Value = 0
			rtStreamElement.Data.Float64Value = 0
			rtStreamElement.Data.RxTime = time.Now().Unix()
			data[count] = rtStreamElement
			count++
		}
	}
	if count > 0 {
		dataBundle := gstypes.RealTimeDataBundle{}
		dataBundle.Data = data[0:count]
		manager.saveToDataStore(dataBundle)
	}
	manager.rtDataStoreMutex.Unlock()
}

func (manager *DataStoreManager) saveToDataStore(dataBundle gstypes.RealTimeDataBundle) {
	for _, element := range dataBundle.Data {
		manager.rtData[element.ParameterName] = element
	}
	manager.sendDataBundle(dataBundle)
}

func (manager *DataStoreManager) sendDataBundle(dataBundle gstypes.RealTimeDataBundle) {
	//check if grpc wants to push a subscriber to the map
	select {
	case <-manager.receiversChannelHolder.Coordinator.Call:
		manager.receiversChannelHolder.Coordinator.Ack <- true
		<-manager.receiversChannelHolder.Coordinator.Done
	default:
	}
	//send the bundle to all subscribers
	for channel := range manager.receiversChannelHolder.Receivers {
		select {
		case *channel <- dataBundle:
		default:
			fmt.Printf("streamerchannel is full \n")
		}
	}
}

func cleanJoin(prefix string, name string) string {
	var fullyFormattedName string

	prefix = strings.TrimSpace(prefix)
	name = strings.TrimSpace(name)

	if prefix != "" {
		s := []string{prefix, name}
		fullyFormattedName = strings.Join(s, " ")
	} else {
		fullyFormattedName = name
	}
	return fullyFormattedName
}

func New(channelsHolder *gsgrpc.ChannelsHolder) (*DataStoreManager, chan<- gstypes.PacketStoreElement) {
	//the channel that will be used to transfer data between the parser and the datastoremanager
	packetStoreChannel := make(chan gstypes.PacketStoreElement, 64)
	signalChannel := make(chan bool)
	checkerSignalChannel := make(chan bool)
	storeManager := &DataStoreManager{
		signalChannel:          signalChannel,
		checkerSignalChannel:   checkerSignalChannel,
		packetStoreCount:       0,
		receiversChannelHolder: channelsHolder,
		packetChannel:          packetStoreChannel,
		rtData:                 map[string]gstypes.RealTimeStreamElement{},
		rtDataStoreMutex:       &sync.Mutex{},
		ticker:                 time.NewTicker(time.Second * 3)}
	return storeManager, packetStoreChannel
}
