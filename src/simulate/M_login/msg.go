package M_login

import (
	"common/Define"
	"common/Log"
	"common/msgProto/MSG_Login"
	"common/msgProto/MSG_MainModule"
	"common/msgProto/MSG_Player"
	"simulate/M_Common"
	"simulate/M_config"
	"strconv"
	"sync"
)

func LoginRun() {
	Log.FmtPrintf("login msg test.")
	MessageRegister()
	var sw sync.WaitGroup

	sw.Add(1)
	//go UserRegister()
	UserLogin()
	//go AlostOfPeopleLogin()
	sw.Wait()
}

func MessageRegister() {

}

func UserRegister() {
	Log.FmtPrintf("user register.")
	loginM := M_Common.NewModule("127.0.0.1:51001", "login")
	for _, item := range M_config.GloginConfig.Get() {
		if item.Register == M_config.CstRegister_No {
			continue
		}
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
	}

}

func UserLogin() {
	Log.FmtPrintf("user login.")
	loginM := M_Common.NewModule("127.0.0.1:51001", "login")
	for _, item := range M_config.GloginConfig.Get() {
		if item.Login == M_config.CstLogin_No {
			continue
		}
		req := &MSG_Login.CS_UserRegister_Req{}
		req.Account = "test1"
		req.Passwd = "abc"
		req.DeviceSerial = "456"
		req.DeviceName = "iso"
		Log.FmtPrintln("UserLogin: ", item.Username, item.Passwd)
		loginM.PushMsg(uint16(Define.ERouteId_ER_Login),
			uint16(MSG_MainModule.MAINMSG_LOGIN),
			uint16(MSG_Login.SUBMSG_CS_Login),
			req)
		go loginM.Run()
		//UserEnter(loginM)
	}
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

func UserEnter(pack *M_Common.TModuleCommon) {
	req := &MSG_Player.CS_EnterServer_Req{}
	pack.PushMsg(uint16(Define.ERouteId_ER_Game),
		uint16(MSG_MainModule.MAINMSG_PLAYER),
		uint16(MSG_Player.SUBMSG_CS_EnterServer),
		req)
	go pack.Run()
}
