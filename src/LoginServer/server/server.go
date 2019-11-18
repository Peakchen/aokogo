package server

import (
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
		LogicMsg.LoginMessageCallBack,
		nil,
		tcpNet.GClient2ServerSession)

	gameSvr.Run()
}
