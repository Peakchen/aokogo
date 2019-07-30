package DBService

import (
	"common/Log"
	"common/MgoService"
	"common/RedisService"
	"time"
)

type TDBProvider struct {
	rconn       *RedisService.TRedisConn
	mconn       *MgoService.AokoMgo
	ServerModel EServerModel
}

func (self *TDBProvider) StartDBService(RedisCfg *TRedisConfig, MgoCfg *TMgoConfig) {
	self.rconn = RedisService.NewRedisConn(RedisCfg.ConnAddr, RedisCfg.DBIndex, RedisCfg.Passwd)
	self.mconn = MgoService.NewMgoConn(MgoCfg.UserName, MgoCfg.Passwd, MgoCfg.Host)

	go self.LoopSyncDBHard()
}

func (self *TDBProvider) LoopSyncDBHard() {
	ticker := time.NewTicker(time.Duration(EDB_DATA_SAVE_INTERVAL) * timer.second)
	for {
		select {
		case <-ticker.C:
			// do something...
			self.SyncDBHard()
		default:
			// nothing...
		}
	}
}

func (self *TDBProvider) SyncDBHard() {
	// TODO: Presist redis...
	if self.rconn == nil {
		return
	}

	DBKey := ":" + self.ServerModel + "_Update_Oper"
	Members, err := self.rconn.Conn.Do("SMEMBERS", DBKey)
	if err != nil {
		Log.Error("DBProvider get redis ", err)
		return
	}

	// TODO: Presist mgo...
	if self.mconn == nil {
		return
	}

	for _, item := range Members.(IDBCache) {
		self.mconn.SaveOne(item)
	}
}
