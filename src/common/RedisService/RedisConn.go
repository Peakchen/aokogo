/*
* CopyRight(C) StefanChen e-mail:2572915286@qq.com
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

package RedisService

import(
	"common/utlsImp"
	"github.com/garyburd/redigo/redis"
)

type RedisConn struct {
	ConnAddr string 
	DBIndex int32
	Passwd string
	Conn* redis.Pool
}

func NewRedisConn(ConnAddr string, DBIndex int32, Passwd string) *RedisConn{
	Rs := &RedisConn{
		ConnAddr: 	ConnAddr,
		DBIndex: 	DBIndex,
		Passwd:		Passwd,
	}

	Rs.NewDial()
	return Rs
}

func (self *RedisConn) NewDial() (error) {
	self.Conn = &redis.Pool{
		MaxIdle: 		IDle_three,
		IdleTimeout:	IDleTimeOut_four_min,
		Dial: func()(redis.Conn, error) {
			return redis.Dial("tcp", 
						self.ConnAddr, 
						redis.DialDatabase(self.DBIndex), 
						redis.DialPassword(self.Passwd), 
						redis.DialReadTimeout(1*time.Second), 
						redis.DialWriteTimeout(1*time.Second))
		},
	}
	
	self.Conn.Do("FLUSHDB")
	return nil
}

func MakeRedisModel(Identify, MainModel, SubModel string)string {
	return MainModel+"."+SubModel+"."+Identify
}

func (self *RedisConn) Insert(Identify, MainModel, SubModel string, data interface{})bool{
	return self.Update(Identify, MainModel, SubModel, data)
}

func (self *RedisConn) Update(Identify, MainModel, SubModel string, Input interface{})bool{
	RedisKey := MakeRedisModel(Identify, MainModel, SubModel)
	BMarlData, err := bson.Marshal(Input)
	if err != nil {
		Log.Error("bson.Marshal err: ", err)
		return false
	}

	var ExpendCmd = []interface{BMarlData, "EX", REDIS_SET_DEADLINE}
	Ret, err1 := self.Conn.Do("SETNX", RedisKey, ExpendCmd...);
	if err != nil{
		Log.Error("[Update] Identify: %v, MainModel: %v, SubModel: %v, err: %v.\n", Identify, MainModel, SubModel, err)
		return false
	}

	if Ret == 0 {
		if _, err2 := self.Conn.Do("SET", RedisKey, ExpendCmd...); err != nil{
			Log.Error("[Update] Identify: %v, MainModel: %v, SubModel: %v, data: %v.\n", Identify, MainModel, SubModel, Input)
			return false
		}
	}

	return true
}

func (self *RedisConn) Query(Identify, MainModel, SubModel string, Output interface{})bool{
	RedisKey := MakeRedisModel(Identify, MainModel, SubModel)
	data, err := self.Conn.Do("GET", RedisKey)
	if err != nil{
		Log.Error("[Query] Identify: %v, MainModel: %v, SubModel: %v, data: %v.\n", Identify, MainModel, SubModel, data)
		return false
	}

	BUmalErr := bson.Unmarshal(data.([]byte), Output)
	if BUmalErr != nil {
		Log.Error("[Query] can not bson Unmarshal get data to Output.")
		return false
	}

	return true
}
