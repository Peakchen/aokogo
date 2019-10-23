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

func A_DataGet(Out public.IDBCache) (err error) {
	// check redis can get db data, if not exist, then from mogo.
	err = GDBProvider.DBGet(Out)
	return
}

func A_DataSet(In public.IDBCache) (err error) {
	// check save data to redis cache or db persistence.
	err = GDBProvider.DBSet(In, ado.EDBOper_Update)
	return
}

func StartDBSerice(RedisCfg *serverConfig.TRedisConfig, MgoCfg *serverConfig.TMgoConfig) {
	GDBProvider.StartDBService(RedisCfg, MgoCfg, false)
}

func init() {
	GDBProvider = &service.TDBProvider{}
}
