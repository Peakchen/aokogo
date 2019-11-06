/*
Copyright (c) <year> <copyright holders>

"Anti 996" License Version 1.0 (Draft)

Permission is hereby granted to any individual or legal entity
obtaining a copy of this licensed work (including the source code,
documentation and/or related items, hereinafter collectively referred
to as the "licensed work"), free of charge, to deal with the licensed
work for any purpose, including without limitation, the rights to use,
reproduce, modify, prepare derivative works of, distribute, publish
and sublicense the licensed work, subject to the following conditions:

1. The individual or the legal entity must conspicuously display,
without modification, this License and the notice on each redistributed
or derivative copy of the Licensed Work.

2. The individual or the legal entity must strictly comply with all
applicable laws, regulations, rules and standards of the jurisdiction
relating to labor and employment where the individual is physically
located or where the individual was born or naturalized; or where the
legal entity is registered or is operating (whichever is stricter). In
case that the jurisdiction has no such laws, regulations, rules and
standards or its laws, regulations, rules and standards are
unenforceable, the individual or the legal entity are required to
comply with Core International Labor Standards.

3. The individual or the legal entity shall not induce, metaphor or force
its employee(s), whether full-time or part-time, or its independent
contractor(s), in any methods, to agree in oral or written form, to
directly or indirectly restrict, weaken or relinquish his or her
rights or remedies under such laws, regulations, rules and standards
relating to labor and employment as mentioned above, no matter whether
such written or oral agreement are enforceable under the laws of the
said jurisdiction, nor shall such individual or the legal entity
limit, in any methods, the rights of its employee(s) or independent
contractor(s) from reporting or complaining to the copyright holder or
relevant authorities monitoring the compliance of the license about
its violation(s) of the said license.

THE LICENSED WORK IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE COPYRIGHT HOLDER BE LIABLE FOR ANY CLAIM,
DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
OTHERWISE, ARISING FROM, OUT OF OR IN ANY WAY CONNECTION WITH THE
LICENSED WORK OR THE USE OR OTHER DEALINGS IN THE LICENSED WORK.
*/

package RedisConn

import (
	"common/Log"
	"common/ado"
	"common/public"
	"fmt"
	"strconv"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/gomodule/redigo/redis"
)

type TRedisConn struct {
	ConnAddr string
	DBIndex  int32
	Passwd   string
	RedPool  *redis.Pool
	us       *TRedisScript
}

func NewRedisConn(ConnAddr string, DBIndex int32, Passwd string) *TRedisConn {
	Rs := &TRedisConn{
		ConnAddr: ConnAddr,
		DBIndex:  DBIndex,
		Passwd:   Passwd,
	}

	Rs.us = &TRedisScript{
		name:   ERedScript_Update,
		script: updateScript(),
	}
	Rs.NewDial()
	return Rs
}

func (this *TRedisConn) NewDial() error {
	this.RedPool = &redis.Pool{
		MaxIdle:     IDle_three,
		IdleTimeout: IDleTimeOut_four_min,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp",
				this.ConnAddr,
				redis.DialDatabase(int(this.DBIndex)),
				redis.DialPassword(this.Passwd),
				redis.DialReadTimeout(1*time.Second),
				redis.DialWriteTimeout(1*time.Second))
		},
	}
	this.RedPool.Get().Do("FLUSHDB")
	return nil
}

/*
	Redis Oper func: Insert
	SaveType: EDBOper_Insert
	purpose: in order to Insert data type EDBOperType to Redis Cache.
*/
func (this *TRedisConn) Insert(Input public.IDBCache) error {
	return this.Update(Input, ado.EDBOper_Insert)
}

/*
	Redis Oper func: Update
	SaveType: EDBOper_Update
	purpose: in order to Update data type EDBOperType to Redis Cache.
*/
func (this *TRedisConn) Update(Input public.IDBCache, SaveType ado.EDBOperType) (err error) {
	RedisKey := MakeRedisModel(Input.Identify(), Input.MainModel(), Input.SubModel())
	BMarlData, err := bson.Marshal(Input)
	if err != nil {
		err = fmt.Errorf("bson.Marshal err: %v.\n", err)
		Log.Error("[Update] err: %v", err)
		return
	}

	this.SaveEx(Input.Identify(), RedisKey, BMarlData, SaveType)
	return
}

