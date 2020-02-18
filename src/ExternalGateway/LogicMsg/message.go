package LogicMsg

// add by stefan

import (
	"common/Define"
	"common/Log"
	"common/msgProto/MSG_HeartBeat"
	"common/msgProto/MSG_MainModule"
	"common/msgProto/MSG_Server"
	"common/tcpNet"
	"fmt"
	"net"

	"github.com/golang/protobuf/proto"
)

func ExternalGatewayMessageCallBack(c net.Conn, mainID uint16, subID uint16, msg proto.Message) {
	Log.FmtPrintf("exec external gateway server message call back: %v, %v.", c.RemoteAddr(), c.LocalAddr())
}

func onSvrRegister(session tcpNet.TcpSession, req *MSG_Server.CS_ServerRegister_Req) (succ bool, err error) {
	Log.FmtPrintf("onSvrRegister, StrIdentify: %v, recv: %v.", session.GetIdentify(), req.ServerType)
	var (
		msgfmt string
	)

	session.Push(Define.ERouteId(req.ServerType))
	for _, id := range req.Msgs {
		mainid, subid := tcpNet.DecodeCmd(uint32(id))
		msgfmt += fmt.Sprintf("[mainid: %v, subid: %v]\t", mainid, subid)
	}

	msgfmt += "\n"
	Log.FmtPrintln("message context: ", msgfmt)
	return tcpNet.RegisterMessageRet(session, uint16(Define.ERouteId_ER_ESG))
}

func onHeartBeat(session tcpNet.TcpSession, req *MSG_HeartBeat.CS_HeartBeat_Req) (succ bool, err error) {
	return tcpNet.ResponseHeartBeat(session, uint16(req.SvrPoint))
}

func init() {
	tcpNet.RegisterMessage(uint16(MSG_MainModule.MAINMSG_SERVER), uint16(MSG_Server.SUBMSG_CS_ServerRegister), onSvrRegister)
	tcpNet.RegisterMessage(uint16(MSG_MainModule.MAINMSG_HEARTBEAT), uint16(MSG_HeartBeat.SUBMSG_CS_HeartBeat), onHeartBeat)
}
