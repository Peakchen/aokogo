package server

import (
	"InnerGateway/LogicMsg"
	"common/Config/serverConfig"
	"common/Define"
	"common/Log"
	"common/tcpNet"
	"context"
	"sync"
)

func StartServer() {
	Log.FmtPrintf("start InnerGateway server.")
	newInnerServer := tcpNet.NewTcpServer(serverConfig.GInnerGWConfig.ListenAddr,
		Define.ERouteId_ER_ISG,
		LogicMsg.InnerGatewayMessageCallBack,
		tcpNet.GServer2ServerSession)

	sw := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	newInnerServer.StartTcpServer(&sw, ctx, cancel)
}
