package server

import (
	"LoginServer/LogicMsg"
	"LoginServer/dbo"
	"common/Config/serverConfig"
	"common/Define"
	"common/tcpNet"
)

func StartServer() {
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
