package LogicMsg

import (
	"common/Log"
	"common/msgProto/MSG_Login"
	"common/msgProto/MSG_MainModule"
	"common/tcpNet"
	"net"

	"github.com/golang/protobuf/proto"
)

func LoginMessageCallBack(c net.Conn, mainID uint16, subID uint16, msg proto.Message) {
	Log.FmtPrintf("exec login server message call back.", c.RemoteAddr(), c.LocalAddr())
}

func onUserBind(key string, req *MSG_Login.CS_UserBind_Req) (succ bool, err error) {
	Log.FmtPrintf("onUserBind recv: %v, %v.", key, req.Account, req.Passwd)

	return
}

func onUserRegister(session *tcpNet.TcpSession, req *MSG_Login.CS_UserRegister_Req) (succ bool, err error) {
	Log.FmtPrintf("onUserRegister recv: %v, %v.", session.SessionID, req.Account, req.Passwd, req.DeviceSerial, req.DeviceName)

	return
}

func init() {
	tcpNet.RegisterMessage(uint16(MSG_MainModule.MAINMSG_LOGIN), uint16(MSG_Login.SUBMSG_CS_UserBind), onUserBind)
	tcpNet.RegisterMessage(uint16(MSG_MainModule.MAINMSG_LOGIN), uint16(MSG_Login.SUBMSG_CS_UserRegister), onUserRegister)
}
