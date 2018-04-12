package server

import (
	"fmt"
	"net"
	"strings"
	"rloop/Go-Ground-Station/gstypes"
	"strconv"
	"rloop/Go-Ground-Station/constants"
	"log"
	"rloop/Go-Ground-Station/helpers"
	"rloop/Go-Ground-Station/parsing"
)

type GSUDPServer interface {
	Run()
}

type UDPBroadcasterServer struct {
	isRunning bool
	doRun bool
	ch <-chan gstypes.Command
}

type UDPListenerServer struct {
	isRunning bool
	doRun bool
	conn      *net.UDPConn
	serverPort int
	ch chan<- gstypes.PacketStoreElement
	loggerChan chan<-gstypes.PacketStoreElement
}

func (srv *UDPListenerServer) open(port int) error {
	srv.serverPort = port
	s := []string{":", strconv.Itoa(port)}

	udpPort := strings.Join(s, "")
	udpAddr, err := net.ResolveUDPAddr("udp4", udpPort)

	if err == nil {
		srv.conn, err = net.ListenUDP("udp4", udpAddr)
		fmt.Printf("\n listening on port: %d \n", port)
	}

	return err
}

func (srv *UDPListenerServer) Run(){
	srv.doRun = true
	if !srv.isRunning {
		srv.listen()
	}
}

func (srv *UDPListenerServer) Stop (){
	srv.doRun = false
}

func (srv *UDPListenerServer) listen() {
	buffer := make([]byte, 1024)
	errCount := 0
	srv.isRunning = true
	for {
		if !srv.doRun {break}
		n, _, err := srv.conn.ReadFromUDP(buffer[0:])

		if err != nil {
			fmt.Printf("Packet error on port: %d\n", srv.serverPort)
			return
		}
		if n > 0 {
			parsing.ParsePacket(srv.serverPort, buffer[:n], srv.ch,srv.loggerChan, &errCount)
		}
	}
	srv.isRunning = false
}

func (srv *UDPBroadcasterServer) Run (){
	srv.doRun = true
	if !srv.isRunning{
		srv.broadcast()
	}
}

func (srv *UDPBroadcasterServer) Stop (){
	srv.doRun = false
}

func (srv *UDPBroadcasterServer) broadcast (){
	var destination *net.UDPAddr
	var conn *net.UDPConn
	var err error
	var packetBytes []byte
	nodesMap := map[string]string{}

	//populate the map with name and port, this way we can just lookup in the map and not loop through the list each time
	for _, node := range constants.Hosts{
		nodesMap[node.Name] = strconv.Itoa(node.Port)
	}

	srv.isRunning = true
	for {
		if !srv.doRun {break}
		//retrieve the next command from the channel
		cmd := <-srv.ch
		//lookup which port is to be used
		port := nodesMap[cmd.Node]
		addr := "255.255.255.255:" + port
		//try to resolve the address
		destination, _ = net.ResolveUDPAddr("udp",addr)
		//dial up, since it's udp shouldn't be a problem
		conn, err = net.DialUDP("udp",nil,destination)
		if err != nil {
			fmt.Println(err)
			err = nil
		}else{
			//if no conflicts on address, serialize the command
			packetBytes,err = serialize(cmd)
		}
		//if there's no error with serialization, send the command
		if err != nil {
			fmt.Println(err)
			err = nil
		} else {
			conn.Write(packetBytes)
		}
		conn.Close()
	}
	srv.isRunning = false
}

func serialize(cmd gstypes.Command) ([]byte, error){
	var bytes []byte
	var err error

	packetType, err := helpers.ParseValueToBytes(cmd.PacketType)
	data := cmd.Data
	bytes = append(packetType,data...)

	return bytes, err
}

func CreateNewUDPServers (channel chan<- gstypes.PacketStoreElement, loggerChannel chan <- gstypes.PacketStoreElement) []*UDPListenerServer{
	//calculate the amount of ports that we'll have to listen to
	amountNodes := len(constants.Hosts)
	//create an array that will hold the port numbers
	nodesPorts := make([]int, amountNodes)
	//create an array that will keep the servers
	serversArray := make([]*UDPListenerServer,amountNodes)
	//populate the nodeports array with the port numbers
	mapIndex := 0
	for k := range constants.Hosts {
		nodesPorts[mapIndex] = k
		mapIndex++
	}
	//create and open all the servers
	for idx:= 0; idx < amountNodes; idx++ {
		srv := &UDPListenerServer{
			isRunning: false,
			doRun: false,
			ch: channel,
			loggerChan:loggerChannel}
		err := srv.open(nodesPorts[idx])
		if err == nil {
			serversArray[idx] = srv
		}else{
			log.Fatalf("Unable to open port %d: %v", nodesPorts[idx],err)
		}
	}
	return serversArray
}

func CreateNewUDPCommandServer(channel <-chan gstypes.Command) *UDPBroadcasterServer{
	srv := &UDPBroadcasterServer{
		ch: channel,
		isRunning:false,
		doRun:false}
	return srv
}
