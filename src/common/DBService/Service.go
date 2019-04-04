package DBService

import (
	"common/RedisService"
	"common/MgoService"
)

type TDBProvider struct {
	rconn 	*RedisService.RedisConn
	mconn   *MgoService.MgoConn
}

func (self *TDBProvider) StartDBService(RedisCfg *TRedisConfig, MgoCfg *TMgoConfig){
	self.rconn = RedisService.NewRedisConn(RedisCfg.ConnAddr, RedisCfg.DBIndex, RedisCfg.Passwd)
	self.mconn = MgoService.NewMgoConn(MgoCfg.UserName, MgoCfg.Passwd, MgoCfg.Host)

	go self.OnTimeSyncDBHard()
}

func (self *TDBProvider) OnTimeSyncDBHard(){
	ticker := time.NewTicker(time.Duration(EDB_DATA_SAVE_INTERVAL)*timer.second)
	for{
		select {
		case <-ticker.C:
			// do something...
			self.SyncDBHard()
		default:
			// nothing...
		}
	}
}

func (self.TDBProvider) SyncDBHard(){
	if self.rconn != nil {
		
	}

	if self.mconn != nil {

	}
}