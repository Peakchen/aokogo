package msg

import (
	"common/tcpNet"

	"github.com/golang/protobuf/proto"
)

type TMsgContent struct {
	msg map[uint32]proto.Message
}

var (
	msgobj *TMsgContent
)

func init() {
	msgobj = &TMsgContent{}
	msgobj.msg = map[uint32]proto.Message{}
}

func (this *TMsgContent) register(mainID, subID uint16, pb proto.Message) {
	cmd := tcpNet.EncodeCmd(mainID, subID)
	this.msg[cmd] = pb
}
