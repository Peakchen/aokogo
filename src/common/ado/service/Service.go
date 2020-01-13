package service

import (
	"common/Config/serverConfig"
	"common/MgoConn"
	"common/RedisConn"
	"context"
	"sync"
	"common/ado/dbCache"
)

type TDBProvider struct {
	rconn  *RedisConn.TAokoRedis
	mconn  *MgoConn.TAokoMgo
	Server string
	ctx    context.Context
	cancle context.CancelFunc
	wg     sync.WaitGroup
}

func (this *TDBProvider) StartDBService(Server string) {
	this.Server = Server
	rediscfg := serverConfig.GRedisconfigConfig.Get()
	this.rconn = RedisConn.NewRedisConn(rediscfg.Connaddr, rediscfg.DBIndex, rediscfg.Passwd)
	dbCache.Init(this.rconn)
	mgocfg := serverConfig.GMgoconfigConfig.Get()
	this.mconn = MgoConn.NewMgoConn(Server, mgocfg.Username, mgocfg.Passwd, mgocfg.Host)
}

func (this *TDBProvider) GetRedisConn()*RedisConn.TAokoRedis{
	return this.rconn
}

func (this *TDBProvider) GetMogoConn()*MgoConn.TAokoMgo{
	return this.mconn
}