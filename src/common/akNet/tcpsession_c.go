// add by stefan

package akNet

import (
	"common/Define"
	"common/Log"
	"common/aktime"
	"common/msgProto/MSG_HeartBeat"
	"common/msgProto/MSG_MainModule"
	"common/msgProto/MSG_Server"
	"common/stacktrace"
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
	//
	Name string
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
		off:        off,
		SvrType:    SvrType,
		//StrIdentify: addr,
	}
}

func (this *ClientTcpSession) Alive() bool {
	return this.isAlive
}

func (this *ClientTcpSession) close(sw *sync.WaitGroup) {
	if this == nil {
		return
	}

	Log.FmtPrintf("session close, svr: %v, regpoint: %v, cache size: %v.", this.SvrType, this.RegPoint, len(this.send))
	GClient2ServerSession.RemoveSession(this.RemoteAddr)
	this.off <- this
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
		this.close(sw)
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
			sendHeartBeat(this)
		}
	}
}

func (this *ClientTcpSession) sendloop(sw *sync.WaitGroup) {
	defer func() {
		sw.Done()
		this.close(sw)
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
		this.close(sw)
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

	defer stacktrace.Catchcrash()

	this.conn.SetWriteDeadline(aktime.Now().Add(writeWait))
	//send...
	//Log.FmtPrintln("[client] begin send response message to server, message length: ", len(data))
	_, err := this.conn.Write(data)
	if err != nil {
		Log.FmtPrintln("send data fail, err: ", err)
		return false
	}

	return true
}

func (this *ClientTcpSession) readMessage() (succ bool) {
	defer func() {
		this.Unlock()
		stacktrace.Catchcrash()
	}()

	this.Lock()

	//this.conn.SetReadDeadline(aktime.Now().Add(pongWait))
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

	var route Define.ERouteId
	mainID, _ := this.pack.GetMessageID()
	if (mainID == uint16(MSG_MainModule.MAINMSG_SERVER) ||
		mainID == uint16(MSG_MainModule.MAINMSG_LOGIN)) && len(this.StrIdentify) == 0 {
		this.StrIdentify = this.RemoteAddr
	}

	if len(this.pack.GetIdentify()) == 0 {
		this.pack.SetIdentify(this.StrIdentify)
	}

	if mainID == uint16(MSG_MainModule.MAINMSG_LOGIN) {
		route = Define.ERouteId_ER_Login
	} else if mainID >= uint16(MSG_MainModule.MAINMSG_PLAYER) {
		route = Define.ERouteId_ER_Game
	}

	if mainID != uint16(MSG_MainModule.MAINMSG_SERVER) &&
		mainID != uint16(MSG_MainModule.MAINMSG_HEARTBEAT) &&
		(this.SvrType == Define.ERouteId_ER_ISG) {
		//Log.FmtPrintf("[client] StrIdentify: %v.", this.StrIdentify)
		succ = innerMsgRouteAct(ESessionType_Client, route, mainID, this.pack.GetSrcMsg())
	} else {
		succ = this.checkmsgProc(route) //路由消息回调处理
	}
	return
}

func (this *ClientTcpSession) checkRegisterRet(route Define.ERouteId) (exist bool) {
	mainID, subID := this.pack.GetMessageID()
	if mainID == uint16(MSG_MainModule.MAINMSG_SERVER) &&
		uint16(MSG_Server.SUBMSG_SC_ServerRegister) == subID {
		this.StrIdentify = this.RemoteAddr
		if this.SvrType == Define.ERouteId_ER_ISG {
			this.Push(Define.ERouteId_ER_ESG)
		} else {
			this.Push(Define.ERouteId_ER_ISG)
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

func (this *ClientTcpSession) checkmsgProc(route Define.ERouteId) (succ bool) {
	//Log.FmtPrintf("recv response, route: %v.", route)
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

	this.Name = fmt.Sprintf("client_%v_%v", GetModuleDef(this.SvrType), this.SessionID)
}

func (this *ClientTcpSession) Push(RegPoint Define.ERouteId) {
	//Log.FmtPrintf("[client] push new sesson, reg point: %v.", RegPoint)
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

func (this *ClientTcpSession) SendMsg(mainid, subid uint16, msg proto.Message) (succ bool, err error) {
	if !this.isAlive {
		err = fmt.Errorf("[client] send msg session disconnection, mainid: %v, subid: %v.", mainid, subid)
		Log.FmtPrintln("send msg err: ", err)
		return succ, err
	}

	data, err := this.pack.PackMsg4Client(mainid, subid, msg)
	if err != nil {
		return succ, err
	}
	this.SetSendCache(data)
	return true, nil
}

func (this *ClientTcpSession) SendSvrMsg(mainid, subid uint16, msg proto.Message) (succ bool, err error) {
	if !this.isAlive {
		err = fmt.Errorf("[client] send svr session disconnection, mainid: %v, subid: %v.", mainid, subid)
		Log.FmtPrintln("send msg err: ", err)
		return false, err
	}

	data, err := this.pack.PackMsg(mainid, subid, msg)
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

	data, err := this.pack.PackMsg(mainid, subid, msg)
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

func (this *ClientTcpSession) IsUser() bool {
	return this.RegPoint == 0
}

func (this *ClientTcpSession) RefreshHeartBeat(mainid, subid uint16) bool {
	return true
}

func (this *ClientTcpSession) GetModuleName() string {
	return this.Name
}
