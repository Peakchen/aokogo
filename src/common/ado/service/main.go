package service

import (
	"common/Config"
)

var (
	GDBProvider *TDBProvider
)

func NewDBProvider() {
	GDBProvider = &TDBProvider{}
}

func Run(RedisCfg *Config.TRedisConfig, MgoCfg *Config.TMgoConfig) {
	GDBProvider.StartDBService(RedisCfg, MgoCfg)
}

func init() {
	NewDBProvider()
}
