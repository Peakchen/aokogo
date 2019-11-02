package service

import (
	"common/Config/serverConfig"
)

var (
	GDBProvider *TDBProvider
)

func NewDBProvider() {
	GDBProvider = &TDBProvider{}
}

func Run(server string, RedisCfg *serverConfig.TRedisConfig, MgoCfg *serverConfig.TMgoConfig) {
	GDBProvider.StartDBService(server, RedisCfg, MgoCfg)
}

func init() {
	NewDBProvider()
}
