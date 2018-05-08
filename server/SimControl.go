package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"rloop/Go-Ground-Station/gstypes"
	"rloop/Go-Ground-Station/proto"
	"rloop/Go-Ground-Station/simproto"
)

type SimController struct {
	doRun       bool
	IsRunning   bool
	conn        *grpc.ClientConn
	client      simproto.SimControlClient
	signalChan  chan bool
	commandChan <-chan *gstypes.SimulatorCommandWithResponse
	configChan  <-chan *gstypes.SimulatorConfigWithResponse
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
		case param := <-client.configChan:
			client.SendNewParams(param)
			break
		case cmd := <-client.commandChan:
			client.SendCommand(cmd)
		case <-client.signalChan:
			break MainLoop
		}
	}
	client.conn.Close()
	client.IsRunning = false
}

func (client *SimController) SendCommand(cmd *gstypes.SimulatorCommandWithResponse){
	var convertedValue simproto.SimCommand_SimCommandEnum
	var simCommand *simproto.SimCommand
	var cmdName string
	var cmdValue int32
	var ack *simproto.Ack
	var err error

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

	ack, err = client.client.ControlSim(context.Background(), simCommand)
	fmt.Printf("sending command: %v \n", cmd.Command)
	if err != nil {
		log.Printf("SimControl send failed: %v \n", err)
		rack.Success = false
	} else {
		rack.Success = ack.Success
	}
	fmt.Printf("Sim Controller Response: %s\n", ack.Message)
	ResponseStatement:
	cmd.ResponseChan <- rack
}

func (controller *SimController) SendNewParams(strct *gstypes.SimulatorConfigWithResponse) {
	var rack *proto.Ack
	var ack *simproto.Ack
	var err error
	var parameters *simproto.Parameter
	var parametersArr []string
	var configLineFormat string

	rack = &proto.Ack{}
	if controller.conn == nil {
		log.Fatalf("Cannot send sim command: Connection is not set \n")
		rack.Success = false
		rack.Message = "Cannot send sim command: Connection is not set"
		goto ReturnStatement
	}

	parametersArr = make([]string, len(strct.Config.Config))

	//the base format in which the parameters will be passed to the simulator
	configLineFormat = "%s: [ %s ]"
	for idx, p := range strct.Config.Config {
		str := fmt.Sprintf(configLineFormat,p.Key, p.Value)
		parametersArr[idx] = str
	}

	parameters = &simproto.Parameter{Config: parametersArr}

	ack, err = controller.client.EditConfig(context.Background(), parameters)

	if err != nil {
		rack.Success = false
		rack.Message = err.Error()
	} else {
		rack.Success = ack.Success
		rack.Message = ack.Message
	}

ReturnStatement:
	strct.ResponseChan <- rack
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

func NewSimController() (*SimController, chan<- *gstypes.SimulatorCommandWithResponse) {
	signalCh := make(chan bool)
	commandCh := make(chan *gstypes.SimulatorCommandWithResponse)
	controller := &SimController{
		signalChan:  signalCh,
		commandChan: commandCh,
		IsRunning:   false,
		doRun:       false}
	return controller, commandCh
}
