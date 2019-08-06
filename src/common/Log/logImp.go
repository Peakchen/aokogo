package Log

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type TAokoLog struct {
	filename   string
	filehandle *os.File
	cancle     context.CancelFunc
	ctx        context.Context
	wg         sync.WaitGroup
}

var (
	GLogObj *TAokoLog
)

func NewLog(filename string) {
	GLogObj := &TAokoLog{}
	RealFileName := fmt.Sprintf("%v_%v", filename, time.Now().Unix())
	filehandler, err := os.OpenFile(RealFileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}

	GLogObj.filehandle = filehandler
	GLogObj.filename = RealFileName
	GLogObj.ctx, GLogObj.cancle = context.WithCancel(context.Background())
	GLogObj.wg.Add(1)
	go GLogObj.loop()

}

func Error(format string, args ...interface{}) {
	if len(format) == 0 {
		return
	}

	WriteLog("[Error]", format, args)
}

func Info(format string, args ...interface{}) {
	if len(format) == 0 {
		return
	}

	WriteLog("[Info]", format, args)
}

func Panic(format string, args ...interface{}) {
	logStr := fmt.Sprintf("%v", debug.Stack())
	GLogObj.filehandle.WriteString(logStr)
	EndLog()
	debug.PrintStack()
}

func WriteLog(title, format string, args ...interface{}) {
	logStr := fmt.Sprintf(title+format, args...)
	if GLogObj == nil {
		panic("log instance is nil.")
	}

	_, err := GLogObj.filehandle.WriteString(logStr)
	if err != nil {
		return
	}
}

func EndLog() {
	if GLogObj.filehandle == nil {
		return
	}
	GLogObj.filehandle.Sync()
	GLogObj.filehandle.Close()
}

func (this *TAokoLog) exit() {
	this.wg.Wait()
	if this.filehandle != nil {
		this.filehandle.Close()
	}
}

func (this *TAokoLog) loop() {
	defer this.wg.Done()
	for {
		select {
		case <-this.ctx.Done():
			this.exit()
			return
		default:

		}
	}
}
