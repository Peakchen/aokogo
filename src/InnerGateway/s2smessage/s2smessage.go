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

package s2smessage

import(
	"common/S2SMessage"
	"log"
	//"net"
	"github.com/gorilla/websocket"
)


func init(){
	S2SMessage.S2SRouteRegister(int32(S2SMessage.ServerId_SID_ISG), OnS2SDispatchMessage)

	
}

/* dispatch s2s message
*/
func OnS2SDispatchMessage(routemsg *S2SMessage.C2SMsgRoute, c* websocket.Conn){
	// todo:
	log.Printf("OnS2SDispatchMessage: Operid[%d].", routemsg.Operid)

/* 	cb, ok := messageHandler[*msg_base.Operid]
	if ok {
		cb(msg_base, c)
	} */
}
