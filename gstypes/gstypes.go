package gstypes

import (
	"go/types"
	"sync"
)

type Node struct {
	IP   string
	Port int
	Name string
}

type Param struct {
	Name      string
	Type      types.BasicKind
	Units     string
	Size      int
	BeginLoop bool
	EndLoop   bool
}

type NodeInfo struct {
	Name            string
	ParameterPrefix string
	Node            string
}

type PacketDefinition struct {
	MetaData   map[string]NodeInfo
	PacketType int
	DAQ        bool
	Parameters []Param
}

type PacketDAQ struct {
	Name       string
	PacketType int
	Node       string
	DAQ        bool
	DataType   types.BasicKind
	DataSize   int
}

type PacketStoreElement struct {
	PacketName string
	PacketType int
	ParameterPrefix string
	Port int
	RxTime int64
	Parameters []DataStoreElement
}

type DataStoreElement struct {
	ParameterName string
	Units string
	Data DataStoreUnit
}

type DataStoreUnit struct {
	ValueIndex   int
	Int8Value    int8
	Int16Value   int16
	Int32Value   int32
	Int64Value   int64
	Uint8Value   uint8
	Uint16Value  uint16
	Uint32Value  uint32
	Uint64Value  uint64
	FloatValue   float32
	Float64Value float64
}

type RealTimeDataStoreUnit struct {
	RxTime int64
	IsStale bool
	Units string
	ValueIndex int
	Int64Value int64
	Uint64Value uint64
	Float64Value float64
}

type RealTimeStreamElement struct {
	PacketName string
	ParameterName string
	Data RealTimeDataStoreUnit
}

type RealTimeDataBundle struct{
	Data []RealTimeStreamElement
}

type Command struct {
	Node string
	PacketType int32
	Data []byte
}

type ReceiversCoordination struct {
	Call chan bool
	Ack chan bool
	Done chan bool
}

type ServiceStatus struct {
	DataStoreMutex sync.RWMutex
	DataStoreManagerRunning bool

	GRPCMutex sync.RWMutex
	GRPCServerRunning bool

	BroadcasterMutex sync.RWMutex
	BroadcasterRunning bool

	GSLoggerMutex sync.RWMutex
	GSLoggerRunning bool

	PortsListeningMutex sync.RWMutex
	PortsListening map[int]bool
}

func NewServiceStatus() ServiceStatus {
	serviceStatus := ServiceStatus{
		DataStoreManagerRunning: false,
		GRPCServerRunning:true,
		BroadcasterRunning:false,
		GSLoggerRunning:false,
		PortsListening: map[int]bool{}}
	return serviceStatus
}