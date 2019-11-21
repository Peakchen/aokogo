package service

import (
	"common/Config/serverConfig"
	"common/Log"
	"common/MgoConn"
	"common/RedisConn"
	"common/ado"
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/gomodule/redigo/redis"
)

/*
	db module: a lot of redis sessions
	purpose: 建立指定数量的redis 链接，不同玩家唯一认证与之关联，定时快速写入mgo，保证数据文档安全.
*/

type TClusterDBProvider struct {
	//redConn []*RedisConn.TRedisConn
	redConn     *RedisConn.TRedisConn
	mgoConn     *MgoConn.AokoMgo
	mgoSessions []*mgo.Session
	Server      string
	ctx         context.Context
	cancle      context.CancelFunc
	wg          sync.WaitGroup
}

func (this *TClusterDBProvider) init(Server string, RedisCfg *serverConfig.TRedisConfig, MgoCfg *serverConfig.TMgoConfig) {
	this.Server = Server
	this.mgoSessions = make([]*mgo.Session, ado.EMgo_Thread_Cnt)
}

func (this *TClusterDBProvider) Start(Server string, RedisCfg *serverConfig.TRedisConfig, MgoCfg *serverConfig.TMgoConfig) {
	defer func() {
		http.ListenAndServe(RedisCfg.ConnAddr, nil)
	}()
	this.init(Server, RedisCfg, MgoCfg)
	this.runDBloop(Server, RedisCfg, MgoCfg)
}

func (this *TClusterDBProvider) Exit() {
	if this.mgoConn != nil {
		this.mgoConn.Exit()
	}

	if this.redConn != nil {
		this.redConn.Exit()
	}
}

func (this *TClusterDBProvider) runDBloop(Server string, RedisCfg *serverConfig.TRedisConfig, MgoCfg *serverConfig.TMgoConfig) {
	this.redConn = RedisConn.NewRedisConn(RedisCfg.ConnAddr, RedisCfg.DBIndex, RedisCfg.Passwd)
	this.ctx, this.cancle = context.WithCancel(context.Background())

	this.mgoConn = MgoConn.NewMgoConn(Server, MgoCfg.UserName, MgoCfg.Passwd, MgoCfg.Host)
	session, err := this.mgoConn.GetMgoSession()
	if err != nil {
		Log.Error(err)
		return
	}

	for midx := int32(0); midx < ado.EMgo_Thread_Cnt; midx++ {
		this.mgoSessions[midx] = session.Copy()
	}

	this.wg.Add(1)
	go this.LoopDBUpdate(&this.wg)
	this.wg.Wait()
}

func (this *TClusterDBProvider) LoopDBUpdate(wg *sync.WaitGroup) {
	defer func() {
		this.Exit()
		wg.Done()
	}()

	ticker := time.NewTicker(time.Duration(ado.EDB_DATA_SAVE_INTERVAL) * time.Second)
	for {
		select {
		case <-this.ctx.Done():
			ticker.Stop()
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

	c := this.redConn.RedPool.Get()
	if c == nil {
		Log.Error("redis conn invalid or disconntion.")
		return
	}

	var (
		ridx int32
	)
	for {
		if ridx == ado.EMgo_Thread_Cnt {
			break
		}

		this.dbupdate(ridx, c)
		ridx++
	}
}

func (this *TClusterDBProvider) dbupdate(ridx int32, c redis.Conn) {
	//Log.FmtPrintln("db update idx: ", ridx)
	// TODO: Presist redis...
	if this.redConn == nil {
		Log.Error("redis conn invalid or conn number invalid, info: ", this.redConn, ridx)
		return
	}

	updateidx := strconv.Itoa(int(ridx))
	onekey := RedisConn.ERedScript_Update + updateidx
	members, err := c.Do("HKEYS", onekey)
	if err != nil || members == nil {
		Log.Error("ClusterDBProvider get redis,err: ", err)
		return
	}

	// TODO: Presist mgo...
	mgosession := this.mgoSessions[ridx]
	if mgosession == nil {
		Log.Error("mgoConn invalid or disconntion.")
		return
	}

	for _, item := range members.([]interface{}) {
		dstkey := string(item.([]byte))
		dstval, err := c.Do("GET", dstkey)
		if err != nil {
			Log.Error("get fail, err: ", err)
			continue
		}

		bsdata := bson.Raw{Kind: byte(0), Data: dstval.([]byte)}
		err = MgoConn.Save(mgosession, this.Server, dstkey, bsdata)
		if err != nil {
			Log.Error("mgo update err: ", err)
		}
	}
}
