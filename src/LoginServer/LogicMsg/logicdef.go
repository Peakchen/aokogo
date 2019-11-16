package LogicMsg

import (
	"common/Log"
	"net"

	"github.com/golang/protobuf/proto"
)

func LoginMessageCallBack(c net.Conn, mainID uint16, subID uint16, msg proto.Message) {
	Log.FmtPrintf("exec login server message call back.", c.RemoteAddr(), c.LocalAddr())
}
