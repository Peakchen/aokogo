package pprof

// add by stefan 20190606 16:12
import (
	"common/Log"
	"fmt"
	"os"
	"path"
	"reflect"
	"runtime/pprof"
	"strings"
	"syscall"
	"time"

	//"log"
	"context"
	"sync"
)

const (
	const_PProfWriteInterval = int32(60 * 1)
)

type TPProfMgr struct {
	ctx context.Context
	wg  sync.WaitGroup
	cpu *os.File
	mem *os.File
}

var (
	_pprofobj *TPProfMgr
)

func init() {
	_pprofobj = &TPProfMgr{}
}

func Run(ctx context.Context) {
	_pprofobj.StartPProf(ctx)
}

func Exit() {
	_pprofobj.Exit()
}

func (this *TPProfMgr) StartPProf(ctx context.Context) {
	this.ctx = ctx
	this.wg.Add(1)
	checkcreateTempDir()
	this.cpu = createCpu()
	this.mem = createMem()
	go this.loop()
}

func (this *TPProfMgr) Exit() {
	Log.FmtPrintln("pprof exist.")
	this.flush()
}

func (this *TPProfMgr) flush() {
	//Log.FmtPrintln("pprof flush.")
	if this.cpu != nil {
		pprof.StopCPUProfile()
		this.cpu.Close()
	}
	if this.mem != nil {
		pprof.WriteHeapProfile(this.mem)
		this.mem.Close()
	}
}

func (this *TPProfMgr) loop() {
	defer this.wg.Done()
	t := time.NewTicker(time.Duration(const_PProfWriteInterval) * time.Second)
	for {
		select {
		case <-this.ctx.Done():
			this.Exit()
		case <-t.C:
			// do nothing...
			this.flush()
		}
	}
}

func Newpprof(file string) (retfile string) {
	timeformat := time.Now().Format("2006-01-02")
	retfile = timeformat + "_" + file
	execpath, err := os.Executable()
	if err != nil {
		return
	}
	execpath = strings.Replace(execpath, "\\", "/", -1)
	_, sfile := path.Split(execpath)
	arrfile := strings.Split(sfile, ".")
	retfile = fmt.Sprintf("./pprof/%s_%v.prof", arrfile[0], retfile)
	return
}

func checkcreateTempDir() {
	err := os.Mkdir("pprof", os.ModePerm)
	if err != nil {
		if reflect.TypeOf(err) != reflect.TypeOf(&os.PathError{}) {
			Log.FmtPrintln("err dir type: ", reflect.TypeOf(err))
			return
		}
		perror := err.(*os.PathError)
		if perror.Err != syscall.ERROR_ALREADY_EXISTS {
			Log.FmtPrintf("pprof mkdir fail, dir: %v, errcode: %v, err: %v.\n", perror.Err, err.Error())
			return
		}
	}
}

func createCpu() (file *os.File) {
	cpuf := Newpprof("cpu")
	f, err := os.OpenFile(cpuf, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		Log.FmtPrintln("cpu pprof open fail, err: ", err)
		return
	}
	pprof.StartCPUProfile(f)
	return f
}

func createMem() (file *os.File) {
	cpuf := Newpprof("mem")
	f, err := os.OpenFile(cpuf, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		Log.FmtPrintln("mem pprof open fail, err: ", err)
		return
	}
	return f
}
