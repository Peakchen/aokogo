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

package RedisEx

import(
	"common/utlsImp"
	"github.com/garyburd/redigo/redis"
)

var(
	m_RedisInfo = &RedisInfo{nServerPort: 443, szRedisServerHost:":6379", szServerUUID: utlsImp.GetUUID()  }
)

func GetRedisInfo() *RedisInfo{
	return m_RedisInfo
}

func SetRedisInfo(nPort int32, szServerHost string, szUUID string){
	m_RedisInfo.nServerPort = nPort
	m_RedisInfo.szRedisServerHost = szServerHost
	m_RedisInfo.szServerUUID = szUUID 
}

func DialDefaultServer() (redis.Conn, error) {
	c, err := redis.Dial("tcp", Addr, redis.DialReadTimeout(1*time.Second), redis.DialWriteTimeout(1*time.Second))
	if err != nil {
		return nil, err
	}
	
	c.Do("FLUSHDB")
	return c, nil
}

type RedisITF struct {

}

type RedisMgr struct {
	conn *redis.Conn
}

func (self *RedisMgr) Insert(data interface{}, key string){
	if _, err := self.conn.Do("HSET", key, data); err != nil{

	}
}

func (self *RedisMgr) Update(data interface{}, key string){

}
