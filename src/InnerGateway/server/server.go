package server

import (
	"InnerGateway/LogicMsg"
	"common/Config/serverConfig"
	"common/Define"
	"common/Log"
	"common/tcpNet"
)

func StartServer() {
	Log.FmtPrintf("start InnerGateway server.")
	Innergw := serverConfig.GInnergwconfigConfig.Get()
	newInnerServer := tcpNet.NewTcpServer(Innergw.Listenaddr,
		Innergw.Pprofaddr,
		Define.ERouteId_ER_ISG,
		LogicMsg.InnerGatewayMessageCallBack,
		tcpNet.GServer2ServerSession)

	newInnerServer.Run()
}
