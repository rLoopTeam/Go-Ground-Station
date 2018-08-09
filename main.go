package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	//"net/http"
	_ "net/http/pprof"
	"os"
	"rloop/Go-Ground-Station-1/datastore"
	"rloop/Go-Ground-Station-1/gsgrpc"
	"rloop/Go-Ground-Station-1/gstypes"
	"rloop/Go-Ground-Station-1/helpers"
	"rloop/Go-Ground-Station-1/logging"
	"rloop/Go-Ground-Station-1/server"
	"strconv"
)

func main() {
	/*
		go func() {
			log.Println(http.ListenAndServe("localhost:8080", nil))
		}()
	*/
	fmt.Println("Backend version 09-08-2018")

	//errors
	var networkConfigError error
	var grpcError error

	//networking
	var networkConfig gstypes.Networking
	var hostsTolisten []gstypes.Host
	var hostsToCommand []gstypes.Host
	var GrpcPort int
	var nodesMap map[int]gstypes.Host
	var nodesPorts []int
	var grpcConn net.Listener

	//services
	var serviceManager *server.ServiceManager
	var simController *server.SimController
	var gsLogger *logging.Gslogger
	var dataStoreManager *datastore.DataStoreManager
	var udpBroadCasterServer *server.UDPBroadcasterServer
	var udpListenerServers []*server.UDPListenerServer
	var grpcServer *grpc.Server

	//channels
	var serviceChannel chan<- gstypes.ServerControlWithTimeout
	var simCommandChannel chan<- *gstypes.SimulatorCommandWithResponse
	var simInitChannel chan<- *gstypes.SimulatorInitWithResponse
	var loggerChannel chan<- gstypes.PacketStoreElement
	var grpcChannelsHolder *gsgrpc.ChannelsHolder
	var dataStoreChannel chan<- gstypes.PacketStoreElement
	var commandChannel chan<- gstypes.Command

	networkConfig, networkConfigError = helpers.DecodeNetworkingFile("./config/networking.json")

	if networkConfigError != nil {
		log.Fatalf("No config is defined, please define a config file: %v \n", networkConfigError)
		os.Exit(1)
	}
	hostsTolisten = networkConfig.HostsToListen
	hostsToCommand = networkConfig.HostsToCommand
	GrpcPort = networkConfig.Grpc
	nodesMap = map[int]gstypes.Host{}

	//create the nodes map, for efficiency
	arrlength := len(hostsTolisten)
	nodesPorts = make([]int, arrlength)
	for idx := 0; idx < arrlength; idx++ {
		currPort := hostsTolisten[idx].Port
		nodesMap[currPort] = hostsTolisten[idx]
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
			if err == nil {
				nodesPorts[i] = port
			}
		}
	}

	//create a servicemanager and get the channel where control commands will be issued
	serviceManager, serviceChannel = server.NewServiceManager()

	if networkConfig.WithSim && networkConfig.PySim != "" {
		simController, simCommandChannel, simInitChannel = server.NewSimController()
		simController.Connect(networkConfig.PySim)
		serviceManager.SetSimController(simController)
	}

	gsLogger, loggerChannel = logging.New()
	//struct that will contain the channels that will be used to communicate between the datastoremanager and stream server
	grpcChannelsHolder = gsgrpc.GetChannelsHolder()
	//Create the datastoremanager server
	dataStoreManager, dataStoreChannel = datastore.New(grpcChannelsHolder)
	//create the broadcasting server that will send the commands to the rpod
	udpBroadCasterServer, commandChannel = server.CreateNewUDPCommandServer(hostsToCommand)
	//Create the UDPListenerServers that will listen to the packets sent by the rpod
	udpListenerServers = server.CreateNewUDPListenerServers(dataStoreChannel, loggerChannel, nodesPorts, nodesMap)
	//Create the gsgrpc stream server
	grpcConn, grpcServer, grpcError = gsgrpc.NewGoGrpcServer(GrpcPort, grpcChannelsHolder, commandChannel, simCommandChannel, simInitChannel, serviceChannel, serviceManager)

	serviceManager.SetDatastoreManager(dataStoreManager)
	serviceManager.SetGsLogger(gsLogger)
	serviceManager.SetUDPListenerServers(udpListenerServers)
	serviceManager.SetUDPBroadcaster(udpBroadCasterServer)
	if grpcError != nil {
		panic("unable to start gsgrpc server")
	} else {
		serviceManager.SetGrpcServer(grpcServer, grpcConn)
	}
	serviceManager.RunAll()
	select {}
}
