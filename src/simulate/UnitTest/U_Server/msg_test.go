package U_Server

import (
	"common/Define"
	"common/Log"
	"common/msgProto/MSG_MainModule"
	"common/msgProto/MSG_Server"
	"simulate/TestCommon"
	"testing"
)

var serverhost string = "127.0.0.1:51001"

func TestServer(t *testing.T) {
	Log.FmtPrintf("server msg test.")
	req := &MSG_Server.CS_EnterServer_Req{}
	req.Enter = 2
	serverM := TestCommon.NewModule(serverhost, "server")
	serverM.PushMsg(uint16(Define.ERouteId_ER_Login),
		uint16(MSG_MainModule.MAINMSG_SERVER),
		uint16(MSG_Server.SUBMSG_CS_ServerRegister),
		req)
	serverM.Run()
}
