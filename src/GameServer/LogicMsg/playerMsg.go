package LogicMsg

import (
	"GameServer/logic/Player"
	"common/Log"
	"common/msgProto/MSG_MainModule"
	"common/msgProto/MSG_Player"
	"common/tcpNet"
)

func onGetPlayerInfo(session tcpNet.TcpSession, req *MSG_Player.CS_PlayerInfo_Req) (succ bool, err error) {
	Log.FmtPrintf("[onGetPlayerInfo] SessionID: %v.", session.GetIdentify())

	rsp := &MSG_Player.SC_PlayerInfo_Rsp{}
	rsp.Ret = MSG_Player.ErrorCode_Success
	data := Player.GetPlayer(session.GetIdentify())
	if data == nil {
		return
	}
	Log.FmtPrintf("get player info: %v.", data.BaseInfo[MSG_Player.EmBaseInfo_Name])
	return session.SendInnerMsg(session.GetIdentify(),
		uint16(MSG_MainModule.MAINMSG_PLAYER),
		uint16(MSG_Player.SUBMSG_SC_PlayerInfo),
		rsp)
}

func init() {
	tcpNet.RegisterMessage(uint16(MSG_MainModule.MAINMSG_PLAYER), uint16(MSG_Player.SUBMSG_CS_PlayerInfo), onGetPlayerInfo)
}
