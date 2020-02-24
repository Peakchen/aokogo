package server

// add by stefan

import (
	"common/Config/serverConfig"
	"common/ado/dbStatistics"
	"common/ado/service"
	"flag"
)

func init() {
	var CfgPath string
	flag.StringVar(&CfgPath, "serverconfig", "serverconfig", "default path for configuration files")
	serverConfig.LoadSvrAllConfig(CfgPath)
	dbStatistics.InitDBStatistics()
}

/*
	run db server.
*/
func StartDBServer() {
	server := "sever1"
	service.StartMultiDBProvider(server)
	dbStatistics.DBStatisticsStop()
}
