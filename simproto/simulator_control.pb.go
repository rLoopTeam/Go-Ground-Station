// Code generated by protoc-gen-go. DO NOT EDIT.
// source: simulator_control.proto

/*
Package simproto is a generated protocol buffer package.

It is generated from these files:
	simulator_control.proto

It has these top-level messages:
	SimCommand
	Ack
	SimInit
	PusherCommand
	Parameter
*/
package simproto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type SimCommand_SimCommandEnum int32

const (
	SimCommand_RunSimulator   SimCommand_SimCommandEnum = 0
	SimCommand_PauseSimulator SimCommand_SimCommandEnum = 1
	SimCommand_StopSimulator  SimCommand_SimCommandEnum = 2
	SimCommand_StartPush      SimCommand_SimCommandEnum = 3
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
	return proto.EnumName(SimCommand_SimCommandEnum_name, int32(x))
}
func (SimCommand_SimCommandEnum) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type PusherCommand_PusherCommandEnum int32

const (
	PusherCommand_StartPush PusherCommand_PusherCommandEnum = 0
)

var PusherCommand_PusherCommandEnum_name = map[int32]string{
	0: "StartPush",
}
var PusherCommand_PusherCommandEnum_value = map[string]int32{
	"StartPush": 0,
}

func (x PusherCommand_PusherCommandEnum) String() string {
	return proto.EnumName(PusherCommand_PusherCommandEnum_name, int32(x))
}
func (PusherCommand_PusherCommandEnum) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor0, []int{3, 0}
}

type SimCommand struct {
	Command SimCommand_SimCommandEnum `protobuf:"varint,1,opt,name=command,enum=simproto.SimCommand_SimCommandEnum" json:"command,omitempty"`
}

func (m *SimCommand) Reset()                    { *m = SimCommand{} }
func (m *SimCommand) String() string            { return proto.CompactTextString(m) }
func (*SimCommand) ProtoMessage()               {}
func (*SimCommand) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *SimCommand) GetCommand() SimCommand_SimCommandEnum {
	if m != nil {
		return m.Command
	}
	return SimCommand_RunSimulator
}

