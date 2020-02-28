package Log

// add by stefan

import (
	"common/aktime"
	"common/public"
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
		go run(aokoLog[logtype])
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
		RealFileName = fmt.Sprintf("./logInfo/%v_Info_No%v_%v.log", filename, aokoLog.FileNo, aktime.Now().Local().Format(public.CstTimeDate))
	case EnLogType_Error:
		RealFileName = fmt.Sprintf("./logError/%v_Error_No%v_%v.log", filename, aokoLog.FileNo, aktime.Now().Local().Format(public.CstTimeDate))
	case EnLogType_Fail:
		RealFileName = fmt.Sprintf("./logFail/%v_Fail_No%v_%v.log", filename, aokoLog.FileNo, aktime.Now().Local().Format(public.CstTimeDate))
	case EnLogType_Debug:
		RealFileName = fmt.Sprintf("./logDebug/%v_Debug_No%v_%v.log", filename, aokoLog.FileNo, aktime.Now().Local().Format(public.CstTimeDate))
	default:

	}

	exepath := utls.GetExeFilePath()
	filepath := exepath + "/" + PathDir
	exist, err := utls.IsPathExisted(filepath)
	if err != nil {
		panic("check path exist err: " + err.Error())
		return
	}

	if false == exist {
		err = os.Mkdir(filepath, os.ModePerm)
		if err != nil {
			panic("log mkdir fail, err: " + err.Error())
			return
		}
	}

	filehandler, err := os.OpenFile(RealFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
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
	go aokoLog.loop()
	aokoLog.wg.Wait()
}

func Error(args ...interface{}) {
	timeFormat := aktime.Now().Local().Format(public.CstTimeFmt)
	format := "time: " + timeFormat + ", Info: %v."
	WriteLog(EnLogType_Error, "[Error]\t", format, args)
}

func ErrorIDCard(identify string, args ...interface{}) {
	format := fmt.Sprintf("identify: %v, %v.", identify, args)
	timeFormat := aktime.Now().Local().Format(public.CstTimeFmt)
	WriteLog(EnLogType_Error, "[Error]\t", timeFormat, format)
}

func ErrorModule(data public.IDBCache, args ...interface{}) {
	format := fmt.Sprintf("main: %v, sub: %v, identify: %v, %v.", data.MainModel(), data.SubModel(), data.Identify(), args)
	timeFormat := aktime.Now().Local().Format(public.CstTimeFmt)
	WriteLog(EnLogType_Error, "[Error]\t", timeFormat, format)
}

func Info(format string, args ...interface{}) {
	timeFormat := aktime.Now().Local().Format(public.CstTimeFmt)
	WriteLog(EnLogType_Info, "[Info]\t", timeFormat+format, args)
}

func Fail(args ...interface{}) {
	format := aktime.Now().Local().Format(public.CstTimeFmt)
	WriteLog(EnLogType_Fail, "[Fail]\t", format, args)
}

func Debug(format string, args ...interface{}) {
	WriteLog(EnLogType_Debug, "[Debug]\t", format, args)
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

	/*
		print(a,b,c...)
	*/
	if len(format) == 0 && len(args) > 0 {
		logStr += fmt.Sprintf(title + format)
		for i, data := range args {
			if i+1 <= len(args) {
				logStr += fmt.Sprintf("%v", data)
			}
		}
		logStr += "\n"
	} else if len(args) == 0 && len(format) > 0 { //print("aaa,bbb,ccc.")
		logStr = fmt.Sprintf(title + format)
	} else if len(format) > 0 && len(args) > 0 { //print("a: %v, b: %v.",a,b)
		logStr = fmt.Sprintf(title+format, args...)
		logStr += "\n"
	}

	if len(logStr) == 0 {
		return
	}

	if aokoLog.filesize >= EnAKLogFileMaxLimix {
		FmtPrintf("log file: %v over max limix.", aokoLog.filename)
		aokoLog.FileNo++
		initLogFile(logtype, aokoLog)
		aokoLog.filesize = 0
	}

	aokoLog.filesize += uint64(len(logStr))
	aokoLog.logNum++
	aokoLog.data <- logStr

	if aokoLog.logNum%EnLogDataChanMax == 0 {
		aokoLog.flush()
		aokoLog.data = make(chan string, EnLogDataChanMax)
	}
}

func (this *TAokoLog) endLog() {
	if this.filehandle != nil {
		this.filehandle.Sync()
		this.filehandle.Close()
	}
}

func (this *TAokoLog) exit() {
	fmt.Println("log exit: ", <-this.data, this.filesize, this.logNum)
	this.flush()
	this.endLog()
	close(this.data)
	this.sw.Wait()
}

func (this *TAokoLog) loop() {
	defer func() {
		this.exit()
		this.sw.Done()
		time.Sleep(time.Duration(3) * time.Second)
	}()

	tick := time.NewTicker(time.Duration(10 * time.Second))
	for {
		select {
		case <-this.ctx.Done():
			tick.Stop()
			return
		case log, ok := <-this.data:
			if !ok {
				continue
			}
			this.writelog(log)
		case s, ok := <-exitchan:
			if !ok {
				continue
			}
			fmt.Println("Got signal:", s)
			return
		case <-tick.C:
			this.flush()
		}
	}
}

func (this *TAokoLog) writelog(src string) {
	_, err := this.filehandle.WriteString(src)
	if err != nil {
		return
	}
}

func (this *TAokoLog) flush() {
	this.writelog(<-this.data)
}
