package datastore

import (
	"rloop/Go-Ground-Station-1/constants"
	"fmt"
	"rloop/Go-Ground-Station-1/gsgrpc"
	"rloop/Go-Ground-Station-1/gstypes"
	"runtime"
	"strings"
	"sync"
	"time"
)

type DataStoreManager struct {
	isRunningMutex         sync.RWMutex
	isRunning              bool
	doRunMutex             sync.RWMutex
	doRun                  bool
	signalChannel          chan bool
	checkerSignalChannel   chan bool
	packetChannel          <-chan gstypes.PacketStoreElement
	receiversChannelHolder *gsgrpc.ChannelsHolder
	ticker                 *time.Ticker
	packetStoreCount       int64

	rtDataStoreMutex *sync.Mutex
	rtData           map[string]gstypes.DataStoreElement
}

func (manager *DataStoreManager) Start() {
	manager.doRun = true
	if !manager.isRunning {
		fmt.Println("go run manager run")
		go manager.run()
		fmt.Println("go run checker")
		go manager.checker()
	}
}

func (manager *DataStoreManager) Stop() {
	manager.doRun = false
	if manager.isRunning {
		manager.signalChannel <- true
		manager.checkerSignalChannel <- true
	}
}

func (manager *DataStoreManager) run() {
	manager.isRunning = true
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
	manager.isRunning = false
}

func (manager *DataStoreManager) initDataStore() {
	var dataBundle gstypes.DataStoreBundle
	var arr []gstypes.DataStoreElement
	var preCount = 0
	var count = 0

	for _, definition := range constants.PacketDefinitions {
		preCount += len(definition.Parameters)
	}
	preCount = preCount * 2
	arr = make([]gstypes.DataStoreElement, preCount)

	for _, definition := range constants.PacketDefinitions {
		for _, node := range definition.MetaData {
			for _, param := range definition.Parameters {
				element := gstypes.DataStoreElement{}
				element.Data = gstypes.DataStoreUnit{}
				element.Data.ValueIndex = 4
				element.PacketName = node.Name
				element.FullParameterName = cleanJoin(node.ParameterPrefix, param.Name)
				element.IsStale = true
				element.ParameterName = param.Name
				element.Units = param.Units
				arr[count] = element
				count++
			}
		}
	}
	dataBundle = gstypes.DataStoreBundle{}
	dataBundle.Data = arr[:count]
	manager.rtDataStoreMutex.Lock()
	manager.saveToDataStore(dataBundle)
	manager.rtDataStoreMutex.Unlock()
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
	dataBundle := gstypes.DataStoreBundle{}
	dataBundle.Data = make([]gstypes.DataStoreElement, paramCount)

	for idx := 0; idx < paramCount; idx++ {
		parameter := parameters[idx]
		fullyFormattedName := cleanJoin(prefix, parameter.ParameterName)
		parameter.ParameterName = fullyFormattedName
		parameter.FullParameterName = fullyFormattedName
		parameter.RxTime = rxTime
		parameter.PacketName = packetName
		manager.UpdateDatastoreElement(&parameter)
		dataBundle.Data[idx] = parameter
	}
	manager.rtDataStoreMutex.Lock()
	manager.saveToDataStore(dataBundle)
	manager.sendDatastoreUpdate()
	manager.rtDataStoreMutex.Unlock()
	manager.packetStoreCount++
}

func (manager *DataStoreManager) UpdateDatastoreElement(element *gstypes.DataStoreElement) {
	switch element.Data.ValueIndex {
	case 1:
		element.Data.Int64Value = int64(element.Data.Int8Value)
		element.Data.ValueIndex = 4
	case 2:
		element.Data.Int64Value = int64(element.Data.Int16Value)
		element.Data.ValueIndex = 4
	case 3:
		element.Data.Int64Value = int64(element.Data.Int32Value)
		element.Data.ValueIndex = 4
	case 5:
		element.Data.Uint64Value = uint64(element.Data.Uint8Value)
		element.Data.ValueIndex = 8
	case 6:
		element.Data.Uint64Value = uint64(element.Data.Uint16Value)
		element.Data.ValueIndex = 8
	case 7:
		element.Data.Uint64Value = uint64(element.Data.Uint32Value)
		element.Data.ValueIndex = 8
	case 9:
		element.Data.Float64Value = float64(element.Data.FloatValue)
		element.Data.ValueIndex = 10
	}
}

