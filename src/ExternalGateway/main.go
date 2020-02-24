// add by stefan

package main

import (
	"ExternalGateway/LogicMsg"
	"common/Config/serverConfig"
	"common/Define"
	"common/Log"
	"common/ado/dbStatistics"
	"common/tcpNet"
	"flag"
)

func init() {
	var CfgPath string
	flag.StringVar(&CfgPath, "serverconfig", "serverconfig", "default path for configuration files")
	serverConfig.LoadSvrAllConfig(CfgPath)
	dbStatistics.InitDBStatistics()
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
	dbStatistics.DBStatisticsStop()
}
