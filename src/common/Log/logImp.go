package Log

import (
	"common/utls"
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"
)

type TAokoLog struct {
	filename   string
	filehandle *os.File
	cancle     context.CancelFunc
	ctx        context.Context
	wg         sync.WaitGroup
	filesize   uint64
	logNum     uint64
	data       chan string
}

const (
	EnAKLogFileMaxLimix = 500 * 1024 * 1024
	EnLogDataChanMax    = 1024
)

var (
	aokoLog *TAokoLog
)
var exitchan = make(chan os.Signal, 1)

func NewLog() {
	aokoLog = &TAokoLog{}
	initLogFile()
	//aokoLog.ctx, aokoLog.cancle = context.WithCancel(context.Background())
}

func initLogFile() {
	filename := utls.GetExeFileName()
	RealFileName := fmt.Sprintf("./log/%v_%v.log", filename, time.Now().Unix())
	filehandler, err := os.OpenFile(RealFileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}

	aokoLog.filehandle = filehandler
	aokoLog.filename = RealFileName
	aokoLog.data = make(chan string, EnLogDataChanMax)

}

func Run(sw *sync.WaitGroup, ctx context.Context) {
	signal.Notify(exitchan, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGSEGV)
	aokoLog.wg.Add(1)
	go aokoLog.loop(sw, ctx)
}

func Error(format string, args ...interface{}) {
	WriteLog("[Error]\t\t\t", format, args)
}

func Info(format string, args ...interface{}) {
	WriteLog("[Info]\t\t\t", format, args)
}

func Panic(format string, args ...interface{}) {
	logStr := fmt.Sprintf("%v", debug.Stack())
	aokoLog.filehandle.WriteString(logStr)
	aokoLog.endLog()
	debug.PrintStack()
}

func WriteLog(title, format string, args ...interface{}) {
	if aokoLog == nil {
		Panic("log instance is nil.")
	}

	var (
		logStr string
	)

	if len(format) == 0 && len(args) > 0 {
		logStr += title
		for i, data := range args {
			if i+1 <= len(args) {
				logStr += fmt.Sprintf("%v", data)
			}
		}
		logStr += "\n"
	} else if len(args) == 0 && len(format) > 0 {
		logStr = fmt.Sprintf(title + format)
	} else if len(format) > 0 && len(args) > 0 {
		logStr = fmt.Sprintf(title+format, args...)
	}

	if len(logStr) == 0 {
		return
	}

	if aokoLog.filesize >= EnAKLogFileMaxLimix {
		FmtPrintf("log file: %v over max limix.", aokoLog.filename)
		initLogFile()
	}

	aokoLog.filesize += uint64(len(logStr))
	aokoLog.logNum++
	aokoLog.data <- logStr
}

func (this *TAokoLog) endLog() {
	if this.filehandle != nil {
		this.filehandle.Sync()
		this.filehandle.Close()
	}
}

func (this *TAokoLog) exit(sw *sync.WaitGroup) {
	fmt.Println("log exit: ", <-this.data, aokoLog.filesize, aokoLog.logNum)
	this.flush()
	this.endLog()
	close(this.data)
	sw.Wait()
}

func (this *TAokoLog) loop(sw *sync.WaitGroup, ctx context.Context) {
	defer this.exit(sw)
	tick := time.NewTicker(time.Duration(30 * time.Second))
	for {
		if s, ok := <-exitchan; ok {
			tick.Stop()
			this.exit(sw)
			time.Sleep(time.Duration(3) * time.Second)
			fmt.Println("Got signal:", s)
			return
		}
		select {
		case <-ctx.Done():
			tick.Stop()
			return
		case <-tick.C:
			go this.flush()
		default:

		}
	}
}

func (this *TAokoLog) flush() {
	for {
		select {
		case val, ok := <-this.data:
			if ok {
				_, err := aokoLog.filehandle.WriteString(val)
				if err != nil {
					return
				}
			}
		}
	}

}
