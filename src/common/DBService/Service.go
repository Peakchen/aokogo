package DBService

import (
	"common/RedisService"
	"common/MgoService"
)

type TDBProvider struct {
	rconn 	*RedisService.RedisConn
	mconn   *MgoService.MgoConn
	ServerModel EServerModel 
}

func (self *TDBProvider) StartDBService(RedisCfg *TRedisConfig, MgoCfg *TMgoConfig){
	self.rconn = RedisService.NewRedisConn(RedisCfg.ConnAddr, RedisCfg.DBIndex, RedisCfg.Passwd)
	self.mconn = MgoService.NewMgoConn(MgoCfg.UserName, MgoCfg.Passwd, MgoCfg.Host)

	go self.LoopSyncDBHard()
}

func (self *TDBProvider) LoopSyncDBHard(){
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
		// TODO: Presist redis... 

	}

	if self.mconn != nil {
		// TODO: Presist mgo... 
		
	}
}


