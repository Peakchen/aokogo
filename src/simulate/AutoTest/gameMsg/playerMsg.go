package gameMsg

import (
	"common/Log"
	"common/msgProto/MSG_MainModule"
	"common/msgProto/MSG_Player"
	"simulate/AutoTest/msg"
)

func init() {
	Log.FmtPrintln("run game player test.")

	msg.RegisterMsg(uint16(MSG_MainModule.MAINMSG_PLAYER), uint16(MSG_Player.SUBMSG_CS_EnterServer), &MSG_Player.CS_EnterServer_Req{})
	msg.RegisterMsg(uint16(MSG_MainModule.MAINMSG_PLAYER), uint16(MSG_Player.SUBMSG_SC_EnterServer), &MSG_Player.SC_EnterServer_Rsp{})
}
