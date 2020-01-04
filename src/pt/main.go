package main

import (
	
	//"fmt"
	"net/http"
	"os"
	//"common/cfg"
	"pt/db"
	"pt/regtable"
	"strconv"
	//"strings"
	"common/tool"
	//"errors"
	"time"
	//"os/exec"
	//"bytes"
)

func init() {

}

func getArgs() (cmds []string) {
	cmds = []string{}
	for i, v := range os.Args {
		tool.MyFmtPrint_Info("args[%v]=%v\n", i, v)
		cmds = append(cmds, v)
	}
	return
}

func registHttpHandle(){
	http.HandleFunc("/BatchUpdate", db.BatchUpdate)

	http.HandleFunc("/Update", db.DBUpdate)

	http.HandleFunc("/Query", db.DBQuery)
}

var (
	Addr = string("0.0.0.0")
	beginPort = int(8000)
	ListenPort int
)

func checkListenHttpSvr() {
	// sqlcfg := cfg.GSqlconfig.Get()
	// if sqlcfg == nil {
	// 	tool.MyFmtPrint_Info("sql config load fail.")
	// 	return
	// }

	ListenPort = beginPort
	for tool.CheckPortUsed(ListenPort) {
		ListenPort++
	}

	host := Addr + ":" + strconv.Itoa(ListenPort)
	db.SetListenHost(host)
	regtable.CreateRegPortKey(int32(ListenPort))
	tool.MyFmtPrint_Info("first host: ", host)
	err := http.ListenAndServe(host, nil)
	if err == nil {
		return
	}else{
		tool.MyFmtPrint_Info("listen server fail, info: ", host, err)
	}

	for {
		ListenPort++
		regtable.CreateRegPortKey(int32(ListenPort))
		host := Addr + ":" + strconv.Itoa(ListenPort)
		tool.MyFmtPrint_Info("next host: ", host)
		db.SetListenHost(host)
		err := http.ListenAndServe(host, nil)
		if err == nil {
			return
		}else{
			tool.MyFmtPrint_Info("listen server fail, info: ", host, err)
		}
	}
}

func checkAddrPortloop(){
	tick := time.NewTicker(time.Duration(3)*time.Second)
	for {
		select {
		case <-tick.C:
			//tool.MyFmtPrint_Info("new port: ", ListenPort, regtable.GetRegPoryKey())
			if regtable.GetRegPoryKey() < uint32(ListenPort) {
				regtable.CreateRegPortKey(int32(ListenPort))
			}
		}
	}
}


func main() {
	//cmds := getArgs()
	//if len(cmds) >= 2 {
		//tool.CmdHide()
		//tool.BatHide(cmds)
		//return
	//}else {
		//tool.BatHide([]string{"pt.exe", "d"})
		//tool.HideConsole()
	//}

	if tool.CstHideConsole {
		tool.HideConsole()
	}
	
	db.LoadDB()
	go checkAddrPortloop()
	registHttpHandle()
	checkListenHttpSvr()
}

