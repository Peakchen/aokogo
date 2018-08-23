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

package C2SMessage

import (
	"github.com/gorilla/websocket"
	"github.com/protobuf/proto"
	"log"
)

type funchandler func(msg *C2Sbasemessgae, c* websocket.Conn)

var(
	messageHandler  map[int32] funchandler
)

func Register(id int32, handler funchandler){
	messageHandler[id] = handler
}

func DispatchMessage(msg []byte, c* websocket.Conn){

	var msg_base = &C2Sbasemessgae{}
	var pm = proto.Unmarshal(msg, msg_base)
	if pm == nil {
		log.Fatal("unmarshal message fail.")
		return
	}

	cb, ok := messageHandler[*msg_base.Baseid]
	if ok {
		cb(msg_base, c)
	}
}

func PostMessage(pb proto.Message, c* websocket.Conn){
	
	msg,err := proto.Marshal(pb)
	if err == nil {
		log.Fatal("Marshal message fail.")
		return
	}

	werr := c.WriteMessage(websocket.BinaryMessage, msg)
	if werr != nil{
		log.Fatal("Write close.")
		return
	}
}