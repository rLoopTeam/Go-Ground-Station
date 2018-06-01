package gsgrpc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"rloop/Go-Ground-Station/gstypes"
	"rloop/Go-Ground-Station/helpers"
	"rloop/Go-Ground-Station/proto"
	"strconv"
	"sync"
	//"time"
)

type GRPCServer struct {
	serviceChan            chan<- gstypes.ServerControlWithTimeout
	commandChannel         chan<- gstypes.Command
	simCommandChannel      chan<- *gstypes.SimulatorCommandWithResponse
	simInitChannel         chan<- *gstypes.SimulatorInitWithResponse
	receiversChannelHolder *ChannelsHolder
	statusProvider         StatusProvider
}

type ChannelsHolder struct {
	ReceiverMutex sync.Mutex
	//map that holds the channels to communicate with the grpc clients
	Receivers map[*chan gstypes.DataStoreBundle]*chan gstypes.DataStoreBundle
}

func (srv *GRPCServer) Ping(context.Context, *proto.Ping) (*proto.Pong, error) {
	srvStatus := srv.statusProvider.GetStatus()
	length := len(srvStatus.PortsListening)
	openPorts := make([]*proto.OpenPort, length)
	idx := 0
	for port, serving := range srvStatus.PortsListening {
		openPorts[idx] = &proto.OpenPort{
			Port:    int32(port),
			Serving: serving}
		idx++
	}
	status := &proto.ServerStatus{
		DataStoreManagerRunning: srvStatus.DataStoreManagerRunning,
		GRPCServerRunning:       srvStatus.GRPCServerRunning,
		BroadcasterRunning:      srvStatus.BroadcasterRunning,
		GSLoggerRunning:         srvStatus.GSLoggerRunning,
		OpenPorts:               openPorts}
	return &proto.Pong{Status: status}, nil
}

func (srv *GRPCServer) StreamPackets(req *proto.StreamRequest, stream proto.GroundStationService_StreamPacketsServer) error {
	var err error
	var dataArrayLength int
	var pushCurrentDataStoreElement bool
	var requestedParameters map[string]struct{}
	var dataStoreBundleLength int
	var dataArrayIdx int
	receiverChannel := make(chan gstypes.DataStoreBundle, 32)
	srv.addChannelToDatastoreQueue(receiverChannel)
	dataArrayLength = len(req.Parameters)
	if !req.All {
		requestedParameters = map[string]struct{}{}
		for _, p := range req.Parameters {
			requestedParameters[p] = struct{}{}
		}
	}

MainLoop:
	for dataStoreBundle := range receiverChannel {

		dataStoreBundleLength = len(dataStoreBundle.Data)
		if req.All {
			dataArrayLength = dataStoreBundleLength
			pushCurrentDataStoreElement = true
		}
		dataBundle := proto.DataBundle{}
		dataArray := make([]*proto.Params, dataArrayLength)
		dataArrayIdx = 0
		for idx := 0; idx < dataStoreBundleLength; idx++ {
			if !req.All {
				_, pushCurrentDataStoreElement = requestedParameters[dataStoreBundle.Data[idx].FullParameterName]
			}
			//fmt.Printf("push to element: %v \n", pushCurrentDataStoreElement)
			if pushCurrentDataStoreElement {
				param := proto.Params{}
				param.RxTime = dataStoreBundle.Data[idx].RxTime
				param.ParamName = dataStoreBundle.Data[idx].FullParameterName
				param.PacketName = dataStoreBundle.Data[idx].PacketName
				switch dataStoreBundle.Data[idx].Data.ValueIndex {
				case 4:
					param.Value = &proto.Value{Index: 1, Int64Value: dataStoreBundle.Data[idx].Data.Int64Value}
				case 8:
					param.Value = &proto.Value{Index: 2, Uint64Value: dataStoreBundle.Data[idx].Data.Uint64Value}
				case 10:
					param.Value = &proto.Value{Index: 3, DoubleValue: dataStoreBundle.Data[idx].Data.Float64Value}
				}
				dataArray[dataArrayIdx] = &param
				dataArrayIdx++
			}
		}
		dataBundle.Parameters = dataArray
		//fmt.Printf("sending bundle: %v \n", dataBundle)
		//error will occur when connection is closed; in which case we remove the channel as a receiver and exit the loop
		err = stream.Send(&dataBundle)
		if err != nil {
			srv.removeChannelFromDatastoreQueue(receiverChannel)
			break MainLoop
		} else {
			fmt.Println("sent data to frontend server")
		}
	}
	return err
}

