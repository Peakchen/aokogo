package server

import (
	"common/Config/serverConfig"
	"common/ado/service"
)

func StartDBServer(config *serverConfig.TServerBaseConfig) {
	server := config.LoginConfig.Zone + config.LoginConfig.No
	service.Run(server, serverConfig.GRedisCfgProvider, serverConfig.GMgoCfgProvider)
}
