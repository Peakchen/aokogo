// add by stefan

package main

import (
	"ExternalGateway/LogicMsg"
	"ExternalGateway/ado"
	"common/Config/serverConfig"
	"common/Log"
	"common/ado/dbStatistics"
	"common/akNet"
	"common/define"
	"flag"
)

func init() {
	var CfgPath string
	flag.StringVar(&CfgPath, "serverconfig", "serverconfig", "default path for configuration files")
	serverConfig.LoadSvrAllConfig(CfgPath)
	dbStatistics.InitDBStatistics()
	LogicMsg.Init()
}

func main() {
	Log.FmtPrintf("start ExternalGateWay.")
	ado.StartDBSerice("ExternalGateWay")
	externalgw := serverConfig.GExternalgwconfigConfig.Get()
	newExternalServer := akNet.NewTcpServer(externalgw.Listenaddr,
		externalgw.Pprofaddr,
		define.ERouteId_ER_ESG,
		externalgw.Name)

	newExternalServer.Run()
	dbStatistics.DBStatisticsStop()
}
