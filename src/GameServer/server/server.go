package server

// add by stefan

import (
	"GameServer/LogicMsg"
	"GameServer/dbo"
	"GameServer/rpc"
	"common/Config/serverConfig"
	"common/Define"
	"common/tcpNet"
	"common/HotUpdate"
	"syscall"
	"common/Config/LogicConfig"
)

func init() {
	LogicMsg.Init()
	rpc.Init()
}

func reloadConfig(){
	LogicConfig.LoadLogicAll()
}

func StartServer() {
	Gamecfg := serverConfig.GGameconfigConfig.Get()
	server := Gamecfg.Zone + Gamecfg.No
	dbo.StartDBSerice(server)
	// for kill pid to emit signal to do action...
	HotUpdate.RunHotUpdateCheck(&HotUpdate.TServerHotUpdateInfo{
		Recvsignal: syscall.SIGTERM,
		HUCallback: reloadConfig,
	})
	gameSvr := tcpNet.NewClient(Gamecfg.Listenaddr,
		Gamecfg.Pprofaddr,
		Define.ERouteId_ER_Game,
		nil,
		nil,
		nil)

	gameSvr.Run()
}
