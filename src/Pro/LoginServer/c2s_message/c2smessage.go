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

package c2s_message

import(
	//"common/C2SMessage"
	//"github.com/gorilla/websocket"
	//"github.com/protobuf/proto"
	"log"
	"net/http"
	//"common/define"
)

const(
	strLogin = "/login" 
	strweChatLogin = "/wechatlogin"
)

type loginHandlerfunc func (http.ResponseWriter, *http.Request)

var (
	m_mapLoginHandler map[string] loginHandlerfunc
)

func Register(key string, f loginHandlerfunc){
	var _, ok = m_mapLoginHandler[key]
	if ok {
		return
	}

	m_mapLoginHandler[key] = f
}

func OnDispatchLoginMessage(key string, w http.ResponseWriter, r *http.Request){
	var cb, ok = m_mapLoginHandler[key]
	if !ok {
		return
	} 

	if cb == nil {
		return
	}

	cb(w,r)
}

func init(){
	Register(strLogin, OnLogin)
	Register(strweChatLogin, OnWechatlogin)
	/* C2SMessage.Register(int32(C2SMessage.GameMessageId_msg_req_login), OnLogin)
	C2SMessage.Register(int32(C2SMessage.GameMessageId_msg_req_wechat_login), OnWechatlogin) */
}

/* login game by account and passwd
*/
/* func OnLogin(basemsg *C2SMessage.C2Sbasemessgae, c* websocket.Conn){
	// todo:
	var msg_login = &C2SMessage.Requestlogin{}
	var pm = proto.Unmarshal(basemsg.Data, msg_login)
	if pm == nil {
		log.Fatal("unmarshal message fail.")
		return
	}

	log.Printf("login: player dbid[%d]",msg_login.Dbid)
}
 */
/* login game by wechat
*/
/* func OnWechatlogin(basemsg *C2SMessage.C2Sbasemessgae, c* websocket.Conn){
	// todo:
	var msg_wechatlogin = &C2SMessage.RequestWechatlogin{}
	var pm = proto.Unmarshal(basemsg.Data, msg_wechatlogin)
	if pm == nil {
		log.Fatal("unmarshal message fail.")
		return
	}

	log.Printf("login: player openid[%d]", msg_wechatlogin.Openid)
} */

func OnLogin(w http.ResponseWriter, r *http.Request){
	log.Println("start login.")
	// todo:
	query := r.URL.Query()
	code := query.Get("code")
	nickName := query.Get("nickName")
	avatarURL := query.Get("avatarURL")
	gender := query.Get("gender")
	log.Println("login code : ", code, " , nickName:", nickName, " , avatarURL:", avatarURL, ", gender:", gender)

	
}

/* login game by wechat
*/
func OnWechatlogin(w http.ResponseWriter, r *http.Request){
	// todo:
	
}