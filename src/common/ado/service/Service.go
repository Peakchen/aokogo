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

func (this *TDBProvider) init(Server string, RedisCfg *serverConfig.TRedisConfig, MgoCfg *serverConfig.TMgoConfig) {
	this.Server = Server
	this.rconn = RedisConn.NewRedisConn(RedisCfg.ConnAddr, RedisCfg.DBIndex, RedisCfg.Passwd)
	this.mconn = MgoConn.NewMgoConn(Server, MgoCfg.UserName, MgoCfg.Passwd, MgoCfg.Host)
}

func (this *TDBProvider) StartDBService(Server string, RedisCfg *serverConfig.TRedisConfig, MgoCfg *serverConfig.TMgoConfig) {
	this.init(Server, RedisCfg, MgoCfg)
}
