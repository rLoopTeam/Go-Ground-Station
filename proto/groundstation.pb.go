// Code generated by protoc-gen-go. DO NOT EDIT.
// source: groundstation.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	groundstation.proto

It has these top-level messages:
	StreamRequest
	ServerControl
	StartLogCommand
	Command
	SimCommand
	Ack
	ServerStatus
	OpenPort
	Pong
	Ping
	DataBundle
	DataPacket
	Params
	Value
	SimInit
	ConfigParameter
	SimConfigListRequest
	SimConfigList
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type ServerControl_CommandEnum int32

const (
	ServerControl_LogServiceStart       ServerControl_CommandEnum = 0
	ServerControl_LogServiceStop        ServerControl_CommandEnum = 1
	ServerControl_DataStoreManagerStart ServerControl_CommandEnum = 2
	ServerControl_DataStoreManagerStop  ServerControl_CommandEnum = 3
	ServerControl_BroadcasterStart      ServerControl_CommandEnum = 4
	ServerControl_BroadcasterStop       ServerControl_CommandEnum = 5
)

var ServerControl_CommandEnum_name = map[int32]string{
	0: "LogServiceStart",
	1: "LogServiceStop",
	2: "DataStoreManagerStart",
	3: "DataStoreManagerStop",
	4: "BroadcasterStart",
	5: "BroadcasterStop",
}
var ServerControl_CommandEnum_value = map[string]int32{
	"LogServiceStart":       0,
	"LogServiceStop":        1,
	"DataStoreManagerStart": 2,
	"DataStoreManagerStop":  3,
	"BroadcasterStart":      4,
	"BroadcasterStop":       5,
}

func (x ServerControl_CommandEnum) String() string {
	return proto1.EnumName(ServerControl_CommandEnum_name, int32(x))
}
func (ServerControl_CommandEnum) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

type SimCommand_SimCommandEnum int32

const (
	SimCommand_RunSimulator   SimCommand_SimCommandEnum = 0
	SimCommand_PauseSimulator SimCommand_SimCommandEnum = 1
	// ResumeSimulator = 2;
	SimCommand_StopSimulator SimCommand_SimCommandEnum = 2
	SimCommand_StartPush     SimCommand_SimCommandEnum = 3
)

var SimCommand_SimCommandEnum_name = map[int32]string{
	0: "RunSimulator",
	1: "PauseSimulator",
	2: "StopSimulator",
	3: "StartPush",
}
var SimCommand_SimCommandEnum_value = map[string]int32{
	"RunSimulator":   0,
	"PauseSimulator": 1,
	"StopSimulator":  2,
	"StartPush":      3,
}

func (x SimCommand_SimCommandEnum) String() string {
	return proto1.EnumName(SimCommand_SimCommandEnum_name, int32(x))
}
func (SimCommand_SimCommandEnum) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{4, 0} }

type StreamRequest struct {
	All        bool     `protobuf:"varint,1,opt,name=All" json:"All,omitempty"`
	Parameters []string `protobuf:"bytes,2,rep,name=Parameters" json:"Parameters,omitempty"`
}

func (m *StreamRequest) Reset()                    { *m = StreamRequest{} }
func (m *StreamRequest) String() string            { return proto1.CompactTextString(m) }
func (*StreamRequest) ProtoMessage()               {}
func (*StreamRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *StreamRequest) GetAll() bool {
	if m != nil {
		return m.All
	}
	return false
}

func (m *StreamRequest) GetParameters() []string {
	if m != nil {
		return m.Parameters
	}
	return nil
}

type ServerControl struct {
	Command ServerControl_CommandEnum `protobuf:"varint,1,opt,name=Command,enum=proto.ServerControl_CommandEnum" json:"Command,omitempty"`
}

func (m *ServerControl) Reset()                    { *m = ServerControl{} }
func (m *ServerControl) String() string            { return proto1.CompactTextString(m) }
func (*ServerControl) ProtoMessage()               {}
func (*ServerControl) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ServerControl) GetCommand() ServerControl_CommandEnum {
	if m != nil {
		return m.Command
	}
	return ServerControl_LogServiceStart
}

type StartLogCommand struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *StartLogCommand) Reset()                    { *m = StartLogCommand{} }
func (m *StartLogCommand) String() string            { return proto1.CompactTextString(m) }
func (*StartLogCommand) ProtoMessage()               {}
func (*StartLogCommand) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *StartLogCommand) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Command struct {
	Origin      string  `protobuf:"bytes,1,opt,name=Origin" json:"Origin,omitempty"`
	Node        string  `protobuf:"bytes,2,opt,name=Node" json:"Node,omitempty"`
	CommandName string  `protobuf:"bytes,3,opt,name=CommandName" json:"CommandName,omitempty"`
	CommandId   int32   `protobuf:"varint,4,opt,name=CommandId" json:"CommandId,omitempty"`
	PacketType  int32   `protobuf:"varint,5,opt,name=PacketType" json:"PacketType,omitempty"`
	Data        []int32 `protobuf:"varint,6,rep,packed,name=Data" json:"Data,omitempty"`
}

