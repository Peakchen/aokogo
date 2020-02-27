// add by stefan

package main

import (
	"ExternalGateway/LogicMsg"
	"ExternalGateway/ado"
	"common/Config/serverConfig"
	"common/Define"
	"common/Log"
	"common/ado/dbStatistics"
	"common/akNet"
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
	ado.StartDBSerice("ExternalGateWay")
	externalgw := serverConfig.GExternalgwconfigConfig.Get()
	newExternalServer := akNet.NewTcpServer(externalgw.Listenaddr,
		externalgw.Pprofaddr,
		Define.ERouteId_ER_ESG,
		LogicMsg.ExternalGatewayMessageCallBack,
		akNet.GClient2ServerSession)

	newExternalServer.Run()
	dbStatistics.DBStatisticsStop()
}
