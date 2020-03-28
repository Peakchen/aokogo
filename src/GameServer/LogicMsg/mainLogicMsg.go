package LogicMsg

import (
	"GameServer/logic"
	"common/Log"
	"common/akNet"
	"common/msgProto/MSG_MainModule"
	"common/msgProto/MSG_Player"
)

func Init() {

}

func onEnterServer(session akNet.TcpSession, req *MSG_Player.CS_EnterServer_Req) (succ bool, err error) {
	Log.FmtPrintf("enter Server player(%v) enter game server.", session.GetIdentify())
	logic.EnterGameReady(session)
	rsp := &MSG_Player.SC_EnterServer_Rsp{}
	rsp.Ret = MSG_Player.ErrorCode_Success
	return session.SendInnerClientMsg(uint16(MSG_MainModule.MAINMSG_PLAYER),
		uint16(MSG_Player.SUBMSG_SC_EnterServer),
		rsp)
}

func onLeaveServer(session akNet.TcpSession, req *MSG_Player.CS_LeaveServer_Req) (succ bool, err error) {
	Log.FmtPrintf("leave Server player(%v).", session.GetIdentify())
	return true, nil
}

func init() {
	akNet.RegisterMessage(uint16(MSG_MainModule.MAINMSG_PLAYER), uint16(MSG_Player.SUBMSG_CS_EnterServer), onEnterServer)
	akNet.RegisterMessage(uint16(MSG_MainModule.MAINMSG_PLAYER), uint16(MSG_Player.SUBMSG_CS_LeaveServer), onLeaveServer)
}
