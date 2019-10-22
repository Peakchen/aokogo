package service

import (
	"common/Config/serverConfig"
	"common/Log"
	"common/MgoService"
	"common/RedisService"
	"common/ado"
	"context"
	"sync"
	"time"
)

type TDBProvider struct {
	rconn       *RedisService.TRedisConn
	mconn       *MgoService.AokoMgo
	ServerModel ado.EServerModel
	ctx         context.Context
	cancle      context.CancelFunc
	wg          sync.WaitGroup
}

func (this *TDBProvider) StartDBService(RedisCfg *serverConfig.TRedisConfig, MgoCfg *serverConfig.TMgoConfig, isOpenLoop bool) {
	this.rconn = RedisService.NewRedisConn(RedisCfg.ConnAddr, RedisCfg.DBIndex, RedisCfg.Passwd)
	this.mconn = MgoService.NewMgoConn(MgoCfg.UserName, MgoCfg.Passwd, MgoCfg.Host)

	if isOpenLoop {
		this.ctx, this.cancle = context.WithCancel(context.Background())
		this.wg.Add(1)
		go this.LoopSyncDBHard(&this.wg)
		this.wg.Wait()
	}
}

func (this *TDBProvider) LoopSyncDBHard(wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(time.Duration(ado.EDB_DATA_SAVE_INTERVAL) * time.Second)
	for {
		select {
		case <-this.ctx.Done():
			return
		case <-ticker.C:
			// do something...
			this.SyncDBHard()
		default:
			// nothing...
		}
	}
}

func (this *TDBProvider) SyncDBHard() {
	// TODO: Presist redis...
	if this.rconn == nil {
		return
	}

	DBKey := ":" + this.ServerModel + "_Update_Oper"
	_, err := this.rconn.RedPool.Get().Do("SMEMBERS", DBKey)
	if err != nil {
		Log.Error("DBProvider get redis ", err)
		return
	}

	// TODO: Presist mgo...
	if this.mconn == nil {
		return
	}

	// for _, item := range Members.(public.IDBCache) {
	// 	this.mconn.SaveOne(item)
	// }
}
