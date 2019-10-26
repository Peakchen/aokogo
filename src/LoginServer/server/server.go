package server

import (
	"ExternalGateway/SessionMgr"
	"LoginServer/LogicMsg"
	"LoginServer/dbo"
	"common/Config/serverConfig"
	"common/Define"
	"common/tcpNet"
)

func StartServer(servercfg *serverConfig.TServerBaseConfig) {
	server := servercfg.LoginConfig.Zone + servercfg.LoginConfig.No
	dbo.StartDBSerice(server, servercfg.RedisConfig, servercfg.MgoConfig)

	gameSvr := tcpNet.NewClient(servercfg.LoginConfig.ListenAddr,
		Define.ERouteId_ER_Login,
		Define.ERouteId_ER_Login,
		Define.ERouteId_ER_ESG,
		LogicMsg.LoginMessageCallBack,
		nil,
		SessionMgr.GClient2ServerSession)

	gameSvr.Run()
}
