package server

import (
	"common/Config/serverConfig"
	"common/ado/service"
)

/*
	run db server.
*/
func StartDBServer(config *serverConfig.TServerBaseConfig) {
	server := config.LoginConfig.Zone + config.LoginConfig.No
	service.StartMultiDBProvider(server, serverConfig.GRedisCfgProvider, serverConfig.GMgoCfgProvider)
}
