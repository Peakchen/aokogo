package LogicMsg

import (
	"LoginServer/Logic/UserAccount"
	"common/Log"
	"common/msgProto/MSG_Login"
	"common/msgProto/MSG_MainModule"
	"common/tcpNet"
)

func onUserBind(key string, req *MSG_Login.CS_UserBind_Req) (succ bool, err error) {
	Log.FmtPrintf("onUserBind recv: %v, %v.", key, req.Account, req.Passwd)

	return
}

func onUserRegister(session *tcpNet.TcpSession, req *MSG_Login.CS_UserRegister_Req) (succ bool, err error) {
	Log.FmtPrintf("[onUserRegister] SessionID: %v, Account: %v, Passwd: %v, DeviceSerial: %v, DeviceName: %v.", session.SessionID, req.Account, req.Passwd, req.DeviceSerial, req.DeviceName)
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
	} else {
		session.SetIdentify(acc.Identify())
	}

	return session.SendMsg(uint16(session.SvrType),
		uint16(MSG_MainModule.MAINMSG_LOGIN),
		uint16(MSG_Login.SUBMSG_SC_UserRegister),
		rsp)
}

func onUserLogin(session *tcpNet.TcpSession, req *MSG_Login.CS_Login_Req) (succ bool, err error) {
	Log.FmtPrintf("[onUserLogin] SessionID: %v, Account: %v, Passwd: %v, DeviceSerial: %v, DeviceName: %v.", session.SessionID, req.Account, req.Passwd, req.DeviceSerial, req.DeviceName)

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

	return session.SendInnerMsg(uint16(MSG_MainModule.MAINMSG_LOGIN),
		uint16(MSG_Login.SUBMSG_SC_Login),
		rsp)
}

func init() {
	tcpNet.RegisterMessage(uint16(MSG_MainModule.MAINMSG_LOGIN), uint16(MSG_Login.SUBMSG_CS_UserBind), onUserBind)
	tcpNet.RegisterMessage(uint16(MSG_MainModule.MAINMSG_LOGIN), uint16(MSG_Login.SUBMSG_CS_UserRegister), onUserRegister)
	tcpNet.RegisterMessage(uint16(MSG_MainModule.MAINMSG_LOGIN), uint16(MSG_Login.SUBMSG_CS_Login), onUserLogin)
}
