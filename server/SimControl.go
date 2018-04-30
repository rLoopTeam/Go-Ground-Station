package server

import (
	"google.golang.org/grpc"
	"log"
	"fmt"
	"context"
	"rloop/Go-Ground-Station/simproto"
	"rloop/Go-Ground-Station/proto"
)

type SimController struct {
	doRun bool
	IsRunning bool
	conn *grpc.ClientConn
	client simproto.SimControlServiceClient
	signalChan chan bool
	commandChan <- chan *simproto.SimCommand
}

func(client *SimController) Stop() {
	client.doRun = false
	client.signalChan <- true
}

func (client *SimController) Run(){
	if client.conn == nil {
		fmt.Printf("Sim controller grpc connection is not set \n")
		return
	}
	client.IsRunning = true
	for {
		select{
		case cmd := <-client.commandChan: client.SendCommand(cmd)
		case <-client.signalChan:break
		}
	}
	client.conn.Close()
	client.IsRunning = false;
}

func (client *SimController) SendCommand(cmd *simproto.SimCommand) *proto.Ack{
	rack := &proto.Ack{}
	ack, err := client.client.SendSimCommand(context.Background(),cmd)
	if err != nil {
		log.Printf("SimControl send failed: %v \n",err)
		rack.CommandExecuted = false
	}else{
		rack.CommandExecuted = ack.CommandExecuted
	}
	return rack
}

func (client *SimController) Connect(address string){
	var err error
	client.conn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client.client = simproto.NewSimControlServiceClient(client.conn)
}

func NewSimController() (*SimController,chan <- *simproto.SimCommand){
	signalCh := make(chan bool)
	commandCh := make(chan *simproto.SimCommand)
	controller := &SimController{
		signalChan:signalCh,
		commandChan:commandCh,
		IsRunning:false,
		doRun:false}
	return controller,commandCh
}