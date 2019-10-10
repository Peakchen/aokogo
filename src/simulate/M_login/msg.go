package M_login

import (
	"common/Log"
	"common/msgProto/MSG_Login"
	"common/msgProto/MSG_MainModule"
	"simulate/M_Common"
)

func LoginRun() {
	Log.FmtPrintf("login msg test.")
	req := &MSG_Login.CS_UserRegister_Req{}
	req.Account = "test"
	req.Passwd = "abc"
	req.DeviceSerial = "123"
	req.DeviceName = "androd"
	loginM := M_Common.NewModule("127.0.0.1:51001", "login")
	loginM.PushMsg(uint16(MSG_MainModule.MAINMSG_LOGIN),
		uint16(MSG_Login.SUBMSG_CS_UserRegister),
		req)
	loginM.Run()
}
