// add by stefan
package main

import (
	"InnerGateway/LogicMsg"
	"InnerGateway/ado"
	"InnerGateway/client"
	"InnerGateway/server"
	"common/Config/serverConfig"
	"common/Log"
	"common/ado/dbStatistics"
	"flag"
	"sync"
)

func init() {
	var CfgPath string
	flag.StringVar(&CfgPath, "serverconfig", "serverconfig", "default path for configuration files")
	serverConfig.LoadSvrAllConfig(CfgPath)
	dbStatistics.InitDBStatistics()
	LogicMsg.Init()
}

func startInnerGW() {
	var sw sync.WaitGroup
	sw.Add(2)
	go server.StartServer()
	go client.StartClient()
	sw.Wait()
}

func main() {
	Log.FmtPrintf("start InnerGateway.")
	ado.StartDBSerice("InnerGateway")
	startInnerGW()
	dbStatistics.DBStatisticsStop()
}
