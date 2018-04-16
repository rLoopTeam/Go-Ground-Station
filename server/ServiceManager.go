package server

import (
	"rloop/Go-Ground-Station/logging"
	"rloop/Go-Ground-Station/datastore"
	"rloop/Go-Ground-Station/proto"
	"google.golang.org/grpc"
	"net"
	"fmt"
)

type ServiceManager struct {
	serviceChan        <-chan *proto.ServerControl
	isRunning          bool
	doRun              bool
	dataStoreManager   *datastore.DataStoreManager
	gRPCServer         *grpc.Server
	udpListenerServers []*UDPListenerServer
	udpBroadcaster     *UDPBroadcasterServer
	gsLogger           *logging.Gslogger
	grpcConn           net.Listener
}

func (manager *ServiceManager) SetGsLogger(logger *logging.Gslogger) {
	manager.gsLogger = logger
}

func (manager *ServiceManager) SetGrpcServer(grpcServer *grpc.Server, conn net.Listener) {
	manager.gRPCServer = grpcServer
	manager.grpcConn = conn
}

func (manager *ServiceManager) SetDatastoreManager(datastoreManager *datastore.DataStoreManager) {
	manager.dataStoreManager = datastoreManager
}

func (manager *ServiceManager) SetUDPListenerServers(udpListeners []*UDPListenerServer) {
	manager.udpListenerServers = udpListeners
}

func (manager *ServiceManager) SetUDPBroadcaster(broadcaster *UDPBroadcasterServer) {
	manager.udpBroadcaster = broadcaster
}

func (manager *ServiceManager) RunAll() {
	manager.StartLogger()
	manager.StartDatastoreManager()
	manager.StartUDPListeners()
	manager.StartGrpcServer()
	manager.StartBroadcaster()
	go manager.Run()
}
func (manager *ServiceManager) StopAll() {
	manager.StopLogger()
	manager.StopDatastoreManager()
	manager.StopUDPListeners()
	manager.StopGrpcServer()
	manager.StopBroadcaster()
}

func (manager *ServiceManager) StartUDPListeners() {
	for _, srv := range manager.udpListenerServers {
		go srv.Run()
	}
}
func (manager *ServiceManager) StopUDPListeners() {}

func (manager *ServiceManager) StartDatastoreManager() { manager.dataStoreManager.Start() }
func (manager *ServiceManager) StopDatastoreManager() {	manager.dataStoreManager.Stop() }

func (manager *ServiceManager) StartBroadcaster() {	go manager.udpBroadcaster.Run() }
func (manager *ServiceManager) StopBroadcaster() { manager.udpBroadcaster.Stop() }

func (manager *ServiceManager) StartGrpcServer() {
	if manager.gRPCServer != nil {
		manager.gRPCServer.Serve(manager.grpcConn)
	} else {
		fmt.Println("cannot start grpc service, server is not set")
	}
}
func (manager *ServiceManager) StopGrpcServer() {}

func (manager *ServiceManager) StartLogger() {
	if !manager.gsLogger.IsRunning {
		go manager.gsLogger.Start()
	}
}
func (manager *ServiceManager) StopLogger() { manager.gsLogger.Stop() }

func (manager *ServiceManager) Run() {
	manager.doRun = true
	for {
		if (!manager.doRun) {
			break
		}
		control := <-manager.serviceChan
		manager.executeControl(control)
	}
}

func (manager *ServiceManager) executeControl(control *proto.ServerControl) {
	switch control.Command {
	case proto.ServerControl_LogServiceStart:
		manager.StartLogger()
		break
	case proto.ServerControl_LogServiceStop:
		manager.StopLogger()
		break
	case proto.ServerControl_DataStoreManagerStart:
		manager.StartDatastoreManager()
		break
	case proto.ServerControl_DataStoreManagerStop:
		manager.StopDatastoreManager()
		break
	case proto.ServerControl_BroadcasterStart:
		manager.StartBroadcaster()
		break
	case proto.ServerControl_BroadcasterStop:
		manager.StopBroadcaster()
		break
	}
}

func (manager *ServiceManager) Stop() {
	manager.doRun = false
}

func NewServiceManager() (*ServiceManager, chan<- *proto.ServerControl) {
	srvcChan := make(chan *proto.ServerControl, 8)
	serviceManager := &ServiceManager{
		serviceChan: srvcChan,
		isRunning:   false,
		doRun:       false}
	return serviceManager, srvcChan
}
