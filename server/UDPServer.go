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
	"sync"
)

type GSUDPServer interface {
	Run()
}

type UDPBroadcasterServer struct {
	hosts              []gstypes.Host
	isRunningMutex     sync.RWMutex
	isRunning          bool
	doRunMutex         sync.RWMutex
	doRun              bool
	signalChannel      chan bool
	commandChannel     <-chan gstypes.Command
	podCommandSequence int32
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
	s := []string{"127.0.0.2:", strconv.Itoa(port)}

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
			fmt.Printf("Problem with parsing packet on port %d in: %v \n", nodePort, r)
			//fmt.Printf("errcount on nodeport %d: %d\n", nodePort,*errcount)
		}
	}()

	element, err := parsing.ParsePacket(nodePort, nodeName, packet, errcount)
	if err == nil {
		select {
		case srv.packetStoreChannel <- element:
		default:
		}
		select {
		case srv.loggerChan <- element:
		default:
		}

	} else {
		//fmt.Println(err)
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

func (srv *UDPBroadcasterServer) ResetSequence() {
	srv.podCommandSequence = 0
}

func (srv *UDPBroadcasterServer) broadcast() {
	var cmd gstypes.Command
	//var destination *net.UDPAddr
	var conn *net.UDPConn
	var err error
	var packetBytes []byte
	nodesMap := map[string]string{}

	//populate the map with name and port, this way we can just lookup in the map and not loop through the list each time
	for _, node := range srv.hosts {
		nodesMap[node.Name] = strconv.Itoa(node.Port)
	}

	srv.isRunning = true
BroadCastLoop:
	for {
		select {
		case <-srv.signalChannel:
			fmt.Println("Broadcaster Stop")
			break BroadCastLoop
		case cmd = <-srv.commandChannel:
			goto Broadcast
		}
	Broadcast:
		var connErr error
		//retrieve the next command from the channel
		//lookup which port is to be used
		port := nodesMap[cmd.Node]
		//addr := fmt.Sprintf("127.0.0.1:%s", port)
		//try to resolve the address
		//destination, _ = net.ResolveUDPAddr("udp", addr)
		intport, _ := strconv.Atoi(port)
		//dial up, since it's udp shouldn't be a problem
		conn, err = net.DialUDP("udp", nil, &net.UDPAddr{IP: []byte{127, 0, 0, 1}, Port: intport, Zone: ""})
		if err != nil {
			fmt.Println(err)
		} else {
			//if no conflicts on address, serialize the command
			packetBytes, err = srv.serialize(cmd)
		}
		//if there's no error with serialization, send the command
		if err != nil {
			fmt.Println(err)
		} else {
			//fmt.Printf("\n sending command to node: %s on address %s \n", cmd.Node, destination.String())
			//fmt.Printf("Data: %v \n", cmd.Data)
			//fmt.Printf("command bytes: %v \n", packetBytes)
			_, connErr = conn.Write(packetBytes)
		}

		if connErr != nil {
			fmt.Printf("Command write error: %v", connErr)
		} else {
			srv.podCommandSequence++
		}
		conn.Close()
	}
	srv.isRunning = false
}

func (srv *UDPBroadcasterServer) serialize(cmd gstypes.Command) ([]byte, error) {
	var err error
	var packetTypeBytes []byte
	var sequenceBytes []byte
	var dataBytes []byte
	var crcBytes []byte
	var serializedPacket []byte
	var crcLachFunc []byte

	packetTypeBytes, err = helpers.ParseValueToBytes(cmd.PacketType)
	sequenceBytes, err = helpers.ParseValueToBytes(srv.podCommandSequence)
	dataBytes = cmd.Data
	crcBytes, err = helpers.ParseValueToBytes(cmd.Crc)

	serializedPacket = append(sequenceBytes, packetTypeBytes...)
	serializedPacket = append(serializedPacket, dataBytes...)

	crcLachFunc, err = helpers.Crc16Bytes(serializedPacket, uint32(len(serializedPacket)))

	serializedPacket = append(serializedPacket, crcLachFunc...)

	fmt.Printf("Sequence Value: %#x \n", sequenceBytes)
	fmt.Printf("PacketType Value: %#x \n", packetTypeBytes)
	fmt.Printf("Data Value: %#x \n\n", dataBytes)
	fmt.Printf("Crc Value: %#x \n", crcBytes)
	fmt.Printf("Packet Value: %#x \n\n", serializedPacket)

	fmt.Printf("Crc Value: %#x \n", crcBytes)
	fmt.Printf("Crc Custom Lach Value: %#x \n\n", crcLachFunc)

	return serializedPacket, err
}

func (srv *UDPBroadcasterServer) GetStatus() (bool, bool) {
	defer func() {
		srv.isRunningMutex.RUnlock()
		srv.doRunMutex.RUnlock()
	}()
	srv.isRunningMutex.RLock()
	srv.doRunMutex.RLock()
	return srv.isRunning, srv.doRun
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
	signalChannel := make(chan bool)
	commandChannel := make(chan gstypes.Command, 4)
	srv := &UDPBroadcasterServer{
		signalChannel:  signalChannel,
		hosts:          hosts,
		commandChannel: commandChannel,
		isRunning:      false,
		doRun:          false}
	return srv, commandChannel
}
