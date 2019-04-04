package DBService

import (
	"common/RedisService"
	"common/MgoService"
)

type TDBService struct {
	rconn 	*RedisService.RedisConn
	mconn   *MgoService.MgoConn
}

func (self *TDBService) StartDBService(RedisCfg *TRedisConfig, MgoCfg *TMgoConfig){
	self.rconn = RedisService.NewRedisConn(RedisCfg.ConnAddr, RedisCfg.DBIndex, RedisCfg.Passwd)
	self.mconn = MgoService.NewMgoConn(MgoCfg.UserName, MgoCfg.Passwd, MgoCfg.Host)

	
}