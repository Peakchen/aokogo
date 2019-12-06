package msgImp

import "github.com/golang/protobuf/proto"

func RegisterMsg(mainID, subID uint16, pb string, obj proto.Message) {
	msgobj.register(mainID, subID, pb, obj)
}

func GetMsgPb(cmd uint32) (pb *PbMsg, exist bool) {
	pb, exist = msgobj.Data[cmd]
	return
}
