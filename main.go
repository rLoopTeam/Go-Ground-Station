package main

import (
	"rloop/Go-Ground-Station/datastore"
	"rloop/Go-Ground-Station/gsgrpc"
	"rloop/Go-Ground-Station/gstypes"
	"rloop/Go-Ground-Station/server"
	"rloop/Go-Ground-Station/logging"
	"fmt"
)

func main() {
	fmt.Println("Backend version 13-04-2018")
	gsLogger, loggerChannel := logging.New()
	//the channel that will be used to transfer data between the parser and the datastoremanager
	packetStoreChannel := make(chan gstypes.PacketStoreElement,512)
	commandChannel := make(chan gstypes.Command,32)

	serviceManager, serviceChan := server.NewServiceManager()
	//struct that will contain the channels that will be used to communicate between the datastoremanager and stream server
	grpcChannelsHolder := gsgrpc.GetChannelsHolder()
	//Create the datastoremanager server
	dataStoreManager := datastore.New(grpcChannelsHolder,packetStoreChannel)
	//create the broadcasting server that will send the commands to the rpod
	udpBroadCasterServer := server.CreateNewUDPCommandServer(commandChannel)
	//Create the UDPListenerServers that will listen to the packets sent by the rpod
	udpListenerServers := server.CreateNewUDPServers(packetStoreChannel,loggerChannel)
	//Create the gsgrpc stream server
	conn, grpcServer, err := gsgrpc.NewGoGrpcServer(grpcChannelsHolder,commandChannel,serviceChan)

	serviceManager.SetDatastoreManager(dataStoreManager)
	serviceManager.SetGsLogger(gsLogger)
	serviceManager.SetUDPListenerServers(udpListenerServers)
	serviceManager.SetUDPBroadcaster(udpBroadCasterServer)
	if err != nil {
		panic("unable to start gsgrpc server")
	}else{
		serviceManager.SetGrpcServer(grpcServer,conn)
	}
	serviceManager.RunAll()
	select {}
}
