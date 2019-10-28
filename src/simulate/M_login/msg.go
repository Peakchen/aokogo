package M_login

import (
	"common/Define"
	"common/Log"
	"common/msgProto/MSG_Login"
	"common/msgProto/MSG_MainModule"
	"simulate/M_Common"
	"strconv"
	"sync"
)

func LoginRun() {
	Log.FmtPrintf("login msg test.")
	MessageRegister()
	var sw sync.WaitGroup

	sw.Add(2)
	go UserRegister()
	go UserLogin()
	go AlostOfPeopleLogin()
	sw.Wait()
}

func MessageRegister() {

}

func UserRegister() {
	Log.FmtPrintf("user register.")
	req := &MSG_Login.CS_UserRegister_Req{}
	req.Account = "test"
	req.Passwd = "abc"
	req.DeviceSerial = "123"
	req.DeviceName = "androd"
	loginM := M_Common.NewModule("127.0.0.1:51001", "login")
	loginM.PushMsg(uint16(Define.ERouteId_ER_Login),
		uint16(MSG_MainModule.MAINMSG_LOGIN),
		uint16(MSG_Login.SUBMSG_CS_UserRegister),
		req)
	loginM.Run()
}

func UserLogin() {
	Log.FmtPrintf("user login.")
	req := &MSG_Login.CS_UserRegister_Req{}
	req.Account = "test"
	req.Passwd = "abc"
	req.DeviceSerial = "456"
	req.DeviceName = "iso"
	loginM := M_Common.NewModule("127.0.0.1:51001", "login")
	loginM.PushMsg(uint16(Define.ERouteId_ER_Login),
		uint16(MSG_MainModule.MAINMSG_LOGIN),
		uint16(MSG_Login.SUBMSG_CS_Login),
		req)
	loginM.Run()
}

func AlostOfPeopleLogin() {
	for i := 1; i <= 100; i++ {
		account := "test" + strconv.Itoa(i)
		Log.FmtPrintf("login account: %v.", account)
		req := &MSG_Login.CS_UserRegister_Req{}
		req.Account = account
		req.Passwd = "abc"
		req.DeviceSerial = "456"
		req.DeviceName = "iso"
		loginM := M_Common.NewModule("127.0.0.1:51001", "login")
		loginM.PushMsg(uint16(Define.ERouteId_ER_Login),
			uint16(MSG_MainModule.MAINMSG_LOGIN),
			uint16(MSG_Login.SUBMSG_CS_Login),
			req)
		loginM.RunEx()
	}
}
