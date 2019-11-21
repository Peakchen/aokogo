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
		servercfg.GameConfig.PProfAddr,
		Define.ERouteId_ER_Game,
		LogicMsg.GameMessageCallBack,
		LogicMsg.AfterDialCallBack,
		nil)

	gameSvr.Run()
}
