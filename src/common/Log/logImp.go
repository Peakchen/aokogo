package Log

import (
	"fmt"
	"os"
	"time"
	"runtime/debug"
)

type TLogData struct{
	filename string
	filehandle *os.File
}

var (
	GLogObj *TLogData
)

func NewLog(filename string){
	GLogObj := &TLogData{}
	RealFileName := fmt.Sprintf("%v_%v", filename, time.Now().Unix())
	filehandler,err := os.OpenFile(RealFileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil{
		return
	}

	GLogObj.filehandle = filehandler
	GLogObj.filename = RealFileName
}

func Error(format string, args ...interface{}){
	if len(format) == 0{
		return
	} 

	WriteLog("[Error]", format, args)
}

func Info(format string, args ...interface{}){
	if len(format) == 0{
		return
	} 

	WriteLog("[Info]", format, args)
}

func Panic(format string, args ...interface{}){
	logStr := fmt.Sprintf("%v", debug.Stack())
	GLogObj.filehandle.WriteString(logStr)
	EndLog()
	debug.PrintStack()
}

func WriteLog(title, format string, args ...interface{}){
	logStr := fmt.Sprintf(title+format, args...)
	if GLogObj == nil {
		panic("log instance is nil.")
	}

	_, err := GLogObj.filehandle.WriteString(logStr)
	if err != nil{
		return
	}
}

func EndLog(){
	if GLogObj.filehandle == nil {
		return
	}

	GLogObj.filehandle.Close()
}
