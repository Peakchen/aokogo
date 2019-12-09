package tcpNet

import (
	"common/Define"
	"common/Log"
	"net"
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"
)

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

type IMessagePack interface {
	PackAction(Output []byte) (err error)
	PackAction4Client(Output []byte) (err error)
	PackData(msg proto.Message) (data []byte, err error)
	UnPackMsg4Client(InData []byte) (pos int, err error)
	UnPackData() (msg proto.Message, cb reflect.Value, err error, exist bool)
	GetRouteID() (route uint16)
	GetMessageID() (mainID uint16, subID uint16)
	Clean()
	SetCmd(routepoint, mainid, subid uint16, data []byte)
	PackMsg(routepoint, mainid, subid uint16, msg proto.Message) (out []byte, err error)
	PackMsg4Client(routepoint, mainid, subid uint16, msg proto.Message) (out []byte, err error)
	GetSrcMsg() (data []byte)
	SetIdentify(identify string)
	GetIdentify() string
	UnPackMsg4Svr(InData []byte) (pos int, err error)
	GetDataLen() (datalen uint32)
	SetRemoteAddr(addr string)
	GetRemoteAddr() (addr string)
}

/*
	func: EncodeCmd
	purpose: Encode message mainid and subid to cmd.
*/
func EncodeCmd(mainID, subID uint16) uint32 {
	return (uint32(mainID) << 16) | uint32(subID)
}

/*
	func: DecodeCmd
	purpose: DecodeCmd message cmd to mainid and subid.
*/
func DecodeCmd(cmd uint32) (uint16, uint16) {
	return uint16(cmd >> 16), uint16(cmd)
}

/*

- 协议格式，小端字节序

转发目的位置 |主命令号 | 次命令号 | 包长度 | 包体
------------|--------| --------| -------- | ----
2字节   |2字节   |  2字节  | 4字节  | 包体
*/
const (
	EnMessage_RoutePoint    = 2  //转发位置
	EnMessage_MainIDPackLen = 2  //主命令
	EnMessage_SubIDPackLen  = 2  //次命令
	EnMessage_DataPackLen   = 6  //真实数据长度 (转发位置+主命令+次命令)
	EnMessage_NoDataLen     = 10 //非data数据长度(包体之前的)->(转发位置+主命令+次命令+datalen)

	EnMessage_SvrDataPackLen  = 50 //真实数据长度 (转发位置+主命令+次命令+ Identify长度 + Identify内容+RemoteAddr 长度+RemoteAddr 内容)
	EnMessage_SvrNoDataLen    = 54 //非data数据长度(包体之前的)->(转发位置+主命令+次命令+ Identify长度 + Identify内容+datalen+RemoteAddr 长度+RemoteAddr 内容)
	EnMessage_IdentifyEixst   = 1  //Identify 长度
	EnMessage_IdentifyLen     = 21 //Identify 内容
	EnMessage_RemoteAddrEixst = 1  //RemoteAddr 长度
	EnMessage_RemoteAddrLen   = 21 //RemoteAddr 内容
)

// session, data, data len
type MessageCb func(c net.Conn, mainID uint16, subID uint16, msg proto.Message)

type TcpSession interface {
	GetRemoteAddr() string
	GetRegPoint() (RegPoint Define.ERouteId)
	GetIdentify() string
	SetSendCache(data []byte)
	Push(RegPoint Define.ERouteId)
	SetIdentify(StrIdentify string)
	SendSvrMsg(route, mainid, subid uint16, msg proto.Message) (succ bool, err error)
	SendMsg(route, mainid, subid uint16, msg proto.Message) (succ bool, err error)
	SendInnerMsg(identify string, mainid, subid uint16, msg proto.Message) (succ bool, err error)
	WriteMessage(data []byte) (succ bool)
	Alive() bool
	GetPack() (obj IMessagePack)
}

// after dial connect todo action.
type AfterDialAct func(s TcpSession)

type TConnSession struct {
	Connsess *TcpSession
	Svr      int32
}

type ESessionGetType int

const (
	ESessionGetType_Identify ESessionGetType = 1
	ESessionGetType_RegPoint ESessionGetType = 2
)

type IProcessConnSession interface {
	RemoveSession(key interface{})
	AddSession(key interface{}, session TcpSession)
	GetSession(key interface{}) (session TcpSession)
	GetSessionByIdentify(key interface{}) (session TcpSession)
}

type ESessionType int8

const (
	ESessionType_Server ESessionType = 1
	ESessionType_Client ESessionType = 2
)

//session begin number with 10000
const (
	ESessionBeginNum = uint64(10000)
)

//client reconnect check interval (ms)
const (
	cstClientSessionCheckMs = 5000
)

const (
	cstKeepLiveHeartBeatSec = 10
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 3 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 4096
	//offline session
	maxOfflineSize = 1024
)

func catchRecover() {
	if r := recover(); r != nil {
		Log.Error("catch recover: ", r)
	}
}
