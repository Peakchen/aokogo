package client

import (
	"common/Config/serverConfig"
	"common/Log"
	"common/akNet"
	"common/define"
)

func StartClient() {
	Log.FmtPrintf("start InnerGateway client.")
	Innergw := serverConfig.GInnergwconfigConfig.Get()
	gameSvr := akNet.NewClient(Innergw.Connectaddr,
		Innergw.Pprofaddr,
		define.ERouteId_ER_ISG,
		nil,
		Innergw.Name)

	gameSvr.Run()
}
