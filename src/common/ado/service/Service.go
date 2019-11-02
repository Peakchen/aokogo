package service

import (
	"common/Config/serverConfig"
	"common/MgoService"
	"common/RedisService"
	"context"
	"sync"
)

type TDBProvider struct {
	rconn  *RedisService.TRedisConn
	mconn  *MgoService.AokoMgo
	Server string
	ctx    context.Context
	cancle context.CancelFunc
	wg     sync.WaitGroup
}

func (this *TDBProvider) init(Server string, RedisCfg *serverConfig.TRedisConfig, MgoCfg *serverConfig.TMgoConfig) {
	this.Server = Server
	this.rconn = RedisService.NewRedisConn(RedisCfg.ConnAddr, RedisCfg.DBIndex, RedisCfg.Passwd)
	this.mconn = MgoService.NewMgoConn(Server, MgoCfg.UserName, MgoCfg.Passwd, MgoCfg.Host)
}

func (this *TDBProvider) StartDBService(Server string, RedisCfg *serverConfig.TRedisConfig, MgoCfg *serverConfig.TMgoConfig) {
	this.init(Server, RedisCfg, MgoCfg)
}
