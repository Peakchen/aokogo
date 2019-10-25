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

func A_DBRead(Out public.IDBCache) {
	// check redis can get db data, if not exist, then from mogo.
	GDBProvider.Get(Out)
}

func A_DBUpdate(In public.IDBCache) {
	// check save data to redis cache or db persistence.
	GDBProvider.Update(In, ado.EDBOper_Update)
}

func StartDBSerice(RedisCfg *serverConfig.TRedisConfig, MgoCfg *serverConfig.TMgoConfig) {
	GDBProvider.StartDBService(RedisCfg, MgoCfg, false)
}

func init() {
	GDBProvider = &service.TDBProvider{}
}
