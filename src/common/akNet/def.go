package akNet

import (
	"common/Define"
	"net"
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"
)

// add by stefan

type IMessagePack interface {
	PackAction(Output []byte) (err error)
	PackAction4Client(Output []byte) (err error)
	PackData(msg proto.Message) (data []byte, err error)
	UnPackMsg4Client(InData []byte) (pos int, err error)
	UnPackData() (msg proto.Message, cb reflect.Value, err error, exist bool)
	GetRouteID() (route uint16)
	GetMessageID() (mainID uint16, subID uint16)
	Clean()
	SetCmd(mainid, subid uint16, data []byte)
	PackMsg(mainid, subid uint16, msg proto.Message) (out []byte, err error)
	PackMsg4Client(mainid, subid uint16, msg proto.Message) (out []byte, err error)
	GetSrcMsg() (data []byte)
	SetIdentify(identify string)
	GetIdentify() string
	UnPackMsg4Svr(InData []byte) (pos int, err error)
	GetDataLen() (datalen uint32)
	SetRemoteAddr(addr string)
	GetRemoteAddr() (addr string)
}

/*

- 协议格式，小端字节序

转发目的位置 |主命令号 | 次命令号 | 包长度 | 包体
------------|--------| --------| -------- | ----
2字节   |2字节   |  2字节  | 4字节  | 包体
*/
const (
	//EnMessage_RoutePoint    = 2  //转发位置
	EnMessage_MainIDPackLen = 2 //主命令
	EnMessage_SubIDPackLen  = 2 //次命令
	EnMessage_DataPackLen   = 4 //真实数据长度 (转发位置+主命令+次命令)
	EnMessage_NoDataLen     = 8 //非data数据长度(包体之前的)->(转发位置+主命令+次命令+datalen)

	EnMessage_SvrDataPackLen  = 48 //真实数据长度 (转发位置+主命令+次命令+ Identify长度 + Identify内容+RemoteAddr 长度+RemoteAddr 内容)
	EnMessage_SvrNoDataLen    = 52 //非data数据长度(包体之前的)->(转发位置+主命令+次命令+ Identify长度 + Identify内容+datalen+RemoteAddr 长度+RemoteAddr 内容)
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
	SendSvrMsg(mainid, subid uint16, msg proto.Message) (succ bool, err error)
	SendMsg(mainid, subid uint16, msg proto.Message) (succ bool, err error)
	SendInnerMsg(identify string, mainid, subid uint16, msg proto.Message) (succ bool, err error)
	WriteMessage(data []byte) (succ bool)
	Alive() bool
	GetPack() (obj IMessagePack)
	IsUser() bool
	RefreshHeartBeat(mainid, subid uint16) bool
	GetModuleName() string
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

//client reconnect check interval (s)
const (
	cstClientSessionCheckMs = 5
)

const (
	cstKeepLiveHeartBeatSec     = 180 //180 3min
	cstCheckHeartBeatMonitorSec = cstKeepLiveHeartBeatSec / 2
	cstSvrDisconnectionSec      = 3 * cstKeepLiveHeartBeatSec //s
	cstClientDisconnectionSec   = 6 * cstKeepLiveHeartBeatSec //s
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

var (
	_svrDefs = map[Define.ERouteId]string{
		Define.ERouteId_ER_ESG:        "ExternalGateway",
		Define.ERouteId_ER_ISG:        "InnerGateway",
		Define.ERouteId_ER_DB:         "DB",
		Define.ERouteId_ER_BigWorld:   "BigWorld",
		Define.ERouteId_ER_Login:      "Login",
		Define.ERouteId_ER_SmallWorld: "SmallWorld",
		Define.ERouteId_ER_DBProxy:    "DBProxy",
		Define.ERouteId_ER_Game:       "Game",
		Define.ERouteId_ER_Client:     "Client",
		Define.ERouteId_ER_Max:        "Max",
	}
)

func GetModuleDef(routeid Define.ERouteId) string {
	name, ok := _svrDefs[routeid]
	if !ok {
		name = "Unknow"
	}
	return name
}
