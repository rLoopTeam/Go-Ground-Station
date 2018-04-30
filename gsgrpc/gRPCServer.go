package gsgrpc

import (
	"rloop/Go-Ground-Station/proto"
	"rloop/Go-Ground-Station/gstypes"
	"fmt"
	"sync"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"golang.org/x/net/context"
	"bytes"
	"encoding/binary"
	"rloop/Go-Ground-Station/simproto"
)

type GRPCServer struct {
	serviceChan chan <- *proto.ServerControl
	commandChannel chan <- gstypes.Command
	simCommandChannel chan <- *simproto.SimCommand
	receiversChannelHolder *ChannelsHolder
	statusProvider StatusProvider
}

type ChannelsHolder struct {
	ReceiverMutex sync.Mutex
	//struct that will prevent multiple operations on the channelholder at the same time, sort of mutex
	Coordinator gstypes.ReceiversCoordination
	Receivers map[*chan gstypes.RealTimeDataBundle]*chan gstypes.RealTimeDataBundle
}

func (srv *GRPCServer) Ping(context.Context, *proto.Ping) (*proto.Pong, error) {
	srvStat := srv.statusProvider.GetStatus()
	length := len(srvStat.PortsListening)
	openPorts := make([]*proto.OpenPort,length)
	idx := 0
	for port,serving := range srvStat.PortsListening{
		openPorts[idx] = &proto.OpenPort{
			Port:int32(port),
			Serving:serving}
		idx++
	}
	status := &proto.ServerStatus{
		DataStoreManagerRunning:srvStat.DataStoreManagerRunning,
		GRPCServerRunning:srvStat.GRPCServerRunning,
		BroadcasterRunning:srvStat.BroadcasterRunning,
		GSLoggerRunning:srvStat.GSLoggerRunning,
		OpenPorts:openPorts}
	return &proto.Pong{Status:status},nil
}

func (srv *GRPCServer) StreamPackets(req *proto.StreamRequest, stream proto.GroundStationService_StreamPacketsServer) error {
	var err error
	receiverChannel := make( chan gstypes.RealTimeDataBundle,32)
	srv.addChannelToDatastoreQueue(receiverChannel)
	fmt.Println("gsgrpc channel pushed to map")
		for element := range receiverChannel{
			dataBundle := proto.DataBundle{}
			dataArray := make([]*proto.Params,len(element.Data))
			for idx := 0; idx < len(element.Data); idx++{
				param := proto.Params{}
				param.RxTime = element.Data[idx].Data.RxTime
				param.ParamName = element.Data[idx].ParameterName
				param.PacketName = element.Data[idx].PacketName
				switch element.Data[idx].Data.ValueIndex {
					case 1:	param.Value = &proto.Value{Index: 1, Int64Value:element.Data[idx].Data.Int64Value}
					case 2: param.Value = &proto.Value{Index: 2, Uint64Value:element.Data[idx].Data.Uint64Value}
					case 3: param.Value = &proto.Value{Index: 3, DoubleValue:element.Data[idx].Data.Float64Value}
				}
				dataArray[idx] = &param
			}
			dataBundle.Parameters = dataArray
			err = stream.Send(&dataBundle)
			if err != nil {
				srv.removeChannelFromDatastoreQueue(receiverChannel)
				break
			}else {
				fmt.Println("sent data to frontend server")
			}
		}
	return err
}

func (srv *GRPCServer) SendCommand(ctx context.Context, cmd *proto.Command) (*proto.Ack, error) {
	ack := &proto.Ack{}
	var err error
	fmt.Printf("Request for command: %v\n", cmd)
	node := cmd.Node
	packetType := cmd.PacketType
	data := cmd.Data
	dataLength := len(data)

	dataBytesArray := [][]byte {{0,0,0,0},{0,0,0,0},{0,0,0,0},{0,0,0,0}}

	//convert the data values to bytes
	for idx := 0; idx < dataLength; idx++ {
		buf := new(bytes.Buffer)
		value := data[idx]
		err := binary.Write(buf, binary.LittleEndian, value)
		if err != nil {
			ack = nil
			break
		}else{
			dataBytesArray[idx] = buf.Bytes()
		}
	}
	//if there's no data or not enough data populate the remaining byte slots with zero value
	for idx := dataLength; idx < 4; idx++{
		var value int32 = 0
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, value)
		if err != nil {
			ack = nil
			break
		}else{
			dataBytesArray[idx] = buf.Bytes()
		}
	}

	dataBytes := append(dataBytesArray[0],dataBytesArray[1]...)
	dataBytes = append(dataBytes,dataBytesArray[2]...)
	dataBytes = append(dataBytes,dataBytesArray[3]...)

	if err == nil {
		command := gstypes.Command{
			Node: node,
			PacketType:packetType,
			Data:dataBytes,
		}
		srv.commandChannel <- command
	}

	return ack, err
}

