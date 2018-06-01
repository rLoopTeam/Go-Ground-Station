package gstypes

import (
	"context"
	"go/types"
	"rloop/Go-Ground-Station/proto"
	"sync"
)

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

// used by packetparser
type PacketStoreElement struct {
	PacketName      string
	PacketType      int
	ParameterPrefix string
	Port            int
	RxTime          int64
	Parameters      []DataStoreElement
}

type DataStoreElement struct {
	PacketName        string
	FullParameterName string
	ParameterName     string
	RxTime            int64
	IsStale           bool
	Units             string
	Data              DataStoreUnit
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

type DataStoreBundle struct {
	Data []DataStoreElement
}

type Command struct {
	Node       string
	PacketType int32
	Data       []byte
}

type ReceiversCoordination struct {
	Call chan bool
	Ack  chan bool
	Done chan bool
}

type ServiceStatus struct {
	TimeMutex   sync.RWMutex
	LastUpdated int64

	DataStoreMutex          sync.RWMutex
	DataStoreManagerRunning bool

	GRPCMutex         sync.RWMutex
	GRPCServerRunning bool

	BroadcasterMutex   sync.RWMutex
	BroadcasterRunning bool

	GSLoggerMutex   sync.RWMutex
	GSLoggerRunning bool

	PortsListeningMutex sync.RWMutex
	PortsListening      map[int]bool
}

type Config struct {
	Networking Networking `json:Networking`
}

type Networking struct {
	HostsToListen  []Host `json:HostsToListen`
	HostsToCommand []Host `json:HostsToCommand`
	Grpc           int    `json:Grpc`
	PySim          string `json:PySim`
	WithSim        bool   `json:WithSim`
}

type Host struct {
	Ip   string `json:Ip`
	Port int    `json:Port`
	Name string `json:Name`
}

type ServerControlWithTimeout struct {
	Ctx          context.Context
	ResponseChan chan<- Ack
	Control      ServerControl_CommandEnum
}

type SimulatorInitWithResponse struct {
	SimInit      *proto.SimInit
	ResponseChan chan *proto.Ack
}

type SimulatorCommandWithResponse struct {
	Command      *proto.SimCommand
	ResponseChan chan *proto.Ack
}

type Ack struct {
	Success bool
	Message string
}

type ServerControl_CommandEnum int32

const (
	ServerControl_LogServiceStart       ServerControl_CommandEnum = 0
	ServerControl_LogServiceStop        ServerControl_CommandEnum = 1
	ServerControl_DataStoreManagerStart ServerControl_CommandEnum = 2
	ServerControl_DataStoreManagerStop  ServerControl_CommandEnum = 3
	ServerControl_BroadcasterStart      ServerControl_CommandEnum = 4
	ServerControl_BroadcasterStop       ServerControl_CommandEnum = 5
)

func NewServiceStatus() ServiceStatus {
	serviceStatus := ServiceStatus{
		LastUpdated:             0,
		DataStoreManagerRunning: false,
		GRPCServerRunning:       true,
		BroadcasterRunning:      false,
		GSLoggerRunning:         false,
		PortsListening:          map[int]bool{}}
	return serviceStatus
}
