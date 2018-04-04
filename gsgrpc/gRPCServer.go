package gsgrpc

import (
	"rloop/Go-Ground-Station/proto"
	"rloop/Go-Ground-Station/gstypes"
	"fmt"
	"sync"
	"google.golang.org/grpc"
	"net"
	"rloop/Go-Ground-Station/constants"
	"strconv"
	"golang.org/x/net/context"
	"bytes"
	"encoding/binary"
)

type GRPCServer struct {
	commandChannel chan <- gstypes.Command
	receiversChannelHolder ChannelsHolder
}

func (srv *GRPCServer) StreamPackets (req *proto.StreamRequest,stream proto.GroundStationService_StreamPacketsServer) error {
	var err error
	receiverChannel := make( chan gstypes.RealTimeDataBundle,64)
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

func (srv *GRPCServer ) addChannelToDatastoreQueue(receiverChannel chan gstypes.RealTimeDataBundle){
	srv.receiversChannelHolder.Coordinator.Call <- true
	<- srv.receiversChannelHolder.Coordinator.Ack
	srv.receiversChannelHolder.Receivers[&receiverChannel] = &receiverChannel
	srv.receiversChannelHolder.Coordinator.Done <- true
}

func (srv *GRPCServer) removeChannelFromDatastoreQueue(receiverChannel chan gstypes.RealTimeDataBundle){
	srv.receiversChannelHolder.Coordinator.Call <- true
	<- srv.receiversChannelHolder.Coordinator.Ack
	close(receiverChannel)
	delete(srv.receiversChannelHolder.Receivers, &receiverChannel)
	fmt.Println("closing receiver channel")
	srv.receiversChannelHolder.Coordinator.Done <- true
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
			break;
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
			break;
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

func NewGroundStationGrpcServer (grpcChannelsHolder ChannelsHolder) *GRPCServer{
	server := &GRPCServer{
		receiversChannelHolder:grpcChannelsHolder}
	return server
}

func NewGoGrpcServer (GSserver *GRPCServer) (net.Listener, *grpc.Server, error){
	var err error
	var grpcServer *grpc.Server

	//initialize grpcserver
	strPort := ":" + strconv.Itoa(constants.GrpcPort)
	conn, err := net.Listen("tcp", strPort)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}else {
		grpcServer = grpc.NewServer()
		proto.RegisterGroundStationServiceServer(grpcServer,GSserver)
	}

	return conn,grpcServer,err
}

func GetChannelsHolder () ChannelsHolder {
	callChannel := make (chan bool)
	ackChannel := make (chan bool)
	doneChannel := make (chan bool)
	coordinator := gstypes.ReceiversCoordination{
		Call:callChannel,
		Ack:ackChannel,
		Done:doneChannel}

	holder := ChannelsHolder{
		Coordinator:coordinator,
		ReceiverMutex: sync.Mutex{},
		Receivers: make(map[*chan gstypes.RealTimeDataBundle]*chan gstypes.RealTimeDataBundle),
	}
	return holder
}

type ChannelsHolder struct {
	ReceiverMutex sync.Mutex
	//struct that will prevent multiple operations on the channelholder at the same time, sort of mutex
	Coordinator gstypes.ReceiversCoordination
	Receivers map[*chan gstypes.RealTimeDataBundle]*chan gstypes.RealTimeDataBundle
}