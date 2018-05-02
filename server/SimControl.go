package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"rloop/Go-Ground-Station/proto"
	"rloop/Go-Ground-Station/simproto"
)

type SimController struct {
	doRun       bool
	IsRunning   bool
	conn        *grpc.ClientConn
	client      simproto.SimControlClient
	signalChan  chan bool
	commandChan <-chan *simproto.SimCommand
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
		case cmd := <-client.commandChan:
			client.SendCommand(cmd)
		case <-client.signalChan:
			break MainLoop
		}
	}
	client.conn.Close()
	client.IsRunning = false
}

func (client *SimController) SendCommand(cmd *simproto.SimCommand) *proto.Ack {
	rack := &proto.Ack{}
	if client.conn == nil {
		log.Fatalf("Cannot send sim command: Connection is not set \n")
		rack.Success = false
		rack.Message = "Cannot send sim command: Connection is not set"
		return rack
	}
	ack, err := client.client.ControlSim(context.Background(), cmd)
	fmt.Printf("sending command: %v \n", cmd.Command)
	if err != nil {
		log.Printf("SimControl send failed: %v \n", err)
		rack.Success = false
	} else {
		rack.Success = ack.Success
	}
	fmt.Printf("Sim Controller Response: %s\n", ack.Message)
	return rack
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

func NewSimController() (*SimController, chan<- *simproto.SimCommand) {
	signalCh := make(chan bool)
	commandCh := make(chan *simproto.SimCommand)
	controller := &SimController{
		signalChan:  signalCh,
		commandChan: commandCh,
		IsRunning:   false,
		doRun:       false}
	return controller, commandCh
}
