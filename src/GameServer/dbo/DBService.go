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

func A_DBRead(Identify string, Out public.IDBCache) {
	// check redis can get db data, if not exist, then from mogo.
	GDBProvider.Get(Identify, Out)
}

func A_DBUpdate(Identify string, In public.IDBCache) {
	// check save data to redis cache or db persistence.
	GDBProvider.Update(Identify, In, ado.EDBOper_Update)
}

func StartDBSerice(server string, RedisCfg *serverConfig.TRedisConfig, MgoCfg *serverConfig.TMgoConfig) {
	GDBProvider.StartDBService(server, RedisCfg, MgoCfg)
}

func init() {
	GDBProvider = &service.TDBProvider{}
}
