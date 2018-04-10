package logging

import (
	"rloop/Go-Ground-Station/gstypes"
	"os"
	"encoding/csv"
	"log"
	//"time"
	"strconv"
	"time"
	"fmt"
)

type Gslogger struct {
	DataChan <- chan gstypes.PacketStoreElement
}

func (gslogger *Gslogger) Run (){
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
	gslogger.checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Write(headers)

	for data := range gslogger.DataChan {
		rxtime := strconv.FormatInt(data.RxTime,10)
		port := strconv.FormatInt(int64(data.Port),10)
		packetType := strconv.FormatInt(int64(data.PacketType),10)
		packetName := data.PacketName
		parameterPrefix := data.ParameterPrefix
		parameters := data.Parameters
		for _,param := range parameters{
			line := []string{rxtime,port,packetType,packetName,parameterPrefix,param.ParameterName,param.Units,param.Data.AsString()}
			err := writer.Write(line)
			gslogger.checkError("Cannot write to file", err)
			writer.Flush()
		}
	}
}

//https://golangcode.com/write-data-to-a-csv-file/
func (Gslogger)checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}