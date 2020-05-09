package server

import (
	"common/Config/serverConfig"
	"common/Log"
	"common/akNet"
	"common/define"
)

func StartServer() {
	Log.FmtPrintf("start InnerGateway server.")
	Innergw := serverConfig.GInnergwconfigConfig.Get()
	newInnerServer := akNet.NewTcpServer(Innergw.Listenaddr,
		Innergw.Pprofaddr,
		define.ERouteId_ER_ISG,
		Innergw.Name)

	newInnerServer.Run()
}
