package tcpNet

import (
	"common/Log"
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
)

/*
	model: ServerProtocol
	server to server, message
*/
type ServerProtocol struct {
	routepoint uint16
	mainid     uint16
	subid      uint16
	length     uint32
	data       []byte
}

func (self *ServerProtocol) PackAction(Output []byte) {
	var pos int32
	binary.LittleEndian.PutUint16(Output[pos:], self.routepoint)
	pos += 2

	binary.LittleEndian.PutUint16(Output[pos:], self.mainid)
	pos += 2

	binary.LittleEndian.PutUint16(Output[pos:], self.subid)
	pos += 2

	binary.LittleEndian.PutUint32(Output[pos:], self.length)
	pos += 4

	Log.FmtPrintln("server PackAction-> data len: ", self.length)
	copy(Output[pos:], self.data)
}

func (self *ServerProtocol) UnPackAction(InData []byte) (pos int32, err error) {
	self.routepoint = binary.LittleEndian.Uint16(InData[pos:])
	pos += 2

	self.mainid = binary.LittleEndian.Uint16(InData[pos:])
	pos += 2

	self.subid = binary.LittleEndian.Uint16(InData[pos:])
	pos += 2

	Log.FmtPrintln("server UnPackAction-> len: ", self.length)
	self.length = binary.LittleEndian.Uint32(InData[pos:])
	pos += 4

	if len(InData) < int(pos+int32(self.length)) {
		err = fmt.Errorf("err: InData len: %v, pos: %v, data len: %v.", len(InData), pos, self.length)
		return
	}

	Log.FmtPrintf("normal: InData len: %v, pos: %v, data len: %v.", len(InData), pos, self.length)
	self.data = InData[pos : pos+int32(self.length)]
	return pos, nil
}

func (self *ServerProtocol) UnPackData() (msg proto.Message, cb reflect.Value, err error) {
	err = nil
	mt, finded := GetMessageInfo(self.mainid, self.subid)
	if !finded {
		err = fmt.Errorf("can not regist message, mainid: %v, subid: %v.", self.mainid, self.subid)
		return
	}

	dst := reflect.New(mt.paramTypes[1].Elem()).Interface()
	err = proto.Unmarshal(self.data, dst.(proto.Message))
	if err != nil {
		return
	}
	msg = dst.(proto.Message)
	cb = mt.proc
	return
}

func (self *ServerProtocol) PackData(msg proto.Message) (data []byte, err error) {
	data, err = proto.Marshal(msg)
	return
}

func (self *ServerProtocol) GetRouteID() (route uint16) {
	return self.routepoint
}

func (self *ServerProtocol) GetMessageID() (mainID uint16, subID uint16) {
	return self.mainid, self.subid
}

func (self *ServerProtocol) Clean() {
	self.length = 0
	self.data = make([]byte, maxMessageSize)
	self.mainid = 0
	self.subid = 0
}

func (self *ServerProtocol) SetCmd(routepoint, mainid, subid uint16, data []byte) {
	self.mainid = mainid
	self.subid = subid
	self.data = data
	self.length = uint32(len(data))
	Log.FmtPrintln("SetCmd data len: ", self.length)
}

func (self *ServerProtocol) PackMsg(routepoint, mainid, subid uint16, msg proto.Message) (out []byte) {
	data, err := proto.Marshal(msg)
	if err != nil {
		Log.FmtPrintln("server proto marshal fail, data: ", err)
		return
	}

	self.SetCmd(routepoint, mainid, subid, data)
	out = make([]byte, len(data)+EnMessage_NoDataLen)
	self.PackAction(out)
	return
}
