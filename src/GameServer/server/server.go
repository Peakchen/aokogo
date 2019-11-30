package server

import (
	"GameServer/LogicMsg"
	"GameServer/dbo"
	"common/Config/serverConfig"
	"common/Define"
	"common/tcpNet"
)

func StartServer() {
	Gamecfg := serverConfig.GGameconfigConfig.Get()
	server := Gamecfg.Zone + Gamecfg.No
	dbo.StartDBSerice(server)
	gameSvr := tcpNet.NewClient(Gamecfg.Listenaddr,
		Gamecfg.Pprofaddr,
		Define.ERouteId_ER_Game,
		LogicMsg.GameMessageCallBack,
		LogicMsg.AfterDialCallBack,
		nil)

	gameSvr.Run()
}
