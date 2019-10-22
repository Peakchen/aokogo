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

func Run(RedisCfg *serverConfig.TRedisConfig, MgoCfg *serverConfig.TMgoConfig) {
	GDBProvider.StartDBService(RedisCfg, MgoCfg, true)
}

func init() {
	NewDBProvider()
}