func (srv *GRPCServer) ControlServer(ctx context.Context, control *proto.ServerControl) (*proto.Ack, error) {
	srv.serviceChan <- control
	return &proto.Ack{}, nil
}

func (srv *GRPCServer) SendSimCommand(ctx context.Context,command *proto.SimCommand) (*proto.Ack, error) {
	var convertedValue simproto.SimCommand_SimCommandEnum
	simCommand := &simproto.SimCommand{}
	cmdName := command.Command.String()
	cmdValue := simproto.SimCommand_SimCommandEnum_value[cmdName]
	convertedValue = simproto.SimCommand_SimCommandEnum(cmdValue)
	simCommand.Command = convertedValue
	srv.simCommandChannel<- simCommand
	return &proto.Ack{},nil
}

func (srv *GRPCServer ) addChannelToDatastoreQueue(receiverChannel chan gstypes.RealTimeDataBundle){
	srv.receiversChannelHolder.Coordinator.Call <- true
	<- srv.receiversChannelHolder.Coordinator.Ack
	srv.receiversChannelHolder.Receivers[&receiverChannel] = &receiverChannel
	srv.receiversChannelHolder.Coordinator.Done <- true
}

func (srv *GRPCServer) removeChannelFromDatastoreQueue(receiverChannel chan gstypes.RealTimeDataBundle){
	srv.receiversChannelHolder.Coordinator.Call <- true
	<- srv.receiversChannelHolder.Coordinator.Ack
	delete(srv.receiversChannelHolder.Receivers, &receiverChannel)
	fmt.Println("closing receiver channel")
	srv.receiversChannelHolder.Coordinator.Done <- true
}

func GetChannelsHolder () *ChannelsHolder {
	callChannel := make (chan bool)
	ackChannel := make (chan bool)
	doneChannel := make (chan bool)
	coordinator := gstypes.ReceiversCoordination{
		Call:callChannel,
		Ack:ackChannel,
		Done:doneChannel}

	holder := &ChannelsHolder{
		Coordinator:coordinator,
		ReceiverMutex: sync.Mutex{},
		Receivers: make(map[*chan gstypes.RealTimeDataBundle]*chan gstypes.RealTimeDataBundle),
	}
	return holder
}

func newGroundStationGrpcServer (grpcChannelsHolder *ChannelsHolder,commandChannel chan <- gstypes.Command, simCommandChannel chan<-*simproto.SimCommand, serviceChan chan<- *proto.ServerControl, statusProvider StatusProvider) *GRPCServer{
	srv := &GRPCServer{
		receiversChannelHolder:grpcChannelsHolder,
		commandChannel:commandChannel,
		serviceChan:serviceChan,
		statusProvider:statusProvider,
		simCommandChannel:simCommandChannel}
	return srv
}

func NewGoGrpcServer (port int, grpcChannelsHolder *ChannelsHolder, commandChannel chan <- gstypes.Command, simCommandChannel chan<-*simproto.SimCommand, serviceChan chan<- *proto.ServerControl,statusProvider StatusProvider) (net.Listener, *grpc.Server, error){
	GSserver := newGroundStationGrpcServer(grpcChannelsHolder, commandChannel, simCommandChannel, serviceChan, statusProvider)
	var err error
	var grpcServer *grpc.Server

	//initialize grpcserver
	strPort := ":" + strconv.Itoa(port)
	conn, err := net.Listen("tcp", strPort)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}else {
		grpcServer = grpc.NewServer()
		proto.RegisterGroundStationServiceServer(grpcServer,GSserver)
	}

	return conn,grpcServer,err
}

type StatusProvider interface {
	GetStatus() gstypes.ServiceStatus
}