func (m *Command) Reset()                    { *m = Command{} }
func (m *Command) String() string            { return proto1.CompactTextString(m) }
func (*Command) ProtoMessage()               {}
func (*Command) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Command) GetOrigin() string {
	if m != nil {
		return m.Origin
	}
	return ""
}

func (m *Command) GetNode() string {
	if m != nil {
		return m.Node
	}
	return ""
}

func (m *Command) GetCommandName() string {
	if m != nil {
		return m.CommandName
	}
	return ""
}

func (m *Command) GetCommandId() int32 {
	if m != nil {
		return m.CommandId
	}
	return 0
}

func (m *Command) GetPacketType() int32 {
	if m != nil {
		return m.PacketType
	}
	return 0
}

func (m *Command) GetData() []int32 {
	if m != nil {
		return m.Data
	}
	return nil
}

type SimCommand struct {
	Command SimCommand_SimCommandEnum `protobuf:"varint,1,opt,name=Command,enum=proto.SimCommand_SimCommandEnum" json:"Command,omitempty"`
}

func (m *SimCommand) Reset()                    { *m = SimCommand{} }
func (m *SimCommand) String() string            { return proto1.CompactTextString(m) }
func (*SimCommand) ProtoMessage()               {}
func (*SimCommand) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *SimCommand) GetCommand() SimCommand_SimCommandEnum {
	if m != nil {
		return m.Command
	}
	return SimCommand_RunSimulator
}

type Ack struct {
	Success bool   `protobuf:"varint,1,opt,name=Success" json:"Success,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=Message" json:"Message,omitempty"`
}

func (m *Ack) Reset()                    { *m = Ack{} }
func (m *Ack) String() string            { return proto1.CompactTextString(m) }
func (*Ack) ProtoMessage()               {}
func (*Ack) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Ack) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *Ack) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type ServerStatus struct {
	DataStoreManagerRunning bool        `protobuf:"varint,1,opt,name=DataStoreManagerRunning" json:"DataStoreManagerRunning,omitempty"`
	GRPCServerRunning       bool        `protobuf:"varint,2,opt,name=GRPCServerRunning" json:"GRPCServerRunning,omitempty"`
	BroadcasterRunning      bool        `protobuf:"varint,3,opt,name=BroadcasterRunning" json:"BroadcasterRunning,omitempty"`
	GSLoggerRunning         bool        `protobuf:"varint,4,opt,name=GSLoggerRunning" json:"GSLoggerRunning,omitempty"`
	OpenPorts               []*OpenPort `protobuf:"bytes,5,rep,name=OpenPorts" json:"OpenPorts,omitempty"`
}

func (m *ServerStatus) Reset()                    { *m = ServerStatus{} }
func (m *ServerStatus) String() string            { return proto1.CompactTextString(m) }
func (*ServerStatus) ProtoMessage()               {}
func (*ServerStatus) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *ServerStatus) GetDataStoreManagerRunning() bool {
	if m != nil {
		return m.DataStoreManagerRunning
	}
	return false
}

func (m *ServerStatus) GetGRPCServerRunning() bool {
	if m != nil {
		return m.GRPCServerRunning
	}
	return false
}

func (m *ServerStatus) GetBroadcasterRunning() bool {
	if m != nil {
		return m.BroadcasterRunning
	}
	return false
}

func (m *ServerStatus) GetGSLoggerRunning() bool {
	if m != nil {
		return m.GSLoggerRunning
	}
	return false
}

func (m *ServerStatus) GetOpenPorts() []*OpenPort {
	if m != nil {
		return m.OpenPorts
	}
	return nil
}

type OpenPort struct {
	Port    int32 `protobuf:"varint,1,opt,name=Port" json:"Port,omitempty"`
	Serving bool  `protobuf:"varint,2,opt,name=Serving" json:"Serving,omitempty"`
}

func (m *OpenPort) Reset()                    { *m = OpenPort{} }
func (m *OpenPort) String() string            { return proto1.CompactTextString(m) }
func (*OpenPort) ProtoMessage()               {}
func (*OpenPort) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *OpenPort) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *OpenPort) GetServing() bool {
	if m != nil {
		return m.Serving
	}
	return false
}

type Pong struct {
	Status *ServerStatus `protobuf:"bytes,1,opt,name=Status" json:"Status,omitempty"`
}

func (m *Pong) Reset()                    { *m = Pong{} }
func (m *Pong) String() string            { return proto1.CompactTextString(m) }
func (*Pong) ProtoMessage()               {}
func (*Pong) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *Pong) GetStatus() *ServerStatus {
	if m != nil {
		return m.Status
	}
	return nil
}

