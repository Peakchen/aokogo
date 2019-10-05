package LogicMsg

import (
	"common/Log"
	"net"

	"github.com/golang/protobuf/proto"
)

func GameMessageCallBack(c net.Conn, mainID uint16, subID uint16, msg proto.Message) {
	Log.FmtPrintf("Exec game server message call back.", c.RemoteAddr(), c.LocalAddr())
}

func AfterDialCallBack() {
	Log.FmtPrintf("After dial call back.")
}
