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

func A_DBRead(Out public.IDBCache) (err error, exist bool) {
	// check redis can get db data, if not exist, then from mogo.
	err, exist = GDBProvider.Get(Out)
	return
}

func A_DBUpdate(In public.IDBCache) (err error) {
	// check save data to redis cache or db persistence.
	err = GDBProvider.Update(In, ado.EDBOper_Update)
	return
}

func A_DBInsert(In public.IDBCache) (err error) {
	// check save data to redis cache or db persistence.
	err = GDBProvider.Insert(In)
	return

}

func StartDBSerice(RedisCfg *serverConfig.TRedisConfig, MgoCfg *serverConfig.TMgoConfig) {
	GDBProvider.StartDBService(RedisCfg, MgoCfg, false)
}

func init() {
	GDBProvider = &service.TDBProvider{}
}
