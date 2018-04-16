package logging

import (
	"rloop/Go-Ground-Station/gstypes"
	"os"
	"encoding/csv"
	"log"
	"strconv"
	"time"
	"fmt"
)

type Gslogger struct {
	DataChan <- chan gstypes.PacketStoreElement
	doRun bool
	IsRunning bool
}

func (gslogger *Gslogger) Start(){
	gslogger.doRun = true
	gslogger.run()
}

func (gslogger *Gslogger) Stop(){
	gslogger.doRun = false
}

func (gslogger *Gslogger) run (){
	currTime := time.Now()
	day := strconv.Itoa(currTime.Day())
	month := currTime.Month().String()
	year := strconv.Itoa(currTime.Year())
	hour := strconv.Itoa(currTime.Hour())
	minute := strconv.Itoa(currTime.Minute())
	//seconds := strconv.Itoa(currTime.Second())

	fileName := fmt.Sprintf("gslog_%s-%s-%s_%s%s.csv",day,month,year,hour,minute)
	headers := []string{"rxtime","port","nodename","packtype","packetname","prefix","parametername","units","value"}
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Error Creating log file:", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Write(headers)
	gslogger.IsRunning = true
	for data := range gslogger.DataChan {
		if(!gslogger.doRun){break}

		rxtime := strconv.FormatInt(data.RxTime,10)
		port := strconv.FormatInt(int64(data.Port),10)
		packetType := strconv.FormatInt(int64(data.PacketType),10)
		packetName := data.PacketName
		parameterPrefix := data.ParameterPrefix
		parameters := data.Parameters
		for _,param := range parameters{
			line := []string{rxtime,port,packetType,packetName,parameterPrefix,param.ParameterName,param.Units,param.Data.AsString()}
			err := writer.Write(line)
			if err != nil {
				log.Fatal("Error Writing log file:", err)
			}
			writer.Flush()
		}
	}
	gslogger.IsRunning = false
}

func New() (*Gslogger, chan <- gstypes.PacketStoreElement){
	dataChan := make(chan gstypes.PacketStoreElement, 256)
	gslogger := &Gslogger{
		DataChan:dataChan,
		doRun:false,
		IsRunning:false}
	return gslogger,dataChan
}