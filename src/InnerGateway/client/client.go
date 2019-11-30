package client

import (
	"common/Config/serverConfig"
	"common/Define"
	"common/Log"
	"common/tcpNet"
)

func StartClient() {
	Log.FmtPrintf("start InnerGateway client.")
	Innergw := serverConfig.GInnergwconfigConfig.Get()
	gameSvr := tcpNet.NewClient(Innergw.Connectaddr,
		Innergw.Pprofaddr,
		Define.ERouteId_ER_ISG,
		nil,
		nil,
		tcpNet.GClient2ServerSession)

	gameSvr.Run()
}
