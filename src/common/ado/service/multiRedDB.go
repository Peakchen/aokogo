package service

import (
	"common/Config/serverConfig"
	"common/Log"
	"common/MgoConn"
	"common/RedisConn"
	"common/ado"
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/globalsign/mgo/bson"
)

/*
	db module: a lot of redis sessions
	purpose: 建立指定数量的redis 链接，不同玩家唯一认证与之关联，定时快速写入mgo，保证数据文档安全.
*/

type TClusterDBProvider struct {
	redConn []*RedisConn.TRedisConn
	mgoConn *MgoConn.AokoMgo
	Server  string
	ctx     context.Context
	cancle  context.CancelFunc
	wg      sync.WaitGroup
}

func (this *TClusterDBProvider) init(Server string, RedisCfg *serverConfig.TRedisConfig, MgoCfg *serverConfig.TMgoConfig) {
	this.Server = Server
	this.redConn = []*RedisConn.TRedisConn{}
	this.mgoConn = MgoConn.NewMgoConn(Server, MgoCfg.UserName, MgoCfg.Passwd, MgoCfg.Host)
}

func (this *TClusterDBProvider) Start(Server string, RedisCfg *serverConfig.TRedisConfig, MgoCfg *serverConfig.TMgoConfig) {
	this.init(Server, RedisCfg, MgoCfg)
	this.runDBloop(RedisCfg)
}

func (this *TClusterDBProvider) runDBloop(RedisCfg *serverConfig.TRedisConfig) {
	var (
		cnt int32
	)
	for {
		if cnt >= ado.EMgo_Thread_Cnt {
			break
		}

		cnt++
		rc := RedisConn.NewRedisConn(RedisCfg.ConnAddr, RedisCfg.DBIndex, RedisCfg.Passwd)
		if rc != nil {
			this.redConn = append(this.redConn, rc)
		}
	}

	this.ctx, this.cancle = context.WithCancel(context.Background())
	this.wg.Add(1)
	go this.LoopDBUpdate(&this.wg)
	this.wg.Wait()

}

func (this *TClusterDBProvider) LoopDBUpdate(wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(time.Duration(ado.EDB_DATA_SAVE_INTERVAL) * time.Second)
	for {
		select {
		case <-this.ctx.Done():
			return
		case <-ticker.C:
			// do something...
			this.flushdb()
		default:
			// nothing...
		}
	}
}

func (this *TClusterDBProvider) flushdb() {
	var (
		ridx int32
	)
	for {
		if ridx == ado.EMgo_Thread_Cnt {
			break
		}

		ridx++
		this.dbupdate(ridx)
	}
}

func (this *TClusterDBProvider) dbupdate(ridx int32) {
	//Log.FmtPrintln("db update idx: ", ridx)
	// TODO: Presist redis...
	if this.redConn == nil || len(this.redConn) < int(ridx) {
		Log.Error("redis conn invalid or conn number invalid, info: ", this.redConn, len(this.redConn), ridx)
		return
	}

	updateidx := strconv.Itoa(int(ridx))
	onekey := RedisConn.ERedScript_Update + updateidx
	c := this.redConn[ridx-1].RedPool.Get()
	if c == nil {
		Log.Error("redis invalid or disconntion, redis conn idx: ", ridx)
		return
	}

	members, err := c.Do("HKEYS", onekey)
	if err != nil || members == nil {
		Log.Error("ClusterDBProvider get redis ", err)
		return
	}

	// TODO: Presist mgo...
	if this.mgoConn == nil {
		Log.Error("mgoConn invalid or disconntion.")
		return
	}

	for _, item := range members.([]interface{}) {
		dstkey := string(item.([]byte))
		dstval, err := c.Do("GET", dstkey)
		if err != nil {
			Log.FmtPrintln("get fail, err: ", err)
			continue
		}

		bsdata := bson.Raw{Kind: byte(0), Data: dstval.([]byte)}
		err = this.mgoConn.Save(dstkey, bsdata)
		if err != nil {
			Log.Error("mgo update err: ", err)
		}
	}
}
