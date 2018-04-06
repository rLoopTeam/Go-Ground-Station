package main

import (
	"rloop/Go-Ground-Station/datastore"
	"rloop/Go-Ground-Station/gsgrpc"
	"rloop/Go-Ground-Station/gstypes"
	"rloop/Go-Ground-Station/server"
	"rloop/Go-Ground-Station/logging"
)

func main() {
	loggerChannel := make(chan gstypes.PacketStoreElement,512)
	logger := logging.Gslogger{
		DataChan:loggerChannel}
		go logger.Run()
	//the channel that will be used to transfer data between the parser and the datastoremanager
	packetStoreChannel := make(chan gstypes.PacketStoreElement,512)
	commandChannel := make(chan gstypes.Command,32)
	//struct that will contain the channels that will be used to communicate between the datastoremanager and stream server
	grpcChannelsHolder := gsgrpc.GetChannelsHolder()
	//Create the datastoremanager server
	dataStoreManager := datastore.New(grpcChannelsHolder,packetStoreChannel)
	dataStoreManager.Start()
	//create the broadcasting server that will send the commands to the rpod
	udpBroadCasterServer := server.CreateNewUDPCommandServer(commandChannel)
	udpBroadCasterServer.Run()
	//Create the UDPListenerServers that will listen to the packets sent by the rpod
	udpListenerServers := server.CreateNewUDPServers(packetStoreChannel,loggerChannel)
	countUdpServers := len(udpListenerServers)
	//run the UDPListenerServers
	for idx := 0; idx < countUdpServers; idx++{
		udpListenerServers[idx].Run()
	}
	//Create the gsgrpc stream server
	groundStationGrpcServer := gsgrpc.NewGroundStationGrpcServer(grpcChannelsHolder)
	conn, grpcServer, err := gsgrpc.NewGoGrpcServer(groundStationGrpcServer)
	if err != nil {
		panic("unable to start gsgrpc server")
	}else{
		grpcServer.Serve(conn)
	}
	select {}
}
