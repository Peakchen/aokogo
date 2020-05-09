package server

import (
	"LoginServer/LogicMsg"
	"LoginServer/dbo"
	"common/Config/serverConfig"
	"common/Log"
	"common/ado/dbStatistics"
	"common/akNet"
	"common/define"
	"flag"
)

func init() {
	var CfgPath string
	flag.StringVar(&CfgPath, "serverconfig", "serverconfig", "default path for configuration files")
	serverConfig.LoadSvrAllConfig(CfgPath)
	dbStatistics.InitDBStatistics()
	LogicMsg.Init()
}

func StartServer() {
	Log.FmtPrintf("start Login server.")
	logincfg := serverConfig.GLoginconfigConfig.Get()
	server := logincfg.Zone + logincfg.No
	dbo.StartDBSerice(server)
	gameSvr := akNet.NewClient(logincfg.Listenaddr,
		logincfg.Pprofaddr,
		define.ERouteId_ER_Login,
		nil,
		logincfg.Name)

	gameSvr.Run()
	dbStatistics.DBStatisticsStop()
}
