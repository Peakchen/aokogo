package server

import (
	"InnerGateway/LogicMsg"
	"common/Config/serverConfig"
	"common/Define"
	"common/Log"
	"common/akNet"
)

func StartServer() {
	Log.FmtPrintf("start InnerGateway server.")
	Innergw := serverConfig.GInnergwconfigConfig.Get()
	newInnerServer := akNet.NewTcpServer(Innergw.Listenaddr,
		Innergw.Pprofaddr,
		Define.ERouteId_ER_ISG,
		LogicMsg.InnerGatewayMessageCallBack,
		akNet.GServer2ServerSession)

	newInnerServer.Run()
}
