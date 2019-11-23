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
	newInnerServer := tcpNet.NewTcpServer(serverConfig.GInnerGWConfig.ListenAddr,
		serverConfig.GInnerGWConfig.PProfAddr,
		Define.ERouteId_ER_ISG,
		LogicMsg.InnerGatewayMessageCallBack,
		tcpNet.GServer2ServerSession)

	newInnerServer.Run()
}
