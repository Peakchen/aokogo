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
	"log"
	"time"
)

func NewPool(addr string)* redis.Pool{
	return &redis.Pool{
		MaxIdle: 		IDle_three,
		IdleTimeout:	IDleTimeOut_four_min,
		Dial: func()(redis.Conn, error) {return redis.Dial("tcp", addr)},
	}
}

func Subscriber(pool *redis.Pool, subcontent interface{}){
	if pool == nil{
		return
	}

	for{
		c := pool.Get()
		psc := redis.PubSubConn{Conn: c}
		psc.Subscribe(subcontent)

		switch recv := psc.Receive().(type) {
		case redis.Message:
			log.Printf("channel: %s, message: %s.", recv.Channel, recv.Data)
		case redis.Subscription:
			log.Printf("channel: %s, Count: %d, kind: %s.", recv.Channel, recv.Count, recv.Kind)
		case redis.PMessage:
			log.Printf("channel: %s, message: %s, Pattern: %s.", recv.Channel, recv.Data, recv.Pattern)
		case error:
			log.Printf("Error.")
			c.Close()
			time.Sleep(Sec_five)
			break
		default:
			log.Printf("default noting.")
		}
	}
}
