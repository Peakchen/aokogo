package service

import (
	"common/Config/serverConfig"
	"common/MgoConn"
	"common/RedisConn"
	"context"
	"sync"
)

type TDBProvider struct {
	rconn  *RedisConn.TRedisConn
	mconn  *MgoConn.AokoMgo
	Server string
	ctx    context.Context
	cancle context.CancelFunc
	wg     sync.WaitGroup
}

func (this *TDBProvider) StartDBService(Server string) {
	this.Server = Server
	rediscfg := serverConfig.GRedisconfigConfig.Get()
	this.rconn = RedisConn.NewRedisConn(rediscfg.Connaddr, rediscfg.DBIndex, rediscfg.Passwd)

	mgocfg := serverConfig.GMgoconfigConfig.Get()
	this.mconn = MgoConn.NewMgoConn(Server, mgocfg.Username, mgocfg.Passwd, mgocfg.Host)
}
