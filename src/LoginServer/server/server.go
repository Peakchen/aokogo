package server

import (
	"LoginServer/LogicMsg"
	"LoginServer/dbo"
	"common/Config/serverConfig"
	"common/Define"
	"common/Log"
	"common/tcpNet"
	"flag"
)

var (
	CfgPath string
)

func init() {
	flag.StringVar(&CfgPath, "serverconfig", "serverconfig", "default path for configuration files")
}

func StartServer() {
	Log.FmtPrintf("start Login server.")
	serverConfig.LoadSvrAllConfig(CfgPath)
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
}
