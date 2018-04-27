package main

import (
	"rloop/Go-Ground-Station/datastore"
	"rloop/Go-Ground-Station/gsgrpc"
	"rloop/Go-Ground-Station/server"
	"rloop/Go-Ground-Station/logging"
	"fmt"
	_ "net/http/pprof"
	"log"
	"net/http"
	"flag"
	"rloop/Go-Ground-Station/constants"
	"strconv"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:8080", nil))
	}()
	fmt.Println("Backend version 13-04-2018")

	//PARSE FLAGS AND DETERMINE WHICH PORTS THE SERVER WILL LISTEN TO
	//the array that will hold the port numbers for the UDP listeners
	flag.Parse()
	var nodesPorts []int
	flags := flag.Args()
	flagLength := len(flags)
	if flagLength > 0 {
		nodesPorts = make([]int, flagLength)
		for i, p := range flags {
			port, err := strconv.Atoi(p)
			if err == nil{
				nodesPorts[i] = port
			}
		}
	}else{
		//calculate the amount of ports that we'll have to listen to
		amountNodes := len(constants.HostsToListen)
		nodesPorts = make([]int, amountNodes)
		mapIndex := 0
		for k := range constants.HostsToListen {
			nodesPorts[mapIndex] = k
			mapIndex++
		}
	}

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
	udpListenerServers := server.CreateNewUDPListenerServers(dataStoreChannel,loggerChannel,nodesPorts)
	//Create the gsgrpc stream server
	conn, grpcServer, err := gsgrpc.NewGoGrpcServer(grpcChannelsHolder,commandChannel,serviceChannel, serviceManager)

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
