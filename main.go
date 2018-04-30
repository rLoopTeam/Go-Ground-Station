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
	"strconv"
	"rloop/Go-Ground-Station/helpers"
	"rloop/Go-Ground-Station/gstypes"
	"os"
	"rloop/Go-Ground-Station/simproto"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:8080", nil))
	}()
	fmt.Println("Backend version 13-04-2018")

	NetworkConfig,err := helpers.DecodeNetworkingFile("./config/networking.json")

	if err != nil {
		log.Fatalf("No config is defined, please define a config file: %v \n", err)
		os.Exit(1)
	}
	HostsTolisten := NetworkConfig.HostsToListen
	HostsToCommand := NetworkConfig.HostsToCommand
	GrpcPort := NetworkConfig.Grpc
	nodesMap := map[int]gstypes.Host{}
	var nodesPorts []int

	//create the nodes map, for efficiency
	arrlength := len(HostsTolisten)
	nodesPorts = make([]int, arrlength)
	for idx := 0; idx < arrlength; idx++ {
		currPort := HostsTolisten[idx].Port
		nodesMap[currPort] = HostsTolisten[idx]
		nodesPorts[idx] = currPort
	}


	//PARSE FLAGS, IF FLAGS ARE GIVEN DETERMINE WHICH PORTS THE SERVER WILL LISTEN TO ELSE DO NOTHING AND USE PORTS FROM CONFIG FILE
	//the array that will hold the port numbers for the UDP listeners
	flag.Parse()

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
	}

	//create a servicemanager and get the channel where control commands will be issued
	serviceManager, serviceChannel := server.NewServiceManager()

	var simCommandChannel chan<-*simproto.SimCommand
	var simController *server.SimController
	if NetworkConfig.WithSim {
		simController, simCommandChannel = server.NewSimController()
		serviceManager.SetSimController(simController)
	}

	gsLogger, loggerChannel := logging.New()
	//struct that will contain the channels that will be used to communicate between the datastoremanager and stream server
	grpcChannelsHolder := gsgrpc.GetChannelsHolder()
	//Create the datastoremanager server
	dataStoreManager, dataStoreChannel := datastore.New(grpcChannelsHolder)
	//create the broadcasting server that will send the commands to the rpod
	udpBroadCasterServer, commandChannel := server.CreateNewUDPCommandServer(HostsToCommand)
	//Create the UDPListenerServers that will listen to the packets sent by the rpod
	udpListenerServers := server.CreateNewUDPListenerServers(dataStoreChannel,loggerChannel,nodesPorts, nodesMap)
	//Create the gsgrpc stream server
	conn, grpcServer, err := gsgrpc.NewGoGrpcServer(GrpcPort, grpcChannelsHolder,commandChannel,simCommandChannel, serviceChannel, serviceManager)

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
