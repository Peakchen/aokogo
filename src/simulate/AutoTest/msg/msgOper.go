package msg

import "github.com/golang/protobuf/proto"

func RegisterMsg(mainID, subID uint16, pb proto.Message) {
	msgobj.register(mainID, subID, pb)
}
