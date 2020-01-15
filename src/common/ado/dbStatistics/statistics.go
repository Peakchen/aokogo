package dbStatistics

/*
	purpose: statistics db write count for check db operate frequency.
	date: 20200114 14:30
*/

import (
	"common/stacktrace"
	"os"
	"common/utls"
	"encoding/gob"
	"bytes"
	"fmt"
	"common/Log"
	"time"
	"common/public"
)

type TModel struct {
	model 	string
	cnt 	int
}

type TDBStatistics struct {
	filehandle 	*os.File
	chbuff 		chan bytes.Buffer

}

const (
	cstSaveLogTickers = 60 //s
	cstStatisticsLog = "DBStatisticsLog"
)

var (
	_dbStatistics *TDBStatistics
)

func InitDBStatistics(){
	_dbStatistics = &TDBStatistics{
		filehandle: nil,
		chbuff: make(chan bytes.Buffer),
	}
	_dbStatistics.Init()
	go _dbStatistics.loop()
}

func (this *TDBStatistics) Init(){
	exename := utls.GetExeFileName()
	_, err := os.Stat(cstStatisticsLog)
	if err != nil {
		Log.Error("err: ", err)
		return
	}

	if os.IsNotExist(err){
		err := os.Mkdir("DBStatisticsLog", 0644)
		if err != nil {
			Log.Error("err: ", err)
			return
		}
	}

	fileName := fmt.Sprintf("./DBStatisticsLog/%v_DBStatistics_%v.log", exename, time.Now().Local().Format(public.CstTimeDate))
	filehandle, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		Log.Error("can not create db statistics log file, err：", err)
		return
	}

	this.filehandle = filehandle
}

func (this *TDBStatistics) loop(){
	for {
		select {
		case data := <-this.chbuff:
			this.filehandle.WriteString(data.String())
		}
	}
}

func (this *TDBStatistics) exit(){
	this.filehandle.Sync()
	this.filehandle.Close()
}

func (this *TDBStatistics) Update(content string) {
	var buff bytes.Buffer
	if err := gob.NewEncoder(&buff).Encode(content); err != nil {
		Log.Error("cover to buffer fail, err: ", err)
		return
	}
	
	this.chbuff <- buff
}

func DBStatistics(identify, model string){
	content := "identify: " + identify + "\r\n"
	content += "model: " + model + "\r\n"
	content += "time: " + time.Now().Local().Format(public.CstTimeFmt) + "\r\n"
	content += "stack log: \r\n" + stacktrace.NormalStackLog() + "\r\n"
	dst := utls.GBKToUTF8(content)
	_dbStatistics.Update(dst)
}

func DBStatisticsStop(){
	_dbStatistics.exit()
}