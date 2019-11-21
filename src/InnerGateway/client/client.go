package client

import (
	"common/Config/serverConfig"
	"common/Define"
	"common/Log"
	"common/tcpNet"
)

func StartClient() {
	Log.FmtPrintf("start InnerGateway client.")
	gameSvr := tcpNet.NewClient(serverConfig.GInnerGWConfig.ConnectAddr,
		serverConfig.GInnerGWConfig.PProfAddr,
		Define.ERouteId_ER_ISG,
		nil,
		nil,
		tcpNet.GClient2ServerSession)

	gameSvr.Run()
}
