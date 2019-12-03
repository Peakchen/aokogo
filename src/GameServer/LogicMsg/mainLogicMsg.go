package LogicMsg

import (
	"GameServer/logic"
	"common/Log"
	"common/msgProto/MSG_MainModule"
	"common/msgProto/MSG_Player"
	"common/tcpNet"
)

func Init() {

}

func onServer(session tcpNet.TcpSession, req *MSG_Player.CS_EnterServer_Req) (succ bool, err error) {
	Log.FmtPrintf("onServer player(%v) enter game server.", session.GetIdentify())
	logic.EnterGameReady(session)
	rsp := &MSG_Player.SC_EnterServer_Rsp{}
	rsp.Ret = MSG_Player.ErrorCode_Success
	return session.SendInnerMsg(session.GetIdentify(),
		uint16(MSG_MainModule.MAINMSG_PLAYER),
		uint16(MSG_Player.SUBMSG_SC_EnterServer),
		rsp)
}

func init() {
	tcpNet.RegisterMessage(uint16(MSG_MainModule.MAINMSG_PLAYER), uint16(MSG_Player.SUBMSG_CS_EnterServer), onServer)
}
