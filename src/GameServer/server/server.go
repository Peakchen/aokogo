package server

import (
	"GameServer/LogicMsg"
	"GameServer/dbo"
	"GameServer/rpc"
	"common/Config/serverConfig"
	"common/Define"
	"common/tcpNet"
)

func init() {
	LogicMsg.Init()
	rpc.Init()
}

func StartServer() {
	Gamecfg := serverConfig.GGameconfigConfig.Get()
	server := Gamecfg.Zone + Gamecfg.No
	dbo.StartDBSerice(server)
	gameSvr := tcpNet.NewClient(Gamecfg.Listenaddr,
		Gamecfg.Pprofaddr,
		Define.ERouteId_ER_Game,
		nil,
		nil,
		nil)

	gameSvr.Run()
}
