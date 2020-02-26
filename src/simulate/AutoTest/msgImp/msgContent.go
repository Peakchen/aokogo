package msgImp

import (
	"common/akNet"
	"reflect"

	"github.com/golang/protobuf/proto"
)

type PbMsg struct {
	Pb  string
	Msg reflect.Value
	MT  reflect.Type
}

type TMsgContent struct {
	Data map[uint32]*PbMsg
}

var (
	msgobj *TMsgContent
)

func init() {
	msgobj = &TMsgContent{}
	msgobj.Data = map[uint32]*PbMsg{}
}

func (this *TMsgContent) register(mainID, subID uint16, pb string, pbobj proto.Message) {
	cmd := akNet.EncodeCmd(mainID, subID)
	this.Data[cmd] = &PbMsg{
		Pb:  pb,
		Msg: reflect.ValueOf(pbobj),
		MT:  reflect.TypeOf(pbobj),
	}
}
