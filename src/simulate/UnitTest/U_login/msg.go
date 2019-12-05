package U_login

import (
	"common/Define"
	"common/Log"
	"common/msgProto/MSG_Login"
	"common/msgProto/MSG_MainModule"
	"common/msgProto/MSG_Player"
	"simulate/TestCommon"
	"simulate/UnitTest/U_config"
	"strconv"
	"sync"
	"time"
)

const (
	cstSendInterval = 200
)

func LoginRun() {
	Log.FmtPrintf("login msg test.")
	MessageRegister()
	var sw sync.WaitGroup

	sw.Add(1)
	go UserRegister()
	//go UserLogin()
	//go AlostOfPeopleLogin()
	sw.Wait()
}

func MessageRegister() {

}

func UserRegister() {
	Log.FmtPrintf("user register.")
	loginM := TestCommon.NewModule("127.0.0.1:51001", "login")
	for _, item := range U_config.GloginConfig.Get() {
		if item.Register != U_config.CstRegister_No {
			req := &MSG_Login.CS_UserRegister_Req{}
			req.Account = item.Username
			req.Passwd = item.Passwd
			req.DeviceSerial = "123"
			req.DeviceName = "androd"
			Log.FmtPrintln("UserRegister: ", item.Username, item.Passwd)
			loginM.PushMsg(uint16(Define.ERouteId_ER_Login),
				uint16(MSG_MainModule.MAINMSG_LOGIN),
				uint16(MSG_Login.SUBMSG_CS_UserRegister),
				req)
			go loginM.Run()
			time.Sleep(time.Duration(cstSendInterval) * time.Millisecond)
		}

		UserLogin(loginM, item)
	}

}

func UserLogin(pack *TestCommon.TModuleCommon, item *U_config.TSimulateLoginBase) {
	Log.FmtPrintf("user login.")
	if item.Login == U_config.CstLogin_No {
		return
	}

	req := &MSG_Login.CS_UserRegister_Req{}
	req.Account = item.Username
	req.Passwd = item.Passwd
	req.DeviceSerial = "456"
	req.DeviceName = "iso"
	Log.FmtPrintln("UserLogin: ", item.Username, item.Passwd)
	pack.PushMsg(uint16(Define.ERouteId_ER_Login),
		uint16(MSG_MainModule.MAINMSG_LOGIN),
		uint16(MSG_Login.SUBMSG_CS_Login),
		req)
	go pack.Run()
	time.Sleep(time.Duration(cstSendInterval) * time.Millisecond)
	UserEnter(pack)

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
		loginM := TestCommon.NewModule("127.0.0.1:51001", "login")
		loginM.PushMsg(uint16(Define.ERouteId_ER_Login),
			uint16(MSG_MainModule.MAINMSG_LOGIN),
			uint16(MSG_Login.SUBMSG_CS_Login),
			req)
		loginM.RunEx()
	}
}

func UserEnter(pack *TestCommon.TModuleCommon) {
	req := &MSG_Player.CS_EnterServer_Req{}
	pack.PushMsg(uint16(Define.ERouteId_ER_Game),
		uint16(MSG_MainModule.MAINMSG_PLAYER),
		uint16(MSG_Player.SUBMSG_CS_EnterServer),
		req)
	go pack.Run()
	time.Sleep(time.Duration(cstSendInterval) * time.Millisecond)
}