func (srv *GRPCServer) SendCommand(ctx context.Context, cmd *proto.Command) (*proto.Ack, error) {
	var ack *proto.Ack
	//the message that will be delivered to the pod
	var dataBytes []byte
	//any error encountered will be pushed into this variable and returned immediately to the sender
	var err error
	//fmt.Printf("Request for command: %v\n", cmd)
	var node string
	var packetType int32
	var sequence int32
	var data []int32
	var crc int32
	var dataLength int
	var command gstypes.Command

	ack = &proto.Ack{}
	node = cmd.Node
	packetType = cmd.PacketType
	data = cmd.Data
	dataLength = len(data)

	dataBytesArray := [][]byte{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}

	//convert the data values to bytes
	for idx := 0; idx < dataLength; idx++ {
		buf := new(bytes.Buffer)
		value := data[idx]
		err := binary.Write(buf, binary.LittleEndian, value)
		if err != nil {
			ack.Success = false
			ack.Message = err.Error()
			goto returnStatement
		} else {
			dataBytesArray[idx] = buf.Bytes()
		}
	}
	//if there's no data or not enough data populate the remaining byte slots with zero value
	//not necessary action
	/*
		for idx := dataLength; idx < 4; idx++ {
			var value int32 = 0
			buf := new(bytes.Buffer)
			err := binary.Write(buf, binary.LittleEndian, value)
			if err != nil {
				ack.Success = false
				ack.Message = err.Error()
				goto returnStatement
			} else {
				dataBytesArray[idx] = buf.Bytes()
			}
		}
	*/

	dataBytes = helpers.AppendVariadic(dataBytesArray...)
	/*
		dataBytes = append(dataBytesArray[0], dataBytesArray[1]...)
		dataBytes = append(dataBytes, dataBytesArray[2]...)
		dataBytes = append(dataBytes, dataBytesArray[3]...)
	*/

	if err == nil {
		command = gstypes.Command{
			Node:       node,
			PacketType: packetType,
			Sequence:   sequence,
			Data:       dataBytes,
			Crc:        int16(crc),
		}
		srv.commandChannel <- command
		ack.Success = true
	}
returnStatement:
	return ack, err
}

func (srv *GRPCServer) ControlServer(ctx context.Context, control *proto.ServerControl) (*proto.Ack, error) {
	var response gstypes.Ack
	var ret *proto.Ack
	var comm chan gstypes.Ack
	var controlStruct gstypes.ServerControlWithTimeout

	comm = make(chan gstypes.Ack)
	controlStruct = gstypes.ServerControlWithTimeout{
		Control:      gstypes.ServerControl_CommandEnum(control.Command),
		ResponseChan: comm,
		Ctx:          ctx}
	srv.serviceChan <- controlStruct

	select {
	case <-ctx.Done():
		fmt.Printf("context done: %v \n", ctx.Err())
	case response = <-comm:
		ret.Success = response.Success
		ret.Message = response.Message
	}
	return ret, nil
}

