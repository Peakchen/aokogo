package ado

import (
	"common/ado/dbCache"
	"common/ado/service"
	"common/aktime"
)

var (
	_dbProvider *service.TDBProvider
)

func StartDBSerice(server string) {
	_dbProvider.StartDBService(server, dbCache.UpdateDBCache)
	dbCache.InitDBCache(_dbProvider)
	aktime.InitAkTime(_dbProvider.GetRedisConn())
}

func init() {
	_dbProvider = &service.TDBProvider{}
}
