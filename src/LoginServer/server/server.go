package server

import (
	"LoginServer/LogicMsg"
	"LoginServer/dbo"
	"common/Config/serverConfig"
	"common/Define"
	"common/Log"
	"common/ado/dbStatistics"
	"common/tcpNet"
	"flag"
)

func init() {
	var CfgPath string
	flag.StringVar(&CfgPath, "serverconfig", "serverconfig", "default path for configuration files")
	serverConfig.LoadSvrAllConfig(CfgPath)
	dbStatistics.InitDBStatistics()
}

func StartServer() {
	Log.FmtPrintf("start Login server.")
	logincfg := serverConfig.GLoginconfigConfig.Get()
	server := logincfg.Zone + logincfg.No
	dbo.StartDBSerice(server)
	gameSvr := tcpNet.NewClient(logincfg.Listenaddr,
		logincfg.Pprofaddr,
		Define.ERouteId_ER_Login,
		LogicMsg.LoginMessageCallBack,
		nil,
		tcpNet.GClient2ServerSession)

	gameSvr.Run()
	dbStatistics.DBStatisticsStop()
}
