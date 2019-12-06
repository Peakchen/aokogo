package U_login

import (
	"common/Define"
	"common/Log"
	"common/msgProto/MSG_Login"
	"common/msgProto/MSG_MainModule"
	"simulate/TestCommon"
	"simulate/UnitTest/U_config"
	"time"
)

const (
	cstSendInterval = 200
)

func UserRegister(host string, module string) {
	Log.FmtPrintf("user register.")
	loginM := TestCommon.NewModule(host, module)
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
}
