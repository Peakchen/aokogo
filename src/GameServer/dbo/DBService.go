package dbo

import (
	"common/ado"
	"common/ado/dbCache"
	"common/ado/service"
	"common/public"
)

var (
	GDBProvider *service.TDBProvider
)

func A_DBRead(Identify string, Out public.IDBCache) (err error, exist bool) {
	// check redis can get db data, if not exist, then from mogo.
	err, exist = GDBProvider.Get(Identify, Out)
	return
}

func A_DBInsert(Identify string, In public.IDBCache) (err error) {
	// check save data to redis cache or db persistence.
	err = GDBProvider.Insert(Identify, In)
	return
}

func A_DBUpdate(Identify string, In public.IDBCache) (err error) {
	// check save data to redis cache or db persistence.
	err = GDBProvider.Update(Identify, In, ado.EDBOper_Update)
	return
}

func StartDBSerice(server string) {
	GDBProvider.StartDBService(server, dbCache.UpdateDBCache)
	dbCache.InitDBCache(GDBProvider)
}

func init() {
	GDBProvider = &service.TDBProvider{}
}