func (manager *DataStoreManager) checkDatastore(currTime time.Time) {
	var data []gstypes.DataStoreElement
	//will count the amount of zeroed parameters and be used to slice the array the right size
	var count = 0
	manager.rtDataStoreMutex.Lock()
	//the current length or amount of parameters in the datastore
	paramLen := len(manager.rtData)
	//make a new array that will be populated with the new values, enough to fit all current parameters
	data = make([]gstypes.DataStoreElement, paramLen)
	//used to calculate the time difference and to set the new time
	//of when the parameters were updated last, only for parameters that will be zeroed
	currentTime := currTime.Unix()
	fmt.Println("checking...")
	for _, dataStoreElement := range manager.rtData {
		recordedTime := dataStoreElement.RxTime
		if (currentTime - recordedTime) > 4 {
			dataStoreElement.IsStale = true
			dataStoreElement.Data.Int64Value = 0
			dataStoreElement.Data.Uint64Value = 0
			dataStoreElement.Data.Float64Value = 0
			dataStoreElement.RxTime = time.Now().Unix()
			data[count] = dataStoreElement
			count++
		}
	}
	if count > 0 {
		dataBundle := gstypes.DataStoreBundle{}
		dataBundle.Data = data[0:count]
		manager.saveToDataStore(dataBundle)
		manager.sendDatastoreUpdate()
	}
	manager.rtDataStoreMutex.Unlock()
}

func (manager *DataStoreManager) saveToDataStore(dataBundle gstypes.DataStoreBundle) {
	for _, element := range dataBundle.Data {
		//fmt.Printf("storing data: %v \n", element)
		manager.rtData[element.FullParameterName] = element
	}
}

func (manager *DataStoreManager) sendDatastoreUpdate() {
	dataBundle := gstypes.DataStoreBundle{}
	dataBundle.Data = make([]gstypes.DataStoreElement, len(manager.rtData))

	idx := 0
	for _, value := range manager.rtData {
		dataBundle.Data[idx] = value
		idx++
	}
	manager.sendDataBundle(dataBundle)
}

func (manager *DataStoreManager) sendDataBundle(dataBundle gstypes.DataStoreBundle) {
	manager.receiversChannelHolder.ReceiverMutex.Lock()
	//send the bundle to all subscribers
	for channel := range manager.receiversChannelHolder.Receivers {
		select {
		case *channel <- dataBundle:
		default:
			fmt.Printf("streamerchannel is full \n")
		}
	}
	manager.receiversChannelHolder.ReceiverMutex.Unlock()
}

func (manager *DataStoreManager) GetStatus() (bool, bool) {
	defer func() {
		manager.isRunningMutex.RUnlock()
		manager.doRunMutex.RUnlock()
	}()
	manager.isRunningMutex.RLock()
	manager.doRunMutex.RLock()
	return manager.isRunning, manager.doRun
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
	dataStoreManager := &DataStoreManager{
		signalChannel:          signalChannel,
		checkerSignalChannel:   checkerSignalChannel,
		packetStoreCount:       0,
		receiversChannelHolder: channelsHolder,
		packetChannel:          packetStoreChannel,
		rtData:                 map[string]gstypes.DataStoreElement{},
		rtDataStoreMutex:       &sync.Mutex{},
		ticker:                 time.NewTicker(time.Second * 3)}
	dataStoreManager.initDataStore()
	return dataStoreManager, packetStoreChannel
}
