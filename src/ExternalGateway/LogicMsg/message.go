package LogicMsg

import (
	"common/Log"
	"common/msgProto/MSG_MainModule"
	"common/msgProto/MSG_Server"
	"common/tcpNet"
	"fmt"
	"net"

	"github.com/golang/protobuf/proto"
)

func ExternalGatewayMessageCallBack(c net.Conn, mainID uint16, subID uint16, msg proto.Message) {
	Log.FmtPrintf("exec external gateway server message call back: %v, %v.", c.RemoteAddr(), c.LocalAddr())
}

func onServer(key string, req *MSG_Server.CS_EnterServer_Req) (succ bool, err error) {
	Log.FmtPrintf("onServer recv: %v, %v.", key, req.Enter)
	return
}

func onSvrRegister(key string, req *MSG_Server.CS_ServerRegister_Req) (succ bool, err error) {
	Log.FmtPrintf("onSvrRegister recv: %v, %v.", key, req.ServerType)
	var (
		msgfmt string
	)
	for _, id := range req.Msgs {
		mainid, subid := tcpNet.DecodeCmd(uint32(id))
		msgfmt += fmt.Sprintf("[mainid: %v, subid: %v]\t", mainid, subid)
	}
	msgfmt += "\n"
	Log.FmtPrintln("message context: ", msgfmt)
	return
}

func init() {
	tcpNet.RegisterMessage(uint16(MSG_MainModule.MAINMSG_SERVER), uint16(MSG_Server.SUBMSG_CS_EnterServer), onServer)
	tcpNet.RegisterMessage(uint16(MSG_MainModule.MAINMSG_SERVER), uint16(MSG_Server.SUBMSG_CS_ServerRegister), onSvrRegister)
}
