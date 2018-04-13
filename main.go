package main

import (
	"rloop/Go-Ground-Station/datastore"
	"rloop/Go-Ground-Station/gsgrpc"
	"rloop/Go-Ground-Station/server"
	"rloop/Go-Ground-Station/logging"
	"fmt"
)

func main() {
	fmt.Println("Backend version 13-04-2018")
	gsLogger, loggerChannel := logging.New()

	//struct that will contain the channels that will be used to communicate between the datastoremanager and stream server
	grpcChannelsHolder := gsgrpc.GetChannelsHolder()
	//create a servicemanager and get the channel where control commands will be issued
	serviceManager, serviceChannel := server.NewServiceManager()
	//Create the datastoremanager server
	dataStoreManager, dataStoreChannel := datastore.New(grpcChannelsHolder)
	//create the broadcasting server that will send the commands to the rpod
	udpBroadCasterServer, commandChannel := server.CreateNewUDPCommandServer()
	//Create the UDPListenerServers that will listen to the packets sent by the rpod
	udpListenerServers := server.CreateNewUDPServers(dataStoreChannel,loggerChannel)
	//Create the gsgrpc stream server
	conn, grpcServer, err := gsgrpc.NewGoGrpcServer(grpcChannelsHolder,commandChannel,serviceChannel)

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
