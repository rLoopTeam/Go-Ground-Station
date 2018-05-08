package server

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"rloop/Go-Ground-Station/datastore"
	"rloop/Go-Ground-Station/gstypes"
	"rloop/Go-Ground-Station/logging"
	"sync"
	"time"
)

type ServiceManager struct {
	serviceChan        <-chan gstypes.ServerControlWithTimeout
	isRunning          bool
	doRun              bool
	dataStoreManager   *datastore.DataStoreManager
	gRPCServer         *grpc.Server
	udpListenerServers []*UDPListenerServer
	udpBroadcaster     *UDPBroadcasterServer
	gsLogger           *logging.Gslogger
	simController      *SimController
	grpcConn           net.Listener
	serviceCheckTicker *time.Ticker
	StatusMutex        sync.RWMutex
	Status             gstypes.ServiceStatus
}

func (manager *ServiceManager) GetStatus() gstypes.ServiceStatus {
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
	//fmt.Println("startlogger")
	//manager.StartLogger()
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
func (manager *ServiceManager) StopDatastoreManager()  { manager.dataStoreManager.Stop() }

func (manager *ServiceManager) StartBroadcaster() { go manager.udpBroadcaster.Run() }
func (manager *ServiceManager) StopBroadcaster()  { manager.udpBroadcaster.Stop() }

func (manager *ServiceManager) StartGrpcServer() {
	if manager.gRPCServer != nil {
		go manager.gRPCServer.Serve(manager.grpcConn)
	} else {
		fmt.Println("cannot start grpc service, server is not set")
	}
}

func (manager *ServiceManager) StopGrpcServer() { manager.gRPCServer.Stop() }

func (manager *ServiceManager) StartLogger() {
	isRunning, _ := manager.gsLogger.GetStatus()
	if !isRunning {
		go manager.gsLogger.Start()
	}
}
func (manager *ServiceManager) StopLogger() { manager.gsLogger.Stop() }

func (manager *ServiceManager) StartSimController() { manager.simController.Run() }

func (manager *ServiceManager) StopSimController() { manager.simController.Stop() }

func (manager *ServiceManager) Run() {
	manager.doRun = true
	for {
		if !manager.doRun {
			break
		}
		control := <-manager.serviceChan
		fmt.Printf("Service Control requested: %v \n", control.Control)
		manager.executeControl(control.Control)
	}
}

func (manager *ServiceManager) checkStatus() {
	//fmt.Println("statuschecker running")
	for t := range manager.serviceCheckTicker.C {
		manager.StatusMutex.Lock()
		//fmt.Printf("status on %d \n", t.Unix())
		//sstring := fmt.Sprintf("GSLogger is running: %t \n" +
		//	"Broadcaster is running: %t \n" +
		//		"DataStoreManager is running: %t \n" +
		//			"Grpc is running: %t \n", manager.Status.GSLoggerRunning, manager.Status.BroadcasterRunning, manager.Status.DataStoreManagerRunning, manager.Status.GRPCServerRunning )
		//fmt.Printf("%s \n",sstring)
		manager.Status.LastUpdated = t.Unix()
		manager.Status.BroadcasterRunning, _ = manager.udpBroadcaster.GetStatus()
		manager.Status.DataStoreManagerRunning, _ = manager.dataStoreManager.GetStatus()
		manager.Status.GSLoggerRunning, _ = manager.gsLogger.GetStatus()
		for _, srv := range manager.udpListenerServers {
			manager.Status.PortsListening[srv.ServerPort] = srv.IsRunning
		}
		manager.StatusMutex.Unlock()
	}
}

func (manager *ServiceManager) executeControl(control gstypes.ServerControl_CommandEnum) bool {
	var result bool = false
	switch control {
	case gstypes.ServerControl_LogServiceStart:
		manager.StartLogger()
		break
	case gstypes.ServerControl_LogServiceStop:
		manager.StopLogger()
		break
	case gstypes.ServerControl_DataStoreManagerStart:
		manager.StartDatastoreManager()
		break
	case gstypes.ServerControl_DataStoreManagerStop:
		manager.StopDatastoreManager()
		break
	case gstypes.ServerControl_BroadcasterStart:
		manager.StartBroadcaster()
		break
	case gstypes.ServerControl_BroadcasterStop:
		manager.StopBroadcaster()
		break
	}
	return result
}

func (manager *ServiceManager) Stop() {
	manager.doRun = false
}

func NewServiceManager() (*ServiceManager, chan<- gstypes.ServerControlWithTimeout) {
	srvcChan := make(chan gstypes.ServerControlWithTimeout, 8)
	serviceManager := &ServiceManager{
		serviceChan:        srvcChan,
		isRunning:          false,
		doRun:              false,
		serviceCheckTicker: time.NewTicker(time.Second * 5),
		StatusMutex:        sync.RWMutex{},
		Status:             gstypes.NewServiceStatus()}

	return serviceManager, srvcChan
}