/*
	Redis Oper func: Query
	purpose: in order to Get data from Redis Cache.
*/
func (this *TRedisConn) Query(Output public.IDBCache) (ret error) {
	ret = nil
	RedisKey := MakeRedisModel(Output.Identify(), Output.MainModel(), Output.SubModel())
	data, err := this.RedPool.Get().Do("GET", RedisKey)
	if err != nil {
		ret = fmt.Errorf("Identify: %v, MainModel: %v, SubModel: %v, data: %v.\n", Output.Identify(), Output.MainModel(), Output.SubModel(), data)
		Log.Error("[Query] err: %v.\n", ret)
		return
	}

	if data == nil {
		ret = fmt.Errorf("Identify: %v, MainModel: %v, SubModel: %v, Nil data is invalid.\n", Output.Identify(), Output.MainModel(), Output.SubModel())
		Log.Error("[Query] err: %v.\n", ret)
		return
	}

	BUmalErr := bson.Unmarshal(data.([]byte), Output)
	if BUmalErr != nil {
		ret = fmt.Errorf("Identify: %v, MainModel: %v, SubModel: %v, data: %v.\n", Output.Identify(), Output.MainModel(), Output.SubModel(), data)
		Log.Error("[Query] can not bson Unmarshal get data to Output, err: %v.\n", ret)
		return
	}

	return
}

func (this *TRedisConn) Save(rolekey, RedisKey string, data interface{}, SaveType ado.EDBOperType) (ret error) {
	ret = nil
	switch SaveType {
	case ado.EDBOper_Insert:
		ExpendCmd := []interface{}{RedisKey, data, "EX", REDIS_SET_DEADLINE}
		Ret, err := this.RedPool.Get().Do("SETNX", ExpendCmd...) // set if not exist
		if err != nil {
			Log.Error("[Save] SETNX data: %v, err: %v.\n", data, err)
			return
		}

		if Ret == 0 {
			// connect key and value.
			if _, err := this.RedPool.Get().Do("SET", ExpendCmd...); err != nil {
				Log.Error("[Save] Insert SET data: %v, err: %v..\n", data, err)
				return
			}
		}

	case ado.EDBOper_Update:
		// connect key and value.
		var ExpendCmd = []interface{}{RedisKey, data, "EX", REDIS_SET_DEADLINE}
		if _, err := this.RedPool.Get().Do("SET", ExpendCmd...); err != nil {
			Log.Error("[Save] Update Set data: %v, err: %v.\n", data, err)
			return
		}

		CollectKey := ":" + RedisKey + "_Update_Oper"
		// Add to collection.
		if _, err := this.RedPool.Get().Do("SADD", CollectKey, RedisKey); err != nil {
			Log.Error("[Save] SADD CollectKey: %v, RedisKey: %v, err: %v.", CollectKey, RedisKey, err)
			return
		}

	case ado.EDBOper_Delete:
		// nothing...
	case ado.EDBOper_DB: //it can be presisted to database.
		// for mogo db.
	default:
		// nothing...

	}

	return
}

func (this *TRedisConn) SaveEx(rolekey, RedisKey string, data interface{}, SaveType ado.EDBOperType) (ret error) {
	ret = this.redSetAct(rolekey, RedisKey, data, false, 0)
	return
}

func (this *TRedisConn) redSetAct(key string, fieldkey string, data interface{}, bsetEx bool, extime int32) (err error) {
	nhashk := RoleKey2Haskey(key)
	strkey := ERedScript_Update + strconv.Itoa(nhashk)
	Log.FmtPrintf("redis act, hashKey: %v, fieldkey: %v.", strkey, fieldkey)
	c := this.RedPool.Get()
	if c == nil {
		err = Log.RetError("red pool get session fail.")
		return
	}

	var (
		exCmd = []interface{}{}
	)

	if bsetEx {
		exCmd = []interface{}{strkey, fieldkey, 1, data, "EX", extime}
	} else {
		exCmd = []interface{}{strkey, fieldkey, 1, data, "", 0}
	}

	_, err = this.us.script.Do(c, exCmd...)
	if err != nil {
		err = Log.RetError("name: %v, ex cmd %v, err: %v", this.us.name, exCmd, err)
		return
	}

	err = nil
	return
}
