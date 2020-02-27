package dbo

import (
	"common/ado"
	"common/ado/dbCache"
	"common/ado/service"
	"common/aktime"
	"common/public"
)

var (
	_dbProvider *service.TDBProvider
)

func A_DBRead(Identify string, Out public.IDBCache) (err error, exist bool) {
	// check redis can get db data, if not exist, then from mogo.
	err, exist = _dbProvider.Get(Identify, Out)
	return
}

func A_DBInsert(Identify string, In public.IDBCache) (err error) {
	// check save data to redis cache or db persistence.
	err = _dbProvider.Insert(Identify, In)
	return
}

func A_DBUpdate(Identify string, In public.IDBCache) (err error) {
	// check save data to redis cache or db persistence.
	err = _dbProvider.Update(Identify, In, ado.EDBOper_Update)
	return
}

func StartDBSerice(server string) {
	_dbProvider.StartDBService(server, dbCache.UpdateDBCache)
	dbCache.InitDBCache(_dbProvider)
	aktime.InitAkTime(_dbProvider.GetRedisConn())
}

func init() {
	_dbProvider = &service.TDBProvider{}
}
