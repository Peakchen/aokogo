/*
Copyright (this) <year> <copyright holders>

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

package tcpNet

import (
	"common/Define"
	"common/Log"
	"common/msgProto/MSG_HeartBeat"
	"common/msgProto/MSG_MainModule"
	"common/msgProto/MSG_Server"

	"fmt"
	"net"
	"sync/atomic"
	"time"

	//"common/S2SMessage"
	"context"
	"sync"

	"github.com/golang/protobuf/proto"
	//. "common/Define"
)

type ClientTcpSession struct {
	sync.Mutex

	RemoteAddr string
	isAlive    bool
	// The net connection.
	conn *net.TCPConn
	// Buffered channel of outbound messages.
	send chan []byte
	// send/recv
	sw  sync.WaitGroup
	ctx context.Context
	// receive message call back
	recvCb MessageCb
	// person offline flag
	off chan *ClientTcpSession
	//message pack
	pack IMessagePack
	// session id
	SessionID uint64
	//Dest point
	SvrType Define.ERouteId
	//src point
	RegPoint Define.ERouteId
	//person StrIdentify
	StrIdentify string
}

func NewClientSession(addr string,
	conn *net.TCPConn,
	ctx context.Context,
	SvrType Define.ERouteId,
	newcb MessageCb,
	off chan *ClientTcpSession,
	pack IMessagePack) *ClientTcpSession {
	return &ClientTcpSession{
		RemoteAddr: addr,
		conn:       conn,
		send:       make(chan []byte, maxMessageSize),
		isAlive:    false,
		ctx:        ctx,
		recvCb:     newcb,
		pack:       pack,
		off:        make(chan *ClientTcpSession, maxOfflineSize),
		SvrType:    SvrType,
		//StrIdentify: addr,
	}
}

func (this *ClientTcpSession) Alive() bool {
	return this.isAlive
}

func (this *ClientTcpSession) exit(sw *sync.WaitGroup) {
	if this == nil {
		return
	}

	Log.FmtPrintf("session exit, svr: %v, regpoint: %v, cache size: %v.", this.SvrType, this.RegPoint, len(this.send))
	GClient2ServerSession.RemoveSession(this.RemoteAddr)
	this.isAlive = false
	this.StrIdentify = ""
	this.off <- this
	this.send <- []byte{}
	//close(this.send)
	this.conn.CloseRead()
	this.conn.CloseWrite()
	this.conn.Close()
}

func (this *ClientTcpSession) SetSendCache(data []byte) {
	this.send <- data
}

func (this *ClientTcpSession) heartbeatloop(sw *sync.WaitGroup) {
	defer func() {
		sw.Done()
		this.exit(sw)
	}()

	ticker := time.NewTicker(time.Duration(cstKeepLiveHeartBeatSec) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-this.ctx.Done():
			return
		case <-ticker.C:
			if this.RegPoint == 0 || len(this.StrIdentify) == 0 {
				continue
			}
			sendHeartBeat(this, uint16(this.SvrType))
		}
	}
}

func (this *ClientTcpSession) sendloop(sw *sync.WaitGroup) {
	defer func() {
		sw.Done()
		this.exit(sw)
	}()

	for {
		select {
		case <-this.ctx.Done():
			return
		case data := <-this.send:
			if !this.WriteMessage(data) {
				return
			}
		}
	}
}

func (this *ClientTcpSession) recvloop(sw *sync.WaitGroup) {
	defer func() {
		sw.Done()
		this.exit(sw)
	}()

	for {
		select {
		case <-this.ctx.Done():
			return
		default:
			if !this.readMessage() {
				return
			}
		}
	}
}

func (this *ClientTcpSession) WriteMessage(data []byte) (succ bool) {
	if !this.isAlive || len(data) == 0 {
		return
	}

	defer catchRecover()

	this.conn.SetWriteDeadline(time.Now().Add(writeWait))
	//send...
	Log.FmtPrintln("[client] begin send response message to server, message length: ", len(data))
	_, err := this.conn.Write(data)
	if err != nil {
		Log.FmtPrintln("send data fail, err: ", err)
		return false
	}

	return true
}

func (this *ClientTcpSession) readMessage() (succ bool) {
	defer catchRecover()

	//this.conn.SetReadDeadline(time.Now().Add(pongWait))
	if len(this.StrIdentify) == 0 &&
		(this.SvrType == Define.ERouteId_ER_ESG || this.SvrType == Define.ERouteId_ER_ISG) {
		succ = UnPackExternalMsg(this.conn, this.pack)
		if !succ {
			return
		}
		this.pack.SetRemoteAddr(this.RemoteAddr)
	} else {
		succ = UnPackInnerMsg(this.conn, this.pack)
		if !succ {
			return
		}
		this.StrIdentify = this.pack.GetIdentify()
	}

	route := this.pack.GetRouteID()
	mainID, _ := this.pack.GetMessageID()
	if (mainID == uint16(MSG_MainModule.MAINMSG_SERVER) ||
		mainID == uint16(MSG_MainModule.MAINMSG_LOGIN)) && len(this.StrIdentify) == 0 {
		this.StrIdentify = this.RemoteAddr
	}

	if len(this.pack.GetIdentify()) == 0 {
		this.pack.SetIdentify(this.StrIdentify)
	}

	if mainID != uint16(MSG_MainModule.MAINMSG_SERVER) &&
		mainID != uint16(MSG_MainModule.MAINMSG_HEARTBEAT) &&
		(this.SvrType == Define.ERouteId_ER_ISG) {
		Log.FmtPrintf("[client] Route (%v), StrIdentify: %v.", route, this.StrIdentify)
		succ = innerMsgRouteAct(route, mainID, this.pack.GetSrcMsg())
	} else {
		succ = this.checkmsgProc(route) //路由消息回调处理
	}
	return
}

func (this *ClientTcpSession) checkRegisterRet(route uint16) (exist bool) {
	mainID, subID := this.pack.GetMessageID()
	if mainID == uint16(MSG_MainModule.MAINMSG_SERVER) &&
		uint16(MSG_Server.SUBMSG_SC_ServerRegister) == subID {
		this.StrIdentify = this.RemoteAddr
		if this.SvrType == Define.ERouteId_ER_ISG {
			this.Push(Define.ERouteId_ER_ESG)
		} else {
			this.Push(Define.ERouteId(route))
		}

		exist = true
	}
	return
}

func (this *ClientTcpSession) checkHeartBeatRet() (exist bool) {
	mainID, subID := this.pack.GetMessageID()
	if mainID == uint16(MSG_MainModule.MAINMSG_HEARTBEAT) &&
		uint16(MSG_HeartBeat.SUBMSG_SC_HeartBeat) == subID {
		exist = true
	}
	return
}

func (this *ClientTcpSession) checkmsgProc(route uint16) (succ bool) {
	Log.FmtPrintf("recv response, route: %v.", route)
	bRegister := this.checkRegisterRet(route)
	bHeartBeat := checkHeartBeatRet(this.pack)
	if bRegister || bHeartBeat {
		succ = true
		return
	}

	succ = msgCallBack(this)
	return
}

func (this *ClientTcpSession) GetPack() (obj IMessagePack) {
	return this.pack
}

func (this *ClientTcpSession) HandleSession(sw *sync.WaitGroup) {
	this.isAlive = true
	atomic.AddUint64(&this.SessionID, 1)
	Log.FmtPrintln("[client] handle new session: ", this.SessionID)
	sw.Add(3)
	go this.recvloop(sw)
	go this.sendloop(sw)
	go this.heartbeatloop(sw)
}

func (this *ClientTcpSession) Push(RegPoint Define.ERouteId) {
	Log.FmtPrintf("[client] push new sesson, reg point: %v.", RegPoint)
	this.RegPoint = RegPoint
	GServer2ServerSession.AddSession(this.RemoteAddr, this)
}

func (this *ClientTcpSession) SetIdentify(StrIdentify string) {
	session := GServer2ServerSession.GetSessionByIdentify(this.StrIdentify)
	if session != nil {
		GServer2ServerSession.RemoveSession(this.StrIdentify)
		this.StrIdentify = StrIdentify
		GServer2ServerSession.AddSession(StrIdentify, session)
	}
}

func (this *ClientTcpSession) Offline() {

}

func (this *ClientTcpSession) SendMsg(route, mainid, subid uint16, msg proto.Message) (succ bool, err error) {
	if !this.isAlive {
		err = fmt.Errorf("[client] session disconnection, route: %v, mainid: %v, subid: %v.", route, mainid, subid)
		Log.FmtPrintln("send msg err: ", err)
		return succ, err
	}

	data, err := this.pack.PackMsg4Client(route,
		mainid,
		subid,
		msg)
	if err != nil {
		return succ, err
	}
	this.SetSendCache(data)
	return true, nil
}

func (this *ClientTcpSession) SendSvrMsg(route, mainid, subid uint16, msg proto.Message) (succ bool, err error) {
	if !this.isAlive {
		err = fmt.Errorf("[client] session disconnection, mainid: %v, subid: %v.", mainid, subid)
		Log.FmtPrintln("send msg err: ", err)
		return false, err
	}

	data, err := this.pack.PackMsg(route,
		mainid,
		subid,
		msg)

	if err != nil {
		return succ, err
	}
	this.SetSendCache(data)
	return true, nil
}

func (this *ClientTcpSession) SendInnerMsg(identify string, mainid, subid uint16, msg proto.Message) (succ bool, err error) {
	if !this.isAlive {
		err = fmt.Errorf("[client] session disconnection, mainid: %v, subid: %v.", mainid, subid)
		Log.FmtPrintln("send msg err: ", err)
		return false, err
	}

	if len(identify) > 0 {
		this.pack.SetIdentify(identify)
	}

	data, err := this.pack.PackMsg(uint16(Define.ERouteId_ER_ISG),
		mainid,
		subid,
		msg)

	if err != nil {
		return succ, err
	}
	this.SetSendCache(data)
	return true, nil
}

func (this *ClientTcpSession) GetIdentify() string {
	return this.StrIdentify
}

func (this *ClientTcpSession) GetRegPoint() (RegPoint Define.ERouteId) {
	return this.RegPoint
}

func (this *ClientTcpSession) GetRemoteAddr() string {
	return this.RemoteAddr
}
