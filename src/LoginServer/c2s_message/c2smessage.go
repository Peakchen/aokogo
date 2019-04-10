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