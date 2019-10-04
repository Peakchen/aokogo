package message

import (
	"common/Log"
	"net"

	"github.com/golang/protobuf/proto"
)

func GameMessageCallBack(c net.Conn, mainID int32, subID int32, msg proto.Message) {
	Log.FmtPrintf("exec game server message call back.", c.RemoteAddr(), c.LocalAddr())
}
