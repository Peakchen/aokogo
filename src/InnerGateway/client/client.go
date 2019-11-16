package client

import (
	"ExternalGateway/SessionMgr"
	"common/Config/serverConfig"
	"common/Define"
	"common/Log"
	"common/tcpNet"
)

func StartClient() {
	Log.FmtPrintf("start InnerGateway client.")
	gameSvr := tcpNet.NewClient(serverConfig.GInnerGWConfig.ConnectAddr,
		Define.ERouteId_ER_ISG,
		nil,
		nil,
		SessionMgr.GClient2ServerSession)

	gameSvr.Run()
}
