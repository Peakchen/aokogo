package LogicMsg

import (
	"LoginServer/Logic/UserAccount"
	"common/Log"
	"common/akNet"
	"common/msgProto/MSG_Login"
	"common/msgProto/MSG_MainModule"
)

func onUserBind(key string, req *MSG_Login.CS_UserBind_Req) (succ bool, err error) {
	Log.FmtPrintf("onUserBind recv: %v, %v.", key, req.Account, req.Passwd)

	return
}

func onUserRegister(session akNet.TcpSession, req *MSG_Login.CS_UserRegister_Req) (succ bool, err error) {
	Log.FmtPrintf("[onUserRegister] identify: %v, Account: %v, Passwd: %v, DeviceSerial: %v, DeviceName: %v.", session.GetIdentify(), req.Account, req.Passwd, req.DeviceSerial, req.DeviceName)
	rsp := &MSG_Login.SC_UserRegister_Rsp{}
	rsp.Ret = MSG_Login.ErrorCode_Success

	acc := &UserAccount.TUserAcc{
		UserName:   req.Account,
		Passwd:     req.Passwd,
		DeviceNo:   req.DeviceSerial,
		DeviceType: req.DeviceName,
	}

	if err, exist := UserAccount.RegisterUseAcc(acc); err != nil || !exist {
		rsp.Ret = MSG_Login.ErrorCode_Fail
	}

	session.SetIdentify(acc.Identify())

	return session.SendInnerMsg(acc.Identify(), //uint16(Define.ERouteId_ER_ISG),
		uint16(MSG_MainModule.MAINMSG_LOGIN),
		uint16(MSG_Login.SUBMSG_SC_UserRegister),
		rsp)
}

func onUserLogin(session akNet.TcpSession, req *MSG_Login.CS_Login_Req) (succ bool, err error) {
	Log.FmtPrintf("[onUserLogin] identify: %v, Account: %v, Passwd: %v, DeviceSerial: %v, DeviceName: %v.", session.GetIdentify(), req.Account, req.Passwd, req.DeviceSerial, req.DeviceName)

	rsp := &MSG_Login.SC_Login_Rsp{}
	rsp.Ret = MSG_Login.ErrorCode_Success

	acc := &UserAccount.TUserAcc{
		UserName:   req.Account,
		Passwd:     req.Passwd,
		DeviceNo:   req.DeviceSerial,
		DeviceType: req.DeviceName,
	}

	if _, exist := UserAccount.GetUserAcc(acc); !exist {
		rsp.Ret = MSG_Login.ErrorCode_UserNotExistOrPasswdErr
	}
	session.SetIdentify(acc.Identify())
	return session.SendInnerMsg(acc.Identify(),
		uint16(MSG_MainModule.MAINMSG_LOGIN),
		uint16(MSG_Login.SUBMSG_SC_Login),
		rsp)
}

func init() {
	akNet.RegisterMessage(uint16(MSG_MainModule.MAINMSG_LOGIN), uint16(MSG_Login.SUBMSG_CS_UserBind), onUserBind)
	akNet.RegisterMessage(uint16(MSG_MainModule.MAINMSG_LOGIN), uint16(MSG_Login.SUBMSG_CS_UserRegister), onUserRegister)
	akNet.RegisterMessage(uint16(MSG_MainModule.MAINMSG_LOGIN), uint16(MSG_Login.SUBMSG_CS_Login), onUserLogin)
}
