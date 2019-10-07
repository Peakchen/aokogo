package Log

import (
	"common/utls"
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
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
	sw         sync.WaitGroup
	FileNo     uint32
}

const (
	EnAKLogFileMaxLimix = 500 * 1024 * 1024
	EnLogDataChanMax    = 1024
)

const (
	EnLogType_Info  string = "logInfo"
	EnLogType_Error string = "logError"
	EnLogType_Fail  string = "logFail"
	EnLogType_Debug string = "logDebug"
)

var (
	aokoLog map[string]*TAokoLog
)
var exitchan = make(chan os.Signal, 1)

func init() {
	aokoLog = map[string]*TAokoLog{}
}

func checkNewLog(logtype string) (logobj *TAokoLog) {
	var (
		ok bool
	)
	logobj, ok = aokoLog[logtype]
	if !ok {
		aokoLog[logtype] = &TAokoLog{
			FileNo: 1,
		}
		initLogFile(logtype, aokoLog[logtype])
		run(aokoLog[logtype])
		logobj = aokoLog[logtype]
	}
	return
}

func initLogFile(logtype string, aokoLog *TAokoLog) {
	var (
		RealFileName string
		PathDir      string = logtype
	)

	filename := utls.GetExeFileName()
	switch logtype {
	case EnLogType_Info:
		RealFileName = fmt.Sprintf("./logInfo/%v_Info_No%v_%v.log", filename, aokoLog.FileNo, time.Now().Local().Format(timeDate))
	case EnLogType_Error:
		RealFileName = fmt.Sprintf("./logError/%v_Error_No%v_%v.log", filename, aokoLog.FileNo, time.Now().Local().Format(timeDate))
	case EnLogType_Fail:
		RealFileName = fmt.Sprintf("./logFail/%v_Fail_No%v_%v.log", filename, aokoLog.FileNo, time.Now().Local().Format(timeDate))
	case EnLogType_Debug:
		RealFileName = fmt.Sprintf("./logDebug/%v_Debug_No%v_%v.log", filename, aokoLog.FileNo, time.Now().Local().Format(timeDate))
	default:

	}

	err := os.Remove(PathDir)
	if err != nil {
		if reflect.TypeOf(err) != reflect.TypeOf(&os.PathError{}) {
			fmt.Println("err dir type: ", reflect.TypeOf(err))
			return
		}
		perror := err.(*os.PathError)
		if perror.Err != syscall.ENOENT &&
			perror.Err != syscall.ERROR_DIR_NOT_EMPTY {
			fmt.Printf("Remove log dir fail, dir: %v, errcode: %v, err: %v.\n", PathDir, perror.Err, err.Error())
			return
		}
	}

	err = os.Mkdir(PathDir, os.ModePerm)
	if err != nil {
		if reflect.TypeOf(err) != reflect.TypeOf(&os.PathError{}) {
			fmt.Println("err dir type: ", reflect.TypeOf(err))
			return
		}
		perror := err.(*os.PathError)
		if perror.Err != syscall.ERROR_ALREADY_EXISTS {
			fmt.Printf("log mkdir fail, dir: %v, errcode: %v, err: %v.\n", PathDir, perror.Err, err.Error())
			return
		}
	}

	filehandler, err := os.OpenFile(RealFileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}

	aokoLog.filehandle = filehandler
	aokoLog.filename = RealFileName
	aokoLog.data = make(chan string, EnLogDataChanMax)

}

func run(aokoLog *TAokoLog) {
	signal.Notify(exitchan, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGSEGV)
	aokoLog.ctx, aokoLog.cancle = context.WithCancel(context.Background())
	aokoLog.wg.Add(1)
	go aokoLog.loop(aokoLog)
}

func Error(args ...interface{}) {
	format := ""
	WriteLog(EnLogType_Error, "[Error]\t\t\t", format, args)
}

func Info(format string, args ...interface{}) {
	WriteLog(EnLogType_Info, "[Info]\t\t\t", format, args)
}

func Fail(args ...interface{}) {
	format := ""
	WriteLog(EnLogType_Fail, "[Fail]\t\t\t", format, args)
}

func Debug(format string, args ...interface{}) {
	WriteLog(EnLogType_Debug, "[Debug]\t\t\t", format, args)
}

func Panic() {
	aokoLog := checkNewLog(EnLogType_Fail)
	if aokoLog != nil {
		debug.PrintStack()
		buf := debug.Stack()
		aokoLog.filehandle.WriteString(string(buf[:]))
		aokoLog.endLog()
		//close(aokoLog.data)
	}
}

func WriteLog(logtype, title, format string, args ...interface{}) {
	aokoLog := checkNewLog(logtype)
	if aokoLog == nil {
		Panic()
		return
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
		aokoLog.FileNo++
		initLogFile(logtype, aokoLog)
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

func (this *TAokoLog) exit(aokoLog *TAokoLog) {
	fmt.Println("log exit: ", <-this.data, aokoLog.filesize, aokoLog.logNum)
	this.flush(aokoLog)
	this.endLog()
	close(this.data)
	this.sw.Wait()
}

func (this *TAokoLog) loop(aokoLog *TAokoLog) {
	defer this.sw.Done()

	tick := time.NewTicker(time.Duration(30 * time.Second))
	for {
		if s, ok := <-exitchan; ok {
			tick.Stop()
			this.exit(aokoLog)
			time.Sleep(time.Duration(3) * time.Second)
			fmt.Println("Got signal:", s)
			return
		}

		select {
		case <-this.ctx.Done():
			tick.Stop()
			return
		case <-tick.C:
			go this.flush(aokoLog)
		default:

		}
	}
}

func (this *TAokoLog) flush(aokoLog *TAokoLog) {
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
