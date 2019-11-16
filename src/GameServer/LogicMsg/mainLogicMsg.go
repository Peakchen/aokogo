package LogicMsg

import (
	"GameServer/logic"
	"common/Log"
	"common/msgProto/MSG_MainModule"
	"common/msgProto/MSG_Player"
	"common/tcpNet"
	"net"

	"github.com/golang/protobuf/proto"
)

func GameMessageCallBack(c net.Conn, mainID uint16, subID uint16, msg proto.Message) {
	Log.FmtPrintf("Exec game server message call back.", c.RemoteAddr(), c.LocalAddr())
}

func AfterDialCallBack(s *tcpNet.TcpSession) {
	Log.FmtPrintf("After dial call back.")
}

func onServer(session *tcpNet.TcpSession, req *MSG_Player.CS_EnterServer_Req) (succ bool, err error) {
	Log.FmtPrintf("SessionID: %v, player(%v) enter game server.", session.SessionID, session.StrIdentify)
	logic.EnterGameReady(session.StrIdentify)
	rsp := &MSG_Player.SC_EnterServer_Rsp{}
	rsp.Ret = MSG_Player.ErrorCode_Success
	return session.SendMsg(uint16(session.SvrType),
		uint16(MSG_MainModule.MAINMSG_PLAYER),
		uint16(MSG_Player.SUBMSG_SC_EnterServer),
		rsp)
}

func init() {
	tcpNet.RegisterMessage(uint16(MSG_MainModule.MAINMSG_PLAYER), uint16(MSG_Player.SUBMSG_CS_EnterServer), onServer)
}
