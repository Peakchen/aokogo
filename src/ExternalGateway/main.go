// add by stefan

package main

import (
	"ExternalGateway/LogicMsg"
	"common/Config/serverConfig"
	"common/Define"
	"common/Log"
	"common/tcpNet"
)

func init() {

}

func main() {
	Log.FmtPrintf("start ExternalGateWay.")

	externalgw := serverConfig.GExternalgwconfigConfig.Get()
	newExternalServer := tcpNet.NewTcpServer(externalgw.Listenaddr,
		externalgw.Pprofaddr,
		Define.ERouteId_ER_ESG,
		LogicMsg.ExternalGatewayMessageCallBack,
		tcpNet.GClient2ServerSession)

	newExternalServer.Run()
}