type Ping struct {
}

func (m *Ping) Reset()                    { *m = Ping{} }
func (m *Ping) String() string            { return proto1.CompactTextString(m) }
func (*Ping) ProtoMessage()               {}
func (*Ping) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

type DataBundle struct {
	Parameters []*Params `protobuf:"bytes,1,rep,name=Parameters" json:"Parameters,omitempty"`
}

func (m *DataBundle) Reset()                    { *m = DataBundle{} }
func (m *DataBundle) String() string            { return proto1.CompactTextString(m) }
func (*DataBundle) ProtoMessage()               {}
func (*DataBundle) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *DataBundle) GetParameters() []*Params {
	if m != nil {
		return m.Parameters
	}
	return nil
}

type DataPacket struct {
	PacketName string    `protobuf:"bytes,1,opt,name=PacketName" json:"PacketName,omitempty"`
	PacketType int32     `protobuf:"varint,2,opt,name=PacketType" json:"PacketType,omitempty"`
	RxTime     int64     `protobuf:"varint,3,opt,name=RxTime" json:"RxTime,omitempty"`
	Parameters []*Params `protobuf:"bytes,4,rep,name=Parameters" json:"Parameters,omitempty"`
}

func (m *DataPacket) Reset()                    { *m = DataPacket{} }
func (m *DataPacket) String() string            { return proto1.CompactTextString(m) }
func (*DataPacket) ProtoMessage()               {}
func (*DataPacket) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *DataPacket) GetPacketName() string {
	if m != nil {
		return m.PacketName
	}
	return ""
}

func (m *DataPacket) GetPacketType() int32 {
	if m != nil {
		return m.PacketType
	}
	return 0
}

func (m *DataPacket) GetRxTime() int64 {
	if m != nil {
		return m.RxTime
	}
	return 0
}

func (m *DataPacket) GetParameters() []*Params {
	if m != nil {
		return m.Parameters
	}
	return nil
}

type Params struct {
	PacketName string `protobuf:"bytes,1,opt,name=PacketName" json:"PacketName,omitempty"`
	ParamName  string `protobuf:"bytes,2,opt,name=ParamName" json:"ParamName,omitempty"`
	RxTime     int64  `protobuf:"varint,3,opt,name=RxTime" json:"RxTime,omitempty"`
	Value      *Value `protobuf:"bytes,4,opt,name=Value" json:"Value,omitempty"`
	Units      string `protobuf:"bytes,5,opt,name=Units" json:"Units,omitempty"`
}

func (m *Params) Reset()                    { *m = Params{} }
func (m *Params) String() string            { return proto1.CompactTextString(m) }
func (*Params) ProtoMessage()               {}
func (*Params) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *Params) GetPacketName() string {
	if m != nil {
		return m.PacketName
	}
	return ""
}

func (m *Params) GetParamName() string {
	if m != nil {
		return m.ParamName
	}
	return ""
}

func (m *Params) GetRxTime() int64 {
	if m != nil {
		return m.RxTime
	}
	return 0
}

func (m *Params) GetValue() *Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Params) GetUnits() string {
	if m != nil {
		return m.Units
	}
	return ""
}

type Value struct {
	Index       int32   `protobuf:"varint,1,opt,name=Index" json:"Index,omitempty"`
	Int64Value  int64   `protobuf:"varint,2,opt,name=Int64Value" json:"Int64Value,omitempty"`
	Uint64Value uint64  `protobuf:"varint,3,opt,name=Uint64Value" json:"Uint64Value,omitempty"`
	DoubleValue float64 `protobuf:"fixed64,4,opt,name=DoubleValue" json:"DoubleValue,omitempty"`
}

func (m *Value) Reset()                    { *m = Value{} }
func (m *Value) String() string            { return proto1.CompactTextString(m) }
func (*Value) ProtoMessage()               {}
func (*Value) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *Value) GetIndex() int32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *Value) GetInt64Value() int64 {
	if m != nil {
		return m.Int64Value
	}
	return 0
}

func (m *Value) GetUint64Value() uint64 {
	if m != nil {
		return m.Uint64Value
	}
	return 0
}

func (m *Value) GetDoubleValue() float64 {
	if m != nil {
		return m.DoubleValue
	}
	return 0
}

type SimInit struct {
	ConfigFiles  []string           `protobuf:"bytes,1,rep,name=config_files,json=configFiles" json:"config_files,omitempty"`
	ConfigParams []*ConfigParameter `protobuf:"bytes,2,rep,name=config_params,json=configParams" json:"config_params,omitempty"`
	OutputDir    string             `protobuf:"bytes,3,opt,name=output_dir,json=outputDir" json:"output_dir,omitempty"`
}

