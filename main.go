package main

import (
	"rloop/Go-Ground-Station/dataStore"
	"rloop/Go-Ground-Station/gsgrpc"
	"rloop/Go-Ground-Station/gstypes"
	"rloop/Go-Ground-Station/server"
)

func main() {
	callChannel := make (chan bool)
	ackChannel := make (chan bool)
	doneChannel := make (chan bool)
	coordinator := gstypes.ReceiversCoordination{
		Call:callChannel,
		Ack:ackChannel,
		Done:doneChannel}
	//the channel that will be used to transfer data between the parser and the datastoremanager
	packetStoreChannel := make(chan gstypes.PacketStoreElement,512)
	commandChannel := make(chan gstypes.Command,32)
	//struct that will contain the channels that will be used to communicate between the datastoremanager and stream server
	grpcChannelsHolder := gsgrpc.GetChannelsHolder()
	//Create the datastoremanager server
	dataStoreManager := dataStore.New(grpcChannelsHolder,packetStoreChannel,coordinator)
	dataStoreManager.Start()
	//create the broadcasting server that will send the commands to the rpod
	udpBroadCasterServer := server.CreateNewUDPCommandServer(commandChannel)
	udpBroadCasterServer.Run()
	//Create the UDPListenerServers that will listen to the packets sent by the rpod
	udpListenerServers := server.CreateNewUDPServers(packetStoreChannel)
	countUdpServers := len(udpListenerServers)
	//run the UDPListenerServers
	for idx := 0; idx < countUdpServers; idx++{
		udpListenerServers[idx].Run()
	}
	//Create the gsgrpc stream server
	groundStationGrpcServer := gsgrpc.NewGroundStationGrpcServer(grpcChannelsHolder,coordinator)
	conn, grpcServer, err := gsgrpc.NewGoGrpcServer(groundStationGrpcServer)
	if err != nil {
		panic("unable to start gsgrpc server")
	}else{
		grpcServer.Serve(conn)
	}
	select {}
}
