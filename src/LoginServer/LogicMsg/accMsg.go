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
	Log.FmtPrintf("[onUserRegister recv] SessionID: %v, Account: %v, Passwd: %v, DeviceSerial: %v, DeviceName: %v.", session.SessionID, req.Account, req.Passwd, req.DeviceSerial, req.DeviceName)
	rsp := &MSG_Login.SC_UserRegister_Rsp{}
	rsp.Ret = MSG_Login.ErrorCode_Success

	acc := &UserAccount.TUserAcc{
		UserName:   req.Account,
		Passwd:     req.Passwd,
		DeviceNo:   req.DeviceSerial,
		DeviceType: req.DeviceName,
	}

	if err, exist := UserAccount.RegisterUseAcc(acc); err != nil && !exist {
		rsp.Ret = MSG_Login.ErrorCode_Fail
	}

	return session.SendMsg(uint16(session.SrcPoint),
		uint16(MSG_MainModule.MAINMSG_LOGIN),
		uint16(MSG_Login.SUBMSG_SC_UserRegister),
		rsp)
}

func init() {
	tcpNet.RegisterMessage(uint16(MSG_MainModule.MAINMSG_LOGIN), uint16(MSG_Login.SUBMSG_CS_UserBind), onUserBind)
	tcpNet.RegisterMessage(uint16(MSG_MainModule.MAINMSG_LOGIN), uint16(MSG_Login.SUBMSG_CS_UserRegister), onUserRegister)
}