func (m *SimInit) Reset()                    { *m = SimInit{} }
func (m *SimInit) String() string            { return proto1.CompactTextString(m) }
func (*SimInit) ProtoMessage()               {}
func (*SimInit) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *SimInit) GetConfigFiles() []string {
	if m != nil {
		return m.ConfigFiles
	}
	return nil
}

func (m *SimInit) GetConfigParams() []*ConfigParameter {
	if m != nil {
		return m.ConfigParams
	}
	return nil
}

func (m *SimInit) GetOutputDir() string {
	if m != nil {
		return m.OutputDir
	}
	return ""
}

type ConfigParameter struct {
	ConfigPath string `protobuf:"bytes,1,opt,name=config_path,json=configPath" json:"config_path,omitempty"`
	Value      string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *ConfigParameter) Reset()                    { *m = ConfigParameter{} }
func (m *ConfigParameter) String() string            { return proto1.CompactTextString(m) }
func (*ConfigParameter) ProtoMessage()               {}
func (*ConfigParameter) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

func (m *ConfigParameter) GetConfigPath() string {
	if m != nil {
		return m.ConfigPath
	}
	return ""
}

func (m *ConfigParameter) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type SimConfigListRequest struct {
}

func (m *SimConfigListRequest) Reset()                    { *m = SimConfigListRequest{} }
func (m *SimConfigListRequest) String() string            { return proto1.CompactTextString(m) }
func (*SimConfigListRequest) ProtoMessage()               {}
func (*SimConfigListRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

type SimConfigList struct {
	ConfigNames []string `protobuf:"bytes,1,rep,name=configNames" json:"configNames,omitempty"`
}

func (m *SimConfigList) Reset()                    { *m = SimConfigList{} }
func (m *SimConfigList) String() string            { return proto1.CompactTextString(m) }
func (*SimConfigList) ProtoMessage()               {}
func (*SimConfigList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

func (m *SimConfigList) GetConfigNames() []string {
	if m != nil {
		return m.ConfigNames
	}
	return nil
}

func init() {
	proto1.RegisterType((*StreamRequest)(nil), "proto.StreamRequest")
	proto1.RegisterType((*ServerControl)(nil), "proto.ServerControl")
	proto1.RegisterType((*StartLogCommand)(nil), "proto.StartLogCommand")
	proto1.RegisterType((*Command)(nil), "proto.Command")
	proto1.RegisterType((*SimCommand)(nil), "proto.SimCommand")
	proto1.RegisterType((*Ack)(nil), "proto.Ack")
	proto1.RegisterType((*ServerStatus)(nil), "proto.ServerStatus")
	proto1.RegisterType((*OpenPort)(nil), "proto.OpenPort")
	proto1.RegisterType((*Pong)(nil), "proto.Pong")
	proto1.RegisterType((*Ping)(nil), "proto.Ping")
	proto1.RegisterType((*DataBundle)(nil), "proto.DataBundle")
	proto1.RegisterType((*DataPacket)(nil), "proto.DataPacket")
	proto1.RegisterType((*Params)(nil), "proto.Params")
	proto1.RegisterType((*Value)(nil), "proto.Value")
	proto1.RegisterType((*SimInit)(nil), "proto.SimInit")
	proto1.RegisterType((*ConfigParameter)(nil), "proto.ConfigParameter")
	proto1.RegisterType((*SimConfigListRequest)(nil), "proto.SimConfigListRequest")
	proto1.RegisterType((*SimConfigList)(nil), "proto.SimConfigList")
	proto1.RegisterEnum("proto.ServerControl_CommandEnum", ServerControl_CommandEnum_name, ServerControl_CommandEnum_value)
	proto1.RegisterEnum("proto.SimCommand_SimCommandEnum", SimCommand_SimCommandEnum_name, SimCommand_SimCommandEnum_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for GroundStationService service

type GroundStationServiceClient interface {
	StreamPackets(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (GroundStationService_StreamPacketsClient, error)
	SendCommand(ctx context.Context, in *Command, opts ...grpc.CallOption) (*Ack, error)
	SendSimCommand(ctx context.Context, in *SimCommand, opts ...grpc.CallOption) (*Ack, error)
	Ping(ctx context.Context, in *Ping, opts ...grpc.CallOption) (*Pong, error)
	ControlServer(ctx context.Context, in *ServerControl, opts ...grpc.CallOption) (*Ack, error)
	InitSim(ctx context.Context, in *SimInit, opts ...grpc.CallOption) (*Ack, error)
	RequestSimConfigList(ctx context.Context, in *SimConfigListRequest, opts ...grpc.CallOption) (*SimConfigList, error)
}

type groundStationServiceClient struct {
	cc *grpc.ClientConn
}

func NewGroundStationServiceClient(cc *grpc.ClientConn) GroundStationServiceClient {
	return &groundStationServiceClient{cc}
}

func (c *groundStationServiceClient) StreamPackets(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (GroundStationService_StreamPacketsClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_GroundStationService_serviceDesc.Streams[0], c.cc, "/proto.GroundStationService/streamPackets", opts...)
	if err != nil {
		return nil, err
	}
	x := &groundStationServiceStreamPacketsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GroundStationService_StreamPacketsClient interface {
	Recv() (*DataBundle, error)
	grpc.ClientStream
}

type groundStationServiceStreamPacketsClient struct {
	grpc.ClientStream
}

func (x *groundStationServiceStreamPacketsClient) Recv() (*DataBundle, error) {
	m := new(DataBundle)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *groundStationServiceClient) SendCommand(ctx context.Context, in *Command, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := grpc.Invoke(ctx, "/proto.GroundStationService/sendCommand", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groundStationServiceClient) SendSimCommand(ctx context.Context, in *SimCommand, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := grpc.Invoke(ctx, "/proto.GroundStationService/sendSimCommand", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groundStationServiceClient) Ping(ctx context.Context, in *Ping, opts ...grpc.CallOption) (*Pong, error) {
	out := new(Pong)
	err := grpc.Invoke(ctx, "/proto.GroundStationService/ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groundStationServiceClient) ControlServer(ctx context.Context, in *ServerControl, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := grpc.Invoke(ctx, "/proto.GroundStationService/controlServer", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groundStationServiceClient) InitSim(ctx context.Context, in *SimInit, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := grpc.Invoke(ctx, "/proto.GroundStationService/InitSim", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groundStationServiceClient) RequestSimConfigList(ctx context.Context, in *SimConfigListRequest, opts ...grpc.CallOption) (*SimConfigList, error) {
	out := new(SimConfigList)
	err := grpc.Invoke(ctx, "/proto.GroundStationService/RequestSimConfigList", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for GroundStationService service

type GroundStationServiceServer interface {
	StreamPackets(*StreamRequest, GroundStationService_StreamPacketsServer) error
	SendCommand(context.Context, *Command) (*Ack, error)
	SendSimCommand(context.Context, *SimCommand) (*Ack, error)
	Ping(context.Context, *Ping) (*Pong, error)
	ControlServer(context.Context, *ServerControl) (*Ack, error)
	InitSim(context.Context, *SimInit) (*Ack, error)
	RequestSimConfigList(context.Context, *SimConfigListRequest) (*SimConfigList, error)
}

func RegisterGroundStationServiceServer(s *grpc.Server, srv GroundStationServiceServer) {
	s.RegisterService(&_GroundStationService_serviceDesc, srv)
}

func _GroundStationService_StreamPackets_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GroundStationServiceServer).StreamPackets(m, &groundStationServiceStreamPacketsServer{stream})
}

type GroundStationService_StreamPacketsServer interface {
	Send(*DataBundle) error
	grpc.ServerStream
}

type groundStationServiceStreamPacketsServer struct {
	grpc.ServerStream
}

func (x *groundStationServiceStreamPacketsServer) Send(m *DataBundle) error {
	return x.ServerStream.SendMsg(m)
}

func _GroundStationService_SendCommand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Command)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroundStationServiceServer).SendCommand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GroundStationService/SendCommand",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroundStationServiceServer).SendCommand(ctx, req.(*Command))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroundStationService_SendSimCommand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SimCommand)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroundStationServiceServer).SendSimCommand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GroundStationService/SendSimCommand",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroundStationServiceServer).SendSimCommand(ctx, req.(*SimCommand))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroundStationService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Ping)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroundStationServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GroundStationService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroundStationServiceServer).Ping(ctx, req.(*Ping))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroundStationService_ControlServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServerControl)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroundStationServiceServer).ControlServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GroundStationService/ControlServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroundStationServiceServer).ControlServer(ctx, req.(*ServerControl))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroundStationService_InitSim_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SimInit)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroundStationServiceServer).InitSim(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GroundStationService/InitSim",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroundStationServiceServer).InitSim(ctx, req.(*SimInit))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroundStationService_RequestSimConfigList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SimConfigListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroundStationServiceServer).RequestSimConfigList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GroundStationService/RequestSimConfigList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroundStationServiceServer).RequestSimConfigList(ctx, req.(*SimConfigListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _GroundStationService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.GroundStationService",
	HandlerType: (*GroundStationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "sendCommand",
			Handler:    _GroundStationService_SendCommand_Handler,
		},
		{
			MethodName: "sendSimCommand",
			Handler:    _GroundStationService_SendSimCommand_Handler,
		},
		{
			MethodName: "ping",
			Handler:    _GroundStationService_Ping_Handler,
		},
		{
			MethodName: "controlServer",
			Handler:    _GroundStationService_ControlServer_Handler,
		},
		{
			MethodName: "InitSim",
			Handler:    _GroundStationService_InitSim_Handler,
		},
		{
			MethodName: "RequestSimConfigList",
			Handler:    _GroundStationService_RequestSimConfigList_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "streamPackets",
			Handler:       _GroundStationService_StreamPackets_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "groundstation.proto",
}

func init() { proto1.RegisterFile("groundstation.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 997 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x55, 0xdd, 0x6f, 0x1b, 0x45,
	0x10, 0xf7, 0xf9, 0x6c, 0xa7, 0x1e, 0xc7, 0x1f, 0xd9, 0x98, 0x60, 0x42, 0x01, 0xb3, 0x12, 0xc2,
	0x52, 0x69, 0x44, 0x1d, 0x84, 0x0a, 0xe5, 0x25, 0x4d, 0x20, 0x44, 0xa4, 0xa9, 0xb5, 0xd7, 0x22,
	0xde, 0xaa, 0xad, 0xbd, 0xbd, 0xac, 0xe2, 0xdb, 0x3d, 0xee, 0xf6, 0xaa, 0xf2, 0x88, 0x78, 0xe1,
	0x1d, 0x09, 0xf1, 0xcc, 0x1b, 0xff, 0x17, 0xff, 0x07, 0x68, 0x3f, 0xee, 0xc3, 0x97, 0x46, 0x79,
	0xba, 0x9d, 0xdf, 0xfc, 0x66, 0x6e, 0x76, 0x66, 0x76, 0x06, 0x76, 0xc3, 0x44, 0x66, 0x62, 0x95,
	0x2a, 0xaa, 0xb8, 0x14, 0x07, 0x71, 0x22, 0x95, 0x44, 0x6d, 0xf3, 0xc1, 0x47, 0xd0, 0x0f, 0x54,
	0xc2, 0x68, 0x44, 0xd8, 0xcf, 0x19, 0x4b, 0x15, 0x1a, 0x81, 0x7f, 0xb4, 0x5e, 0x4f, 0xbc, 0xa9,
	0x37, 0xbb, 0x43, 0xf4, 0x11, 0x7d, 0x08, 0xb0, 0xa0, 0x09, 0x8d, 0x98, 0x62, 0x49, 0x3a, 0x69,
	0x4e, 0xfd, 0x59, 0x97, 0x54, 0x10, 0xfc, 0xaf, 0x07, 0xfd, 0x80, 0x25, 0xaf, 0x59, 0x72, 0x2c,
	0x85, 0x4a, 0xe4, 0x1a, 0x7d, 0x0d, 0x5b, 0xc7, 0x32, 0x8a, 0xa8, 0x58, 0x19, 0x3f, 0x83, 0xf9,
	0xd4, 0xfe, 0xf4, 0x60, 0x83, 0x76, 0xe0, 0x38, 0xdf, 0x8a, 0x2c, 0x22, 0xb9, 0x01, 0xfe, 0xd3,
	0x83, 0x5e, 0x45, 0x81, 0x76, 0x61, 0x78, 0x2e, 0x43, 0x6d, 0xc8, 0x97, 0x2c, 0x50, 0x34, 0x51,
	0xa3, 0x06, 0x42, 0x30, 0xa8, 0x82, 0x32, 0x1e, 0x79, 0xe8, 0x3d, 0x78, 0xe7, 0x84, 0x2a, 0x1a,
	0x28, 0x99, 0xb0, 0x27, 0x54, 0xd0, 0x90, 0x25, 0x96, 0xde, 0x44, 0x13, 0x18, 0x5f, 0x57, 0xc9,
	0x78, 0xe4, 0xa3, 0x31, 0x8c, 0x1e, 0x27, 0x92, 0xae, 0x96, 0x34, 0x55, 0x39, 0xbf, 0xa5, 0xff,
	0xb9, 0x81, 0xca, 0x78, 0xd4, 0xc6, 0x9f, 0xc0, 0xd0, 0xe8, 0xcf, 0x65, 0xe8, 0xe2, 0x43, 0x08,
	0x5a, 0x82, 0x46, 0xcc, 0x5c, 0xb2, 0x4b, 0xcc, 0x19, 0xff, 0xe3, 0x15, 0x97, 0x47, 0x7b, 0xd0,
	0x79, 0x9a, 0xf0, 0x90, 0x0b, 0xc7, 0x70, 0x92, 0xb6, 0xbb, 0x90, 0x2b, 0x36, 0x69, 0x5a, 0x3b,
	0x7d, 0x46, 0xd3, 0xe2, 0xda, 0x17, 0xda, 0xa5, 0x6f, 0x54, 0x55, 0x08, 0xdd, 0x85, 0xae, 0x13,
	0xcf, 0x56, 0x93, 0xd6, 0xd4, 0x9b, 0xb5, 0x49, 0x09, 0xd8, 0x2a, 0x2d, 0xaf, 0x98, 0x7a, 0xf6,
	0x4b, 0xcc, 0x26, 0x6d, 0xa3, 0xae, 0x20, 0xfa, 0x9f, 0x3a, 0x07, 0x93, 0xce, 0xd4, 0x9f, 0xb5,
	0x89, 0x39, 0xe3, 0xbf, 0x3d, 0x80, 0x80, 0x47, 0x79, 0xb8, 0x37, 0x97, 0xad, 0xe0, 0x54, 0x8e,
	0x9b, 0x65, 0xfb, 0x09, 0x06, 0x9b, 0x2a, 0x34, 0x82, 0x6d, 0x92, 0x89, 0x80, 0x47, 0xd9, 0x9a,
	0x2a, 0x99, 0xd8, 0xaa, 0x2d, 0x68, 0x96, 0xb2, 0x12, 0xf3, 0xd0, 0x8e, 0xee, 0x3f, 0x19, 0x97,
	0x50, 0x13, 0xf5, 0xa1, 0x6b, 0x12, 0xbd, 0xc8, 0xd2, 0xcb, 0x91, 0x8f, 0xbf, 0x02, 0xff, 0x68,
	0x79, 0x85, 0x26, 0xb0, 0x15, 0x64, 0xcb, 0x25, 0x4b, 0x53, 0xd7, 0x9b, 0xb9, 0xa8, 0x35, 0x4f,
	0x58, 0x9a, 0xd2, 0x30, 0x4f, 0x68, 0x2e, 0xe2, 0xff, 0x3c, 0xd8, 0xb6, 0x2d, 0x17, 0x28, 0xaa,
	0xb2, 0x14, 0x3d, 0x84, 0x77, 0xeb, 0x8d, 0x40, 0x32, 0x21, 0xb8, 0x08, 0x9d, 0xd3, 0x9b, 0xd4,
	0xe8, 0x33, 0xd8, 0x39, 0x25, 0x8b, 0x63, 0xeb, 0x2d, 0xb7, 0x69, 0x1a, 0x9b, 0xeb, 0x0a, 0x74,
	0x00, 0xa8, 0xd2, 0x40, 0x39, 0xdd, 0x37, 0xf4, 0xb7, 0x68, 0xd0, 0x0c, 0x86, 0xa7, 0xc1, 0xb9,
	0x0c, 0x2b, 0xf1, 0xb4, 0x0c, 0xb9, 0x0e, 0xa3, 0xfb, 0xd0, 0x7d, 0x1a, 0x33, 0xb1, 0x90, 0x89,
	0x4a, 0x27, 0xed, 0xa9, 0x3f, 0xeb, 0xcd, 0x87, 0xae, 0x4a, 0x39, 0x4e, 0x4a, 0x06, 0x7e, 0x08,
	0x77, 0x72, 0x41, 0x77, 0x80, 0xfe, 0x9a, 0x9b, 0xb6, 0x89, 0x39, 0x9b, 0xac, 0xea, 0x57, 0x54,
	0x5c, 0x26, 0x17, 0xf1, 0xa1, 0x66, 0x8b, 0x10, 0xdd, 0x83, 0x8e, 0x4d, 0x9e, 0xb1, 0xeb, 0xcd,
	0x77, 0x37, 0x9e, 0xb2, 0x55, 0x11, 0x47, 0xc1, 0x1d, 0x68, 0x2d, 0xb4, 0xf1, 0x23, 0x00, 0x9d,
	0xc8, 0xc7, 0x99, 0x58, 0xad, 0x19, 0xba, 0xbf, 0x31, 0x40, 0x3c, 0x13, 0x74, 0xdf, 0xb9, 0x31,
	0x8a, 0x74, 0x63, 0x9e, 0xfc, 0xe1, 0x59, 0x6b, 0xdb, 0xbc, 0x65, 0x63, 0x5f, 0x94, 0x4f, 0xad,
	0x82, 0xd4, 0x1a, 0xbf, 0x79, 0xad, 0xf1, 0xf7, 0xa0, 0x43, 0xde, 0x3c, 0xe3, 0xee, 0x4d, 0xf9,
	0xc4, 0x49, 0xb5, 0xa8, 0x5a, 0xb7, 0x45, 0xf5, 0x97, 0x07, 0x1d, 0x0b, 0xdf, 0x1a, 0xd1, 0x5d,
	0xe8, 0x1a, 0xa6, 0x51, 0xdb, 0x96, 0x2c, 0x81, 0x1b, 0xe3, 0xc1, 0xd0, 0xfe, 0x91, 0xae, 0x33,
	0x66, 0x2a, 0xdf, 0x9b, 0x6f, 0xbb, 0x50, 0x0c, 0x46, 0xac, 0x0a, 0x8d, 0xa1, 0xfd, 0x5c, 0x70,
	0x53, 0x79, 0xed, 0xd5, 0x0a, 0xf8, 0x57, 0x0f, 0x4a, 0xfd, 0x99, 0x58, 0xb1, 0x37, 0xae, 0xc6,
	0x56, 0xd0, 0xf1, 0x9e, 0x09, 0xf5, 0xe5, 0x17, 0xd6, 0x7d, 0xd3, 0xfc, 0xb5, 0x82, 0xe8, 0xd1,
	0xf3, 0x9c, 0x97, 0x04, 0x1d, 0x56, 0x8b, 0x54, 0x21, 0xcd, 0x38, 0x91, 0xd9, 0xcb, 0x35, 0x2b,
	0x23, 0xf4, 0x48, 0x15, 0xc2, 0xbf, 0x7b, 0xb0, 0x15, 0xf0, 0xe8, 0x4c, 0x70, 0x85, 0x3e, 0x86,
	0xed, 0xa5, 0x14, 0xaf, 0x78, 0xf8, 0xe2, 0x15, 0x5f, 0x33, 0x5b, 0xf1, 0x2e, 0xe9, 0x59, 0xec,
	0x3b, 0x0d, 0xa1, 0x47, 0xd0, 0x77, 0x94, 0xd8, 0xe4, 0xd4, 0xac, 0x95, 0xde, 0x7c, 0xcf, 0x5d,
	0xfa, 0xd8, 0xe8, 0x8a, 0xec, 0x13, 0xe7, 0xcf, 0xe5, 0xff, 0x03, 0x00, 0x99, 0xa9, 0x38, 0x53,
	0x2f, 0x56, 0x3c, 0x71, 0x93, 0xb2, 0x6b, 0x91, 0x13, 0x9e, 0xe0, 0xef, 0x61, 0x58, 0xb3, 0x47,
	0x1f, 0x41, 0xaf, 0xf8, 0x9d, 0xba, 0xcc, 0x4b, 0x96, 0x3b, 0x55, 0x97, 0x3a, 0x71, 0xaf, 0x8b,
	0xec, 0x74, 0x89, 0x15, 0xf0, 0x1e, 0x8c, 0xcd, 0x50, 0xd3, 0xb4, 0x73, 0x9e, 0x2a, 0xb7, 0x23,
	0xf1, 0x03, 0xe8, 0x6f, 0xe0, 0x3a, 0x3f, 0xd6, 0x99, 0xae, 0x70, 0xed, 0xc2, 0x06, 0x9a, 0xff,
	0xe6, 0xc3, 0xf8, 0xd4, 0xac, 0xe1, 0xc0, 0xae, 0x61, 0xb7, 0xbc, 0xd0, 0x37, 0xd0, 0x4f, 0xcd,
	0x02, 0xb6, 0x0d, 0x94, 0xa2, 0x71, 0xfe, 0xc0, 0xaa, 0x6b, 0x79, 0x7f, 0xc7, 0xa1, 0xe5, 0xb3,
	0xc2, 0x8d, 0xcf, 0x3d, 0x74, 0x0f, 0x7a, 0x29, 0x13, 0xab, 0x7c, 0x82, 0x0f, 0x8a, 0xfc, 0x19,
	0x79, 0x1f, 0x9c, 0x7c, 0xb4, 0xbc, 0xc2, 0x0d, 0xf4, 0x00, 0x06, 0x9a, 0x5c, 0x99, 0xf8, 0x3b,
	0xd7, 0x06, 0x7c, 0xcd, 0x04, 0x43, 0x2b, 0xd6, 0x63, 0xa7, 0x97, 0x3f, 0x0c, 0x2e, 0xc2, 0xfd,
	0x42, 0x90, 0x22, 0xc4, 0x0d, 0x74, 0x68, 0x6a, 0xa9, 0x37, 0xba, 0x9d, 0x09, 0xe5, 0x0d, 0xaa,
	0xdb, 0xbe, 0xe6, 0xf8, 0x53, 0xd8, 0xd2, 0xbd, 0x12, 0xf0, 0xa8, 0x08, 0xda, 0xb5, 0x4f, 0x8d,
	0xf8, 0x03, 0x8c, 0x5d, 0x0e, 0x36, 0x53, 0xfe, 0x7e, 0x35, 0xf4, 0x5a, 0x81, 0xf6, 0xc7, 0x6f,
	0x53, 0xe2, 0xc6, 0xcb, 0x8e, 0x81, 0x0f, 0xff, 0x0f, 0x00, 0x00, 0xff, 0xff, 0x3a, 0xb7, 0x80,
	0xee, 0x12, 0x09, 0x00, 0x00,
}