type Ack struct {
	Success bool   `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
}

func (m *Ack) Reset()                    { *m = Ack{} }
func (m *Ack) String() string            { return proto.CompactTextString(m) }
func (*Ack) ProtoMessage()               {}
func (*Ack) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

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

type SimInit struct {
	ConfigFiles []string `protobuf:"bytes,1,rep,name=config_files,json=configFiles" json:"config_files,omitempty"`
	OutputDir   string   `protobuf:"bytes,2,opt,name=output_dir,json=outputDir" json:"output_dir,omitempty"`
}

func (m *SimInit) Reset()                    { *m = SimInit{} }
func (m *SimInit) String() string            { return proto.CompactTextString(m) }
func (*SimInit) ProtoMessage()               {}
func (*SimInit) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *SimInit) GetConfigFiles() []string {
	if m != nil {
		return m.ConfigFiles
	}
	return nil
}

func (m *SimInit) GetOutputDir() string {
	if m != nil {
		return m.OutputDir
	}
	return ""
}

type PusherCommand struct {
	Command PusherCommand_PusherCommandEnum `protobuf:"varint,1,opt,name=command,enum=simproto.PusherCommand_PusherCommandEnum" json:"command,omitempty"`
}

func (m *PusherCommand) Reset()                    { *m = PusherCommand{} }
func (m *PusherCommand) String() string            { return proto.CompactTextString(m) }
func (*PusherCommand) ProtoMessage()               {}
func (*PusherCommand) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *PusherCommand) GetCommand() PusherCommand_PusherCommandEnum {
	if m != nil {
		return m.Command
	}
	return PusherCommand_StartPush
}

type Parameter struct {
	Config []string `protobuf:"bytes,1,rep,name=config" json:"config,omitempty"`
}

func (m *Parameter) Reset()                    { *m = Parameter{} }
func (m *Parameter) String() string            { return proto.CompactTextString(m) }
func (*Parameter) ProtoMessage()               {}
func (*Parameter) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Parameter) GetConfig() []string {
	if m != nil {
		return m.Config
	}
	return nil
}

func init() {
	proto.RegisterType((*SimCommand)(nil), "simproto.SimCommand")
	proto.RegisterType((*Ack)(nil), "simproto.Ack")
	proto.RegisterType((*SimInit)(nil), "simproto.SimInit")
	proto.RegisterType((*PusherCommand)(nil), "simproto.PusherCommand")
	proto.RegisterType((*Parameter)(nil), "simproto.Parameter")
	proto.RegisterEnum("simproto.SimCommand_SimCommandEnum", SimCommand_SimCommandEnum_name, SimCommand_SimCommandEnum_value)
	proto.RegisterEnum("simproto.PusherCommand_PusherCommandEnum", PusherCommand_PusherCommandEnum_name, PusherCommand_PusherCommandEnum_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for SimControl service

type SimControlClient interface {
	ControlSim(ctx context.Context, in *SimCommand, opts ...grpc.CallOption) (*Ack, error)
	InitSim(ctx context.Context, in *SimInit, opts ...grpc.CallOption) (*Ack, error)
	ControlPusher(ctx context.Context, in *PusherCommand, opts ...grpc.CallOption) (*Ack, error)
	EditConfig(ctx context.Context, in *Parameter, opts ...grpc.CallOption) (*Ack, error)
}

type simControlClient struct {
	cc *grpc.ClientConn
}

func NewSimControlClient(cc *grpc.ClientConn) SimControlClient {
	return &simControlClient{cc}
}

func (c *simControlClient) ControlSim(ctx context.Context, in *SimCommand, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := grpc.Invoke(ctx, "/simproto.SimControl/ControlSim", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simControlClient) InitSim(ctx context.Context, in *SimInit, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := grpc.Invoke(ctx, "/simproto.SimControl/InitSim", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simControlClient) ControlPusher(ctx context.Context, in *PusherCommand, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := grpc.Invoke(ctx, "/simproto.SimControl/ControlPusher", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simControlClient) EditConfig(ctx context.Context, in *Parameter, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := grpc.Invoke(ctx, "/simproto.SimControl/EditConfig", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for SimControl service

type SimControlServer interface {
	ControlSim(context.Context, *SimCommand) (*Ack, error)
	InitSim(context.Context, *SimInit) (*Ack, error)
	ControlPusher(context.Context, *PusherCommand) (*Ack, error)
	EditConfig(context.Context, *Parameter) (*Ack, error)
}

func RegisterSimControlServer(s *grpc.Server, srv SimControlServer) {
	s.RegisterService(&_SimControl_serviceDesc, srv)
}

func _SimControl_ControlSim_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SimCommand)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimControlServer).ControlSim(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/simproto.SimControl/ControlSim",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimControlServer).ControlSim(ctx, req.(*SimCommand))
	}
	return interceptor(ctx, in, info, handler)
}

func _SimControl_InitSim_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SimInit)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimControlServer).InitSim(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/simproto.SimControl/InitSim",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimControlServer).InitSim(ctx, req.(*SimInit))
	}
	return interceptor(ctx, in, info, handler)
}

func _SimControl_ControlPusher_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PusherCommand)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimControlServer).ControlPusher(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/simproto.SimControl/ControlPusher",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimControlServer).ControlPusher(ctx, req.(*PusherCommand))
	}
	return interceptor(ctx, in, info, handler)
}

func _SimControl_EditConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Parameter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimControlServer).EditConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/simproto.SimControl/EditConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimControlServer).EditConfig(ctx, req.(*Parameter))
	}
	return interceptor(ctx, in, info, handler)
}

var _SimControl_serviceDesc = grpc.ServiceDesc{
	ServiceName: "simproto.SimControl",
	HandlerType: (*SimControlServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ControlSim",
			Handler:    _SimControl_ControlSim_Handler,
		},
		{
			MethodName: "InitSim",
			Handler:    _SimControl_InitSim_Handler,
		},
		{
			MethodName: "ControlPusher",
			Handler:    _SimControl_ControlPusher_Handler,
		},
		{
			MethodName: "EditConfig",
			Handler:    _SimControl_EditConfig_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "simulator_control.proto",
}

func init() { proto.RegisterFile("simulator_control.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 371 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x51, 0x3d, 0x6f, 0xe2, 0x40,
	0x10, 0xc5, 0x20, 0x61, 0x3c, 0x87, 0x11, 0xcc, 0x9d, 0x0e, 0x74, 0xd2, 0x49, 0xdc, 0xd2, 0x70,
	0x45, 0x28, 0xa0, 0xa2, 0x48, 0x81, 0x08, 0x91, 0xa2, 0x34, 0xc8, 0x6e, 0xd2, 0x21, 0xc7, 0x2c,
	0x64, 0x05, 0xeb, 0x45, 0xfb, 0x21, 0xe5, 0x0f, 0xe5, 0x77, 0xe5, 0xaf, 0x44, 0xfe, 0x58, 0x19,
	0xc7, 0x49, 0x37, 0xef, 0xcd, 0xbc, 0xdd, 0x37, 0x6f, 0x60, 0xa8, 0x18, 0x37, 0xe7, 0x48, 0x0b,
	0xb9, 0x8b, 0x45, 0xa2, 0xa5, 0x38, 0xcf, 0x2e, 0x52, 0x68, 0x81, 0x1d, 0xc5, 0x78, 0x56, 0x91,
	0x37, 0x07, 0x20, 0x64, 0x7c, 0x2d, 0x38, 0x8f, 0x92, 0x3d, 0xde, 0x82, 0x1b, 0xe7, 0xe5, 0xc8,
	0x19, 0x3b, 0xd3, 0xde, 0x7c, 0x32, 0xb3, 0xa3, 0xb3, 0x72, 0xec, 0xaa, 0xdc, 0x24, 0x86, 0x07,
	0x56, 0x43, 0x9e, 0xa0, 0x57, 0x6d, 0x61, 0x1f, 0xba, 0x81, 0x49, 0x42, 0xeb, 0xa3, 0xdf, 0x40,
	0x84, 0xde, 0x36, 0x32, 0x8a, 0x96, 0x9c, 0x83, 0x03, 0xf0, 0x43, 0x2d, 0x2e, 0x25, 0xd5, 0x44,
	0x1f, 0xbc, 0x50, 0x47, 0x52, 0x6f, 0x8d, 0x7a, 0xe9, 0xb7, 0xc8, 0x12, 0x5a, 0xab, 0xf8, 0x84,
	0x23, 0x70, 0x95, 0x89, 0x63, 0xaa, 0x54, 0xe6, 0xaf, 0x13, 0x58, 0x98, 0x76, 0x38, 0x55, 0x2a,
	0x3a, 0xd2, 0x51, 0x73, 0xec, 0x4c, 0xbd, 0xc0, 0x42, 0xf2, 0x08, 0x6e, 0xc8, 0xf8, 0x43, 0xc2,
	0x34, 0xfe, 0x83, 0x6e, 0x2c, 0x92, 0x03, 0x3b, 0xee, 0x0e, 0xec, 0x4c, 0xd3, 0x37, 0x5a, 0x53,
	0x2f, 0xf8, 0x91, 0x73, 0xf7, 0x29, 0x85, 0x7f, 0x01, 0x84, 0xd1, 0x17, 0xa3, 0x77, 0x7b, 0x26,
	0x8b, 0xa7, 0xbc, 0x9c, 0xb9, 0x63, 0x92, 0xbc, 0x82, 0x9f, 0x3a, 0xa2, 0xd2, 0x26, 0xb6, 0xfe,
	0x9c, 0xd8, 0xff, 0x32, 0xb1, 0xca, 0x64, 0x15, 0x55, 0x73, 0x23, 0x30, 0xa8, 0x75, 0xab, 0x09,
	0x34, 0xc8, 0x04, 0xbc, 0x6d, 0x24, 0x23, 0x4e, 0x35, 0x95, 0xf8, 0x1b, 0xda, 0xb9, 0xe9, 0x62,
	0x85, 0x02, 0xcd, 0xdf, 0xed, 0x39, 0xb3, 0x6b, 0xe3, 0x02, 0xa0, 0x28, 0x43, 0xc6, 0xf1, 0xd7,
	0x57, 0xb7, 0xfc, 0xe3, 0x97, 0xec, 0x2a, 0x3e, 0x91, 0x06, 0xde, 0x80, 0x9b, 0x86, 0x95, 0x2a,
	0x06, 0x15, 0x45, 0xca, 0xd6, 0xc7, 0x97, 0xe0, 0x17, 0x7f, 0xe4, 0x2b, 0xe0, 0xf0, 0x9b, 0x00,
	0xea, 0xd2, 0x39, 0xc0, 0x66, 0xcf, 0xf4, 0x3a, 0xf3, 0x8e, 0x3f, 0xaf, 0x74, 0x76, 0xd1, 0x9a,
	0xe6, 0xb9, 0x9d, 0x81, 0xc5, 0x47, 0x00, 0x00, 0x00, 0xff, 0xff, 0xab, 0xf3, 0x05, 0xac, 0xdc,
	0x02, 0x00, 0x00,
}