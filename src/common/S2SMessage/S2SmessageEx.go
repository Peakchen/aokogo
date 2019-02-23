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

package S2SMessage

import (
	"github.com/gorilla/websocket"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
)

type S2SDispatchMHandler func(msg []byte, c net.Conn)
type C2SRoutehandler func(msg *C2SMsgRoute, c* websocket.Conn)

var(
	messageHandler  map[int32] C2SRoutehandler
	s2sdispatchMsgMap  map[int32] S2SDispatchMHandler 
)

func S2SRouteRegister(id int32, handler C2SRoutehandler){
	messageHandler[id] = handler
}

func S2SDispatchMRegister(id int32, handler S2SDispatchMHandler){
	s2sdispatchMsgMap[id] = handler
}

func DispatchClientMessage(msg []byte, c* websocket.Conn){

	var msg_base = &C2SMsgRoute{}
	var pm = proto.Unmarshal(msg, msg_base)
	if pm == nil {
		log.Fatal("unmarshal message fail.")
		return
	}

	cb, ok := messageHandler[*msg_base.Operid]
	if ok {
		cb(msg_base, c)
	}
}

func DispatchMessage(msg []byte, c net.Conn){

	var msg_base = &C2SMsgRoute{}
	var pm = proto.Unmarshal(msg, msg_base)
	if pm == nil {
		log.Fatal("unmarshal message fail.")
		return
	}

	// cb, ok := messageHandler[*msg_base.Operid]
	// if ok {
	// 	cb(msg_base, c)
	// }
}

func PostMessage(pb proto.Message, c net.Conn){
	
	msg, err := proto.Marshal(pb)
	if err == nil {
		log.Fatal("Marshal message fail.")
		return
	}

	_, err = c.Write(msg)
	if err != nil {
		log.Fatal("Write close.")
		return
	}
}