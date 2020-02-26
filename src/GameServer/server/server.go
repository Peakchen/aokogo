package server

// add by stefan

import (
	"GameServer/LogicMsg"
	"GameServer/dbo"
	"GameServer/rpc"
	"common/Config/LogicConfig"
	"common/Config/serverConfig"
	"common/Define"
	"common/HotUpdate"
	"common/akNet"
	"flag"
	"syscall"
)

func init() {
	var CfgPath string
	flag.StringVar(&CfgPath, "serverconfig", "serverconfig", "default path for configuration files")
	serverConfig.LoadSvrAllConfig(CfgPath)
	LogicMsg.Init()
	rpc.Init()
}

func reloadConfig() {
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
	gameSvr := akNet.NewClient(Gamecfg.Listenaddr,
		Gamecfg.Pprofaddr,
		Define.ERouteId_ER_Game,
		nil,
		nil,
		nil)

	gameSvr.Run()
}
