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
	"common/tcpNet"
	"flag"
	"syscall"
)

var (
	CfgPath string
)

func init() {
	flag.StringVar(&CfgPath, "serverconfig", "serverconfig", "default path for configuration files")
	LogicMsg.Init()
	rpc.Init()
}

func reloadConfig() {
	LogicConfig.LoadLogicAll()
}

func StartServer() {
	serverConfig.LoadSvrAllConfig(CfgPath)
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
