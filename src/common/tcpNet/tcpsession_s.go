// add by stefan

package tcpNet

import (
	"common/Define"
	"common/Log"
	"common/msgProto/MSG_MainModule"

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

type SvrTcpSession struct {
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
	off chan *SvrTcpSession
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

func NewSvrSession(addr string,
	conn *net.TCPConn,
	ctx context.Context,
	SvrType Define.ERouteId,
	newcb MessageCb,
	off chan *SvrTcpSession,
	pack IMessagePack) *SvrTcpSession {
	return &SvrTcpSession{
		RemoteAddr: addr,
		conn:       conn,
		send:       make(chan []byte, maxMessageSize),
		isAlive:    false,
		ctx:        ctx,
		recvCb:     newcb,
		pack:       pack,
		off:        make(chan *SvrTcpSession, maxOfflineSize),
		SvrType:    SvrType,
	}
}

func (this *SvrTcpSession) Alive() bool {
	return this.isAlive
}

func (this *SvrTcpSession) exit(sw *sync.WaitGroup) {
	if this == nil {
		return
	}

	Log.FmtPrintf("session exit, svr: %v, regpoint: %v, cache size: %v.", this.SvrType, this.RegPoint, len(this.send))
	GServer2ServerSession.RemoveSession(this.RemoteAddr)
	this.isAlive = false
	this.StrIdentify = ""
	this.off <- this
	this.send <- []byte{}
	//close(this.send)
	this.conn.CloseRead()
	this.conn.CloseWrite()
	this.conn.Close()
}

func (this *SvrTcpSession) SetSendCache(data []byte) {
	this.send <- data
}

func (this *SvrTcpSession) sendloop(sw *sync.WaitGroup) {
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

func (this *SvrTcpSession) recvloop(sw *sync.WaitGroup) {
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

func (this *SvrTcpSession) WriteMessage(data []byte) (succ bool) {
	if !this.isAlive || len(data) == 0 {
		return
	}

	defer catchRecover()

	this.conn.SetWriteDeadline(time.Now().Add(writeWait))
	//send...
	Log.FmtPrintln("[server] begin send response message to client, message length: ", len(data))
	_, err := this.conn.Write(data)
	if err != nil {
		Log.FmtPrintln("send data fail, err: ", err)
		return false
	}

	return true
}

func (this *SvrTcpSession) readMessage() (succ bool) {
	defer catchRecover()

	//this.conn.SetReadDeadline(time.Now().Add(pongWait))
	if this.RegPoint == 0 {
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
	mainID, SubID := this.pack.GetMessageID()
	Log.FmtPrintf("recv message, route: %v, mainID: %v, subID: %v.", route, mainID, SubID)
	if mainID == uint16(MSG_MainModule.MAINMSG_SERVER) &&
		Define.ERouteId(route) == Define.ERouteId_ER_ISG &&
		this.SvrType == Define.ERouteId_ER_ESG {

		this.RegPoint = Define.ERouteId_ER_ISG
		this.Push(Define.ERouteId(route)) //外网关加入内网关session
		RegisterMessageRet(this, uint16(Define.ERouteId_ER_ESG))
		succ = true
		return
	}

	if (mainID == uint16(MSG_MainModule.MAINMSG_SERVER) ||
		mainID == uint16(MSG_MainModule.MAINMSG_LOGIN)) && len(this.StrIdentify) == 0 {
		this.StrIdentify = this.RemoteAddr
	}

	if len(this.pack.GetIdentify()) == 0 {
		this.pack.SetIdentify(this.StrIdentify)
	}

	if mainID != uint16(MSG_MainModule.MAINMSG_SERVER) && mainID != uint16(MSG_MainModule.MAINMSG_HEARTBEAT) &&
		(this.SvrType == Define.ERouteId_ER_ESG || this.SvrType == Define.ERouteId_ER_ISG) {
		Log.FmtPrintf("[server] Route (%v), StrIdentify: %v.", route, this.StrIdentify)
		if this.SvrType == Define.ERouteId_ER_ESG {
			succ = externalRouteAct(route, this)
		} else {
			succ = innerMsgRouteAct(route, mainID, this.pack.GetSrcMsg())
		}
	} else {
		succ = msgCallBack(this) //路由消息回调处理
	}
	return
}

func (this *SvrTcpSession) GetPack() (obj IMessagePack) {
	return this.pack
}

func (this *SvrTcpSession) HandleSession(sw *sync.WaitGroup) {
	this.isAlive = true
	atomic.AddUint64(&this.SessionID, 1)
	Log.FmtPrintln("[server] handle new session: ", this.SessionID)
	sw.Add(2)
	go this.recvloop(sw)
	go this.sendloop(sw)
}

func (this *SvrTcpSession) Push(RegPoint Define.ERouteId) {
	Log.FmtPrintf("[server] push new sesson, reg point: %v.", RegPoint)
	this.RegPoint = RegPoint
	GServer2ServerSession.AddSession(this.RemoteAddr, this)
}

func (this *SvrTcpSession) SetIdentify(StrIdentify string) {
	session := GServer2ServerSession.GetSessionByIdentify(this.StrIdentify)
	if session != nil {
		GServer2ServerSession.RemoveSession(this.StrIdentify)
		this.StrIdentify = StrIdentify
		GServer2ServerSession.AddSession(StrIdentify, session)
	}
}

func (this *SvrTcpSession) Offline() {

}

func (this *SvrTcpSession) SendMsg(route, mainid, subid uint16, msg proto.Message) (succ bool, err error) {
	if !this.isAlive {
		err = fmt.Errorf("[server] session disconnection, route: %v, mainid: %v, subid: %v.", route, mainid, subid)
		Log.FmtPrintln("send msg err: ", err)
		return false, err
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

func (this *SvrTcpSession) SendSvrMsg(route, mainid, subid uint16, msg proto.Message) (succ bool, err error) {
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

func (this *SvrTcpSession) SendInnerMsg(identify string, mainid, subid uint16, msg proto.Message) (succ bool, err error) {
	if !this.isAlive {
		err = fmt.Errorf("[server] session disconnection, mainid: %v, subid: %v.", mainid, subid)
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

func (this *SvrTcpSession) GetIdentify() string {
	return this.StrIdentify
}

func (this *SvrTcpSession) GetRegPoint() (RegPoint Define.ERouteId) {
	return this.RegPoint
}

func (this *SvrTcpSession) GetRemoteAddr() string {
	return this.RemoteAddr
}
