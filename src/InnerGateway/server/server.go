package server

import (
	"InnerGateway/LogicMsg"
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
	Log.FmtPrintf("start InnerGateway server.")
	serverConfig.LoadSvrAllConfig(CfgPath)
	Innergw := serverConfig.GInnergwconfigConfig.Get()
	newInnerServer := tcpNet.NewTcpServer(Innergw.Listenaddr,
		Innergw.Pprofaddr,
		Define.ERouteId_ER_ISG,
		LogicMsg.InnerGatewayMessageCallBack,
		tcpNet.GServer2ServerSession)

	newInnerServer.Run()
}
