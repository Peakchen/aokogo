package dbStatistics

/*
	purpose: statistics db write count for check db operate frequency.
	date: 20200114 14:30
*/

import (
	"commom/stacktrace"
	"os"
	"common/utls"
	"encoding/gob"
)

type TModel struct {
	model 	string
	cnt 	int
}

type TDBStatistics struct {
	filehandle 	*os.File
	buff 		bytes.Buffer

}

const (
	timeDate string = "2006-01-02"
	cstSaveLogTickers = 60 //s
)

var (
	_dbStatistics *TDBStatistics
)

func InitDBStatistics(){
	_dbStatistics = &TDBStatistics{}
	_dbStatistics.Init()
	go _dbStatistics.loop()
}

func (this *TDBStatistics) Init(){
	exename := utls.GetExeFileName()
	fileName := fmt.Sprintf("./DBStatisticsLog/%v_DBStatistics_%v.log", exename, time.Now().Local().Format(timeDate))
	filehandle, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		Log.Error("can not create db statistics log file.")
		return
	}

	this.filehandle = filehandle
}

func (this *TDBStatistics) loop(){
	tick := time.NewTicker(time.Duration(cstSaveLogTickers)*time.Second)
	for {
		select {
		case <-tick.C:
			if this.buff.Len() == 0 {
				continue
			}
			this.filehandle.WriteString(this.buff.String())
		}
	}
}

func (this *TDBStatistics) Exit(){
	this.filehandle.Sync()
	this.filehandle.Close()
}

func (this *TDBStatistics) Update(content string) {
	if err := gob.NewEncoder(&(this.buff)).Encode(content); err != nil {
		Log.Error("cover to buffer fail, err: ", err)
    }
}

func DBUpdateStatistics(identify, model string){
	var content string
	content = "identify: " + identify + "\r\n"
	content += "model: " + model + "\r\n"
	content += "stack log: \r\n" + stacktrace.NormalStackLog() + "\r\n"
	_dbStatistics.Update(content)
}