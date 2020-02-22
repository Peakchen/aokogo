// add by stefan

package main

import (
	"ExternalGateway/LogicMsg"
	"common/Config/serverConfig"
	"common/Define"
	"common/Log"
	"common/tcpNet"
	"flag"
)

var (
	CfgPath string
)

func init() {
	flag.StringVar(&CfgPath, "serverconfig", "serverconfig", "default path for configuration files")
}

func main() {
	Log.FmtPrintf("start ExternalGateWay.")
	serverConfig.LoadSvrAllConfig(CfgPath)
	externalgw := serverConfig.GExternalgwconfigConfig.Get()
	newExternalServer := tcpNet.NewTcpServer(externalgw.Listenaddr,
		externalgw.Pprofaddr,
		Define.ERouteId_ER_ESG,
		LogicMsg.ExternalGatewayMessageCallBack,
		tcpNet.GClient2ServerSession)

	newExternalServer.Run()
}
