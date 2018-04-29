package server

import (
	"rloop/Go-Ground-Station/logging"
	"rloop/Go-Ground-Station/datastore"
	"rloop/Go-Ground-Station/proto"
	"google.golang.org/grpc"
	"net"
	"fmt"
	"sync"
	"time"
	"rloop/Go-Ground-Station/gstypes"
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
	simController *SimController
	grpcConn           net.Listener
	serviceCheckTicker *time.Ticker
	StatusMutex sync.RWMutex
	Status gstypes.ServiceStatus
}

func (manager *ServiceManager) GetStatus() gstypes.ServiceStatus{
	manager.StatusMutex.RLock()
	defer manager.StatusMutex.RUnlock()
	return manager.Status
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

func (manager *ServiceManager) SetSimController(simController *SimController) {
	manager.simController = simController
}

func (manager *ServiceManager) RunAll() {
	fmt.Println("startlogger")
	manager.StartLogger()
	fmt.Println("startdatastore")
	manager.StartDatastoreManager()
	fmt.Println("startudplisteners")
	manager.StartUDPListeners()
	fmt.Println("startgrpc")
	manager.StartGrpcServer()
	fmt.Println("startbroadcaster")
	manager.StartBroadcaster()
	fmt.Println("starrun")
	go manager.Run()
	fmt.Println("startcheckstatus")
	go manager.checkStatus()
	fmt.Println("startSimController")
	go manager.StartSimController()
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
		go manager.gRPCServer.Serve(manager.grpcConn)
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

func (manager *ServiceManager) StartSimController() {manager.simController.Run()}

func (manager *ServiceManager) StopSimController() {manager.simController.Stop()}

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

func (manager *ServiceManager) checkStatus() {
	//fmt.Println("statuschecker running")
	for t := range manager.serviceCheckTicker.C {
		manager.StatusMutex.Lock()
		fmt.Printf("status on %d \n",t.Unix())
		//fmt.Printf("%v \n",manager.Status)
		manager.Status.BroadcasterRunning = manager.udpBroadcaster.isRunning
		manager.Status.DataStoreManagerRunning = manager.dataStoreManager.IsRunning
		//manager.Status.GRPCServerRunning = manager.gRPCServer.IsRunning
		manager.Status.GSLoggerRunning = manager.gsLogger.IsRunning
		for _, srv := range manager.udpListenerServers {
			manager.Status.PortsListening[srv.ServerPort] = srv.IsRunning
		}
		manager.StatusMutex.Unlock()
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
		doRun:       false,
		serviceCheckTicker: time.NewTicker(time.Second*5),
		StatusMutex: sync.RWMutex{},
		Status: gstypes.NewServiceStatus()}

	return serviceManager, srvcChan
}
