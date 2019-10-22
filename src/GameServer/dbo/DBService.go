package dbo

import (
	"common/Config/serverConfig"
	"common/ado"
	"common/ado/service"
	"common/public"
)

var (
	GDBProvider *service.TDBProvider
)

func A_DataGet(key string, Out public.IDBCache) {
	// check redis can get db data, if not exist, then from mogo.
	GDBProvider.DBGet(key, Out)
}

func A_DataSet(key string, In public.IDBCache) {
	// check save data to redis cache or db persistence.
	GDBProvider.DBSet(key, In, ado.EDBOper_Update)
}

func StartDBSerice(RedisCfg *serverConfig.TRedisConfig, MgoCfg *serverConfig.TMgoConfig) {
	GDBProvider.StartDBService(RedisCfg, MgoCfg, false)
}

func init() {
	GDBProvider = &service.TDBProvider{}
}
