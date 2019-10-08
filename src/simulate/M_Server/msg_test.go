package M_Server

import (
	"common/Log"
	"common/msgProto/MSG_MainModule"
	"common/msgProto/MSG_Server"
	"simulate/M_Common"
	"testing"
)

var serverhost string = "127.0.0.1:51001"

func TestServer(t *testing.T) {
	Log.FmtPrintf("server msg test.")
	req := &MSG_Server.CS_EnterServer_Req{}
	req.Enter = 2
	serverM := M_Common.NewModule(serverhost, "server")
	serverM.PushMsg(uint16(MSG_MainModule.MAINMSG_SERVER),
		uint16(MSG_Server.SUBMSG_CS_ServerRegister),
		req)
	serverM.Run()
}
