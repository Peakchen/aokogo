package client

import (
	"common/Config/serverConfig"
	"common/Define"
	"common/Log"
	"common/akNet"
)

func StartClient() {
	Log.FmtPrintf("start InnerGateway client.")
	Innergw := serverConfig.GInnergwconfigConfig.Get()
	gameSvr := akNet.NewClient(Innergw.Connectaddr,
		Innergw.Pprofaddr,
		Define.ERouteId_ER_ISG,
		nil,
		Innergw.Name)

	gameSvr.Run()
}
