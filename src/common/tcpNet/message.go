package tcpNet

import (
	"common/Log"
	"common/msgProto/MSG_MainModule"
	"errors"
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
)

type TMessageProc struct {
	proc       reflect.Value
	paramTypes []reflect.Type
}

var (
	_MessageTab map[uint32]*TMessageProc = map[uint32]*TMessageProc{}
)

func RegisterMessage(mainID, subID uint16, proc interface{}) {
	_cmd := EncodeCmd(mainID, subID)
	_, ok := _MessageTab[_cmd]
	if ok {
		return
	}

	cbref := reflect.TypeOf(proc)
	if cbref.Kind() != reflect.Func {
		Log.FmtPrintln("proc type not is func, but is: %v.", cbref.Kind())
		return
	}

	if cbref.NumIn() != 2 {
		Log.FmtPrintln("proc num input is not 2, but is: %v.", cbref.NumIn())
		return
	}

	if cbref.NumOut() != 2 {
		Log.FmtPrintln("proc num output is not 2, but is: %v.", cbref.NumOut())
		return
	}

	if cbref.Out(0) != reflect.TypeOf(bool(false)) {
		Log.FmtPrintln("proc num out 1 is not string, but is: %v.", cbref.Out(0))
		return
	}

	if cbref.Out(1).Name() != "error" {
		Log.FmtPrintln("proc num out 2 is not *proto.Message, but is: %v.", cbref.Out(1), reflect.TypeOf(error(nil)), errors.New("0"), fmt.Errorf("0"))
		return
	}

	paramtypes := []reflect.Type{}
	for i := 0; i < cbref.NumIn(); i++ {
		t := cbref.In(i)
		if t.Kind() == reflect.String ||
			t.Implements(reflect.TypeOf((*proto.Message)(nil)).Elem()) {
			paramtypes = append(paramtypes, t)
		}
	}

	_MessageTab[_cmd] = &TMessageProc{
		proc:       reflect.ValueOf(proc),
		paramTypes: paramtypes,
	}

	return
}

func GetMessageInfo(mainID, subID uint16) (proc *TMessageProc, finded bool) {
	_cmd := EncodeCmd(mainID, subID)
	proc, finded = _MessageTab[_cmd]
	return
}

func GetAllMessageIDs() (msgs []int32) {
	msgs = []int32{}
	for msgid, _ := range _MessageTab {
		msgs = append(msgs, int32(msgid))
	}
	return
}

func MessageCallBack(packobj IMessagePack) (succ bool, err error) {
	mainID, subID := packobj.GetMessageID()
	Log.FmtPrintf("mainid: %v, subID: %v.", mainID, subID)
	msg, cb, err := packobj.UnPackData()
	if err != nil {
		Log.Error("unpack data err: ", err)
		return
	}

	switch mainID {
	case uint16(MSG_MainModule.MAINMSG_SERVER):

	default:

	}

	params := []reflect.Value{
		reflect.ValueOf("1"),
		reflect.ValueOf(msg),
	}

	ret := cb.Call(params)
	succ = ret[0].Interface().(bool)
	err = ret[1].Interface().(error)

	return
}
