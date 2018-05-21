package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"rloop/Go-Ground-Station/gstypes"
	"rloop/Go-Ground-Station/proto"
	"rloop/Go-Ground-Station/simproto"
	"time"
)

type SimController struct {
	doRun       bool
	IsRunning   bool
	conn        *grpc.ClientConn
	client      simproto.SimControlClient
	signalChan  chan bool
	commandChan <-chan *gstypes.SimulatorCommandWithResponse
	simInitChan  <-chan *gstypes.SimulatorInitWithResponse
}

func (client *SimController) Stop() {
	if client.IsRunning {
		client.doRun = false
		client.signalChan <- true
	}
}

func (client *SimController) Run() {
	if client.conn == nil {
		fmt.Printf("Sim controller grpc connection is not set \n")
		return
	}
	fmt.Printf("launching sim controller client \n")
	client.IsRunning = true
MainLoop:
	for {
		select {
		case param := <-client.simInitChan:
			client.InitSim(param)
			break
		case cmd := <-client.commandChan:
			client.SendCommand(cmd)
		case <-client.signalChan: //stopsignal
			break MainLoop
		}
	}
	client.conn.Close()
	client.IsRunning = false
}

func (client *SimController) SendCommand(cmd *gstypes.SimulatorCommandWithResponse) {
	//will hold converted value from main proto to simulator proto
	var convertedValue simproto.SimCommand_SimCommandEnum
	var simCommand *simproto.SimCommand
	var cmdName string
	var cmdValue int32
	var ack *simproto.Ack
	var err error
	var ctx context.Context
	var cancel context.CancelFunc

	rack := &proto.Ack{}
	if client.conn == nil {
		log.Fatalf("Cannot send sim command: Connection is not set \n")
		rack.Success = false
		rack.Message = "Cannot send sim command: Connection is not set"
		goto ResponseStatement
	}

	simCommand = &simproto.SimCommand{}
	cmdName = cmd.Command.Command.String()
	cmdValue = simproto.SimCommand_SimCommandEnum_value[cmdName]
	convertedValue = simproto.SimCommand_SimCommandEnum(cmdValue)
	simCommand.Command = convertedValue
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	fmt.Printf("sending command: %v \n", simCommand)
	ack, err = client.client.ControlSim(ctx, simCommand)

	if err != nil {
		log.Printf("SimControl send failed: %v \n", err)
		rack.Success = false
	}else if ack == nil{
		rack.Success = false
		rack.Message = ctx.Err().Error()
	} else {
		rack.Success = ack.Success
	}
	fmt.Printf("Sim Controller Response: %s\n", rack.Message)
ResponseStatement:
	cmd.ResponseChan <- rack
}

func (controller *SimController) InitSim(initWithResponse *gstypes.SimulatorInitWithResponse) {
	var rack *proto.Ack
	var ack *simproto.Ack
	var err error
	var simInitConfig *simproto.SimInit
	var config_files_arr []string
	var config_params_arr []*simproto.ConfigParameter

	simInitConfig = &simproto.SimInit{}
	rack = &proto.Ack{}
	config_files_arr = make([]string, len(initWithResponse.SimInit.ConfigFiles))
	config_params_arr = make([]*simproto.ConfigParameter, len(initWithResponse.SimInit.ConfigParams))

	if controller.conn == nil {
		log.Fatalf("Cannot send sim command: Connection is not set \n")
		rack.Success = false
		rack.Message = "Cannot send sim command: Connection is not set"
		goto ReturnStatement
	}

	for idx, cfg := range initWithResponse.SimInit.ConfigFiles {
		config_files_arr[idx] = cfg
	}

	for idx, cfg := range initWithResponse.SimInit.ConfigParams {
		config_params_arr[idx] = &simproto.ConfigParameter{
			ConfigPath:cfg.ConfigPath,
			Value:cfg.Value}
	}

	simInitConfig.ConfigParams = config_params_arr
	simInitConfig.ConfigFiles = config_files_arr
	simInitConfig.OutputDir = initWithResponse.SimInit.OutputDir

	ack, err = controller.client.InitSim(context.Background(), simInitConfig)

	if err != nil {
		rack.Success = false
		rack.Message = err.Error()
	} else {
		rack.Success = ack.Success
		rack.Message = ack.Message
	}

ReturnStatement:
	initWithResponse.ResponseChan <- rack
}

func (client *SimController) Connect(address string) {
	var err error
	client.conn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Sim Controller Client did not connect: %v\n", err)
	} else {
		client.client = simproto.NewSimControlClient(client.conn)
	}
}

func NewSimController() (*SimController, chan<- *gstypes.SimulatorCommandWithResponse, chan<- *gstypes.SimulatorInitWithResponse) {
	signalCh := make(chan bool)
	commandCh := make(chan *gstypes.SimulatorCommandWithResponse)
	simInitCh := make(chan *gstypes.SimulatorInitWithResponse)
	controller := &SimController{
		signalChan:  signalCh,
		commandChan: commandCh,
		simInitChan:  simInitCh,
		IsRunning:   false,
		doRun:       false}
	return controller, commandCh, simInitCh
}
