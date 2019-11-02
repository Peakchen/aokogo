package server

import (
	"GameServer/LogicMsg"
	"GameServer/dbo"
	"common/Config/serverConfig"
	"common/Define"
	"common/tcpNet"
)

func StartServer(servercfg *serverConfig.TServerBaseConfig) {
	server := servercfg.GameConfig.Zone + servercfg.GameConfig.No
	dbo.StartDBSerice(server, servercfg.RedisConfig, servercfg.MgoConfig)
	gameSvr := tcpNet.NewClient(servercfg.GameConfig.ListenAddr,
		Define.ERouteId_ER_Game,
		Define.ERouteId_ER_Client,
		Define.ERouteId_ER_ESG,
		LogicMsg.GameMessageCallBack,
		LogicMsg.AfterDialCallBack,
		nil)

	gameSvr.Run()
}
