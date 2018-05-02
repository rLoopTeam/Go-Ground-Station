package server

import (
	"fmt"
	"log"
	"net"
	"rloop/Go-Ground-Station/gstypes"
	"rloop/Go-Ground-Station/helpers"
	"rloop/Go-Ground-Station/parsing"
	"strconv"
	"strings"
)

type GSUDPServer interface {
	Run()
}

type UDPBroadcasterServer struct {
	hosts          []gstypes.Host
	isRunning      bool
	doRun          bool
	signalChannel  chan bool
	commandChannel <-chan gstypes.Command
}

type UDPListenerServer struct {
	IsRunning          bool
	doRun              bool
	signalChannel      chan bool
	conn               *net.UDPConn
	ServerPort         int
	NodeName           string
	packetStoreChannel chan<- gstypes.PacketStoreElement
	loggerChan         chan<- gstypes.PacketStoreElement
}

func (srv *UDPListenerServer) open(port int) error {
	srv.ServerPort = port
	s := []string{":", strconv.Itoa(port)}

	udpPort := strings.Join(s, "")
	udpAddr, err := net.ResolveUDPAddr("udp4", udpPort)

	if err == nil {
		srv.conn, err = net.ListenUDP("udp4", udpAddr)
		fmt.Printf("\n listening on port: %d \n", port)
	} else {
		fmt.Printf("UDPSERVER ERROR: %v \n", err)
	}

	return err
}

func (srv *UDPListenerServer) Run() {
	srv.doRun = true
	if !srv.IsRunning {
		srv.listen()
	}
}

func (srv *UDPListenerServer) Stop() {
	if srv.IsRunning {
		srv.doRun = false
		srv.signalChannel <- true
	}
}

func (srv *UDPListenerServer) listen() {
	buffer := make([]byte, 1024)
	errCount := 0
	srv.IsRunning = true
MainLoop:
	for { /*
			select {
			case n, _, err := srv.conn.ReadFromUDP(buffer[0:]):
				if err != nil {
					fmt.Printf("Packet error on port: %d\n", srv.ServerPort)
					continue MainLoop
				}
				if n > 0 {
					srv.ProcessMessage(srv.ServerPort, srv.NodeName, buffer[:n], &errCount)
				}

			case <-srv.signalChannel:
				break MainLoop
			}
		*/
		if !srv.doRun {
			break
		}
		n, _, err := srv.conn.ReadFromUDP(buffer[0:])

		if err != nil {
			fmt.Printf("Packet error on port: %d\n", srv.ServerPort)
			continue MainLoop
		}
		if n > 0 {
			srv.ProcessMessage(srv.ServerPort, srv.NodeName, buffer[:n], &errCount)
		}

	}
	fmt.Printf("UDP SERVER RUNNING = FALSE \n")
	srv.IsRunning = false
}

func (srv *UDPListenerServer) ProcessMessage(nodePort int, nodeName string, packet []byte, errcount *int) {
	defer func() {
		if r := recover(); r != nil {
			*errcount++
			//fmt.Println("Problem with parsing packet in: ", r)
			//fmt.Printf("errcount on nodeport %d: %d\n", nodePort,*errcount)
		}
	}()

	element, err := parsing.ParsePacket(nodePort, nodeName, packet, errcount)
	if err == nil {
		srv.packetStoreChannel <- element
		srv.loggerChan <- element
	} else {
		fmt.Println(err)
	}
}

func (srv *UDPBroadcasterServer) Run() {
	srv.doRun = true
	if !srv.isRunning {
		srv.broadcast()
	}
}

func (srv *UDPBroadcasterServer) Stop() {
	if srv.isRunning {
		srv.signalChannel <- true
		fmt.Println("broadcaster signal sent")
	}
	srv.doRun = false
}

func (srv *UDPBroadcasterServer) broadcast() {
	var cmd gstypes.Command
	var destination *net.UDPAddr
	var conn *net.UDPConn
	var err error
	var packetBytes []byte
	nodesMap := map[string]string{}

	//populate the map with name and port, this way we can just lookup in the map and not loop through the list each time
	for _, node := range srv.hosts {
		nodesMap[node.Name] = strconv.Itoa(node.Port)
	}

	//fmt.Printf("hosts: %v", nodesMap)

	srv.isRunning = true
BroadCastLoop:
	for {
		select {
		case <-srv.signalChannel:
			fmt.Println("Broadcaster Stop")
			break BroadCastLoop
		case cmd = <-srv.commandChannel:
			fmt.Println("received request for command broadcast")
			goto Broadcast
		}
		/*
		select {
		case <-srv.signalChannel:
			fmt.Println("Broadcaster Stop")
			break BroadCastLoop
		default:
		}
		select {
		case cmd = <-srv.commandChannel:
			fmt.Println("received request for command broadcast")
			goto Broadcast
			default: continue
		}
		*/
	Broadcast:
		var connErr error
		//retrieve the next command from the channel
		//lookup which port is to be used
		port := nodesMap[cmd.Node]
		addr := "127.0.0.1:" + port
		//try to resolve the address
		destination, _ = net.ResolveUDPAddr("udp", addr)
		//dial up, since it's udp shouldn't be a problem
		conn, err = net.DialUDP("udp", nil, destination)
		if err != nil {
			fmt.Println(err)
			err = nil
		} else {
			//if no conflicts on address, serialize the command
			packetBytes, err = serialize(cmd)
		}
		//if there's no error with serialization, send the command
		if err != nil {
			fmt.Println(err)
			err = nil
		} else {
			//fmt.Printf("\n sending command to node: %s on address %s \n", cmd.Node, destination.String())
			_, connErr = conn.Write(packetBytes)
		}

		if connErr != nil {
			fmt.Printf("Command write error: %v", connErr)
		}
		conn.Close()
	}
	srv.isRunning = false
}

func serialize(cmd gstypes.Command) ([]byte, error) {
	var bytes []byte
	var err error

	packetType, err := helpers.ParseValueToBytes(cmd.PacketType)
	data := cmd.Data
	bytes = append(packetType, data...)

	return bytes, err
}

func CreateNewUDPListenerServers(channel chan<- gstypes.PacketStoreElement, loggerChannel chan<- gstypes.PacketStoreElement, nodesPorts []int, nodesMap map[int]gstypes.Host) []*UDPListenerServer {
	amountNodes := len(nodesPorts)
	//create an array that will keep the servers
	serversArray := make([]*UDPListenerServer, amountNodes)
	//populate the nodeports array with the port numbers

	//create and open all the servers
	for idx := 0; idx < amountNodes; idx++ {
		srv := &UDPListenerServer{
			IsRunning:          false,
			doRun:              false,
			packetStoreChannel: channel,
			loggerChan:         loggerChannel,
			NodeName:           nodesMap[nodesPorts[idx]].Name}
		err := srv.open(nodesPorts[idx])
		if err == nil {
			serversArray[idx] = srv
		} else {
			log.Fatalf("Unable to open port %d: %v", nodesPorts[idx], err)
		}
	}
	return serversArray
}

func CreateNewUDPCommandServer(hosts []gstypes.Host) (*UDPBroadcasterServer, chan<- gstypes.Command) {
	signalChannel:= make(chan bool)
	commandChannel := make(chan gstypes.Command, 4)
	srv := &UDPBroadcasterServer{
		signalChannel: signalChannel,
		hosts:          hosts,
		commandChannel: commandChannel,
		isRunning:      false,
		doRun:          false}
	return srv, commandChannel
}
