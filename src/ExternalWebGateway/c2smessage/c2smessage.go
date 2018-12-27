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

package c2smessage

import(
	"common/C2SMessage"
	"common/S2SMessage"
	"github.com/gorilla/websocket"
	"log"
)

/* login game by account and passwd
*/
func OnC2SDispatchMessage(basemsg *C2SMessage.C2Sbasemessgae, c* websocket.Conn){
	// todo:
	log.Printf("OnC2SDispatchMessage: player Sid[%d].", basemsg.Sid)
	switch S2SMessage.ServerId(*basemsg.Sid) {
	case S2SMessage.ServerId_SID_Login:
		//send to loginserver
		S2SMessage.DispatchClientMessage(basemsg.Data, c)
	case S2SMessage.ServerId_SID_BigWorld:

	case S2SMessage.ServerId_SID_Game:

	case S2SMessage.ServerId_SID_SmallWorld:

	}
}

func init(){
	C2SMessage.Register(int32(C2SMessage.MessageRoute_R_SID_ESG), OnC2SDispatchMessage)
}