func (srv *GRPCServer) SendSimCommand(ctx context.Context, command *proto.SimCommand) (*proto.Ack, error) {
	var err error
	var ack *proto.Ack
	var req *gstypes.SimulatorCommandWithResponse
	var responseChan chan *proto.Ack

	responseChan = make(chan *proto.Ack)

	req = &gstypes.SimulatorCommandWithResponse{
		ResponseChan: responseChan,
		Command:      command}

	srv.simCommandChannel <- req
	ack = <-responseChan
	return ack, err
}
func (srv *GRPCServer) InitSim(ctx context.Context, in *proto.SimInit) (*proto.Ack, error) {
	var err error
	var ack *proto.Ack
	var req *gstypes.SimulatorInitWithResponse
	var responseChan chan *proto.Ack

	//fmt.Printf("request for new sim config: %v \n", in)
	responseChan = make(chan *proto.Ack)

	req = &gstypes.SimulatorInitWithResponse{
		ResponseChan: responseChan,
		SimInit:      in}

	srv.simInitChannel <- req
	ack = <-responseChan
	return ack, err
}

func (srv *GRPCServer) RequestSimConfigList(context.Context, *proto.SimConfigListRequest) (*proto.SimConfigList, error) {
	arr := []string{"test1", "test2"}
	obj := &proto.SimConfigList{
		ConfigNames: arr}

	return obj, nil
}

func (srv *GRPCServer) addChannelToDatastoreQueue(receiverChannel chan gstypes.DataStoreBundle) {
	srv.receiversChannelHolder.ReceiverMutex.Lock()
	srv.receiversChannelHolder.Receivers[&receiverChannel] = &receiverChannel
	srv.receiversChannelHolder.ReceiverMutex.Unlock()
}

func (srv *GRPCServer) removeChannelFromDatastoreQueue(receiverChannel chan gstypes.DataStoreBundle) {
	srv.receiversChannelHolder.ReceiverMutex.Lock()
	delete(srv.receiversChannelHolder.Receivers, &receiverChannel)
	fmt.Println("closing receiver channel")
	srv.receiversChannelHolder.ReceiverMutex.Unlock()
}

func GetChannelsHolder() *ChannelsHolder {
	holder := &ChannelsHolder{
		ReceiverMutex: sync.Mutex{},
		Receivers:     make(map[*chan gstypes.DataStoreBundle]*chan gstypes.DataStoreBundle),
	}
	return holder
}

func newGroundStationGrpcServer(grpcChannelsHolder *ChannelsHolder, commandChannel chan<- gstypes.Command, simCommandChannel chan<- *gstypes.SimulatorCommandWithResponse, simInitChannel chan<- *gstypes.SimulatorInitWithResponse, serviceChan chan<- gstypes.ServerControlWithTimeout, statusProvider StatusProvider) *GRPCServer {
	srv := &GRPCServer{
		receiversChannelHolder: grpcChannelsHolder,
		commandChannel:         commandChannel,
		serviceChan:            serviceChan,
		statusProvider:         statusProvider,
		simCommandChannel:      simCommandChannel,
		simInitChannel:         simInitChannel}
	return srv
}

func NewGoGrpcServer(port int, grpcChannelsHolder *ChannelsHolder, commandChannel chan<- gstypes.Command, simCommandChannel chan<- *gstypes.SimulatorCommandWithResponse, simInitChannel chan<- *gstypes.SimulatorInitWithResponse, serviceChan chan<- gstypes.ServerControlWithTimeout, statusProvider StatusProvider) (net.Listener, *grpc.Server, error) {
	var GSServer *GRPCServer
	var grpcServer *grpc.Server
	var strPort string
	var err error
	var conn net.Listener

	GSServer = newGroundStationGrpcServer(grpcChannelsHolder, commandChannel, simCommandChannel, simInitChannel, serviceChan, statusProvider)
	//initialize grpcserver
	strPort = ":" + strconv.Itoa(port)
	conn, err = net.Listen("tcp", strPort)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		grpcServer = grpc.NewServer()
		proto.RegisterGroundStationServiceServer(grpcServer, GSServer)
	}

	return conn, grpcServer, err
}

type StatusProvider interface {
	GetStatus() gstypes.ServiceStatus
}
