package logging

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"rloop/Go-Ground-Station/gstypes"
	"strconv"
	"sync"
	"time"
)

type Gslogger struct {
	ticker         *time.Ticker
	DataChan       <-chan gstypes.PacketStoreElement
	signalChan     chan bool
	doRunMutex     sync.RWMutex
	doRun          bool
	isRunningMutex sync.RWMutex
	isRunning      bool
}

func (gslogger *Gslogger) Start() {
	gslogger.doRun = true
	gslogger.run()
}

func (gslogger *Gslogger) Stop() {
	if gslogger.isRunning {
		gslogger.signalChan <- true
		gslogger.doRun = false
	}
}

func (gslogger *Gslogger) run() {
	currTime := time.Now()
	day := strconv.Itoa(currTime.Day())
	month := currTime.Month().String()
	year := strconv.Itoa(currTime.Year())
	hour := strconv.Itoa(currTime.Hour())
	minute := strconv.Itoa(currTime.Minute())
	//seconds := strconv.Itoa(currTime.Second())

	fileName := fmt.Sprintf("gslog_%s-%s-%s_%s%s.csv", day, month, year, hour, minute)
	headers := []string{"rxtime", "port", "nodename", "packtype", "packetname", "prefix", "parametername", "units", "value"}

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Error Creating log file:", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Write(headers)
	gslogger.isRunning = true
MainLoop:
	for {
		select {
		case data := <-gslogger.DataChan:
			gslogger.write(data, writer)
		case <-gslogger.signalChan:
			break MainLoop
		}
	}
	gslogger.isRunning = false
}

func (gsLogger *Gslogger) write(data gstypes.PacketStoreElement, writer *csv.Writer) {
	rxtime := strconv.FormatInt(data.RxTime, 10)
	port := strconv.FormatInt(int64(data.Port), 10)
	packetType := strconv.FormatInt(int64(data.PacketType), 10)
	packetName := data.PacketName
	parameterPrefix := data.ParameterPrefix
	parameters := data.Parameters
	for _, param := range parameters {
		line := []string{rxtime, port, packetType, packetName, parameterPrefix, param.ParameterName, param.Units, param.Data.AsString()}
		err := writer.Write(line)
		if err != nil {
			log.Fatal("Error Writing log file:", err)
		}
	}
	select {
	case <-gsLogger.ticker.C:
		writer.Flush()
	default:
	}
}

func (gsLogger *Gslogger) GetStatus() (bool, bool) {
	defer func() {
		gsLogger.isRunningMutex.RUnlock()
		gsLogger.doRunMutex.RUnlock()
	}()
	gsLogger.isRunningMutex.RLock()
	gsLogger.doRunMutex.RLock()
	return gsLogger.isRunning, gsLogger.doRun
}

func New() (*Gslogger, chan<- gstypes.PacketStoreElement) {
	dataChan := make(chan gstypes.PacketStoreElement, 512)
	signalChan := make(chan bool)
	gslogger := &Gslogger{
		DataChan:   dataChan,
		doRun:      false,
		isRunning:  false,
		signalChan: signalChan,
		ticker:     time.NewTicker(10 * time.Second)}
	return gslogger, dataChan
}
