package message

import (
	"common/Log"
	"net"

	"github.com/golang/protobuf/proto"
)

func ExternalGatewayMessageCallBack(c net.Conn, mainID int32, subID int32, msg proto.Message) {
	Log.FmtPrintf("exec external gateway server message call back: %v, %v.", c.RemoteAddr(), c.LocalAddr())
}
