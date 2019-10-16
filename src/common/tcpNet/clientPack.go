package tcpNet

import (
	"common/Log"
	"encoding/binary"
	"reflect"

	"github.com/golang/protobuf/proto"
)

/*
	model: ClientProtocol
	Client to Server, message
*/
type ClientProtocol struct {
	routepoint uint16
	mainid     uint16
	subid      uint16
	length     uint32
	data       []byte
}

func (self *ClientProtocol) PackAction(Output []byte) {
	var pos int32
	binary.LittleEndian.PutUint16(Output[pos:], self.routepoint)
	pos += 2

	binary.LittleEndian.PutUint16(Output[pos:], self.mainid)
	pos += 2

	binary.LittleEndian.PutUint16(Output[pos:], self.subid)
	pos += 2

	binary.LittleEndian.PutUint32(Output[pos:], self.length)
	pos += 4
	Log.FmtPrintln("client PackAction-> data len: ", self.length)
	copy(Output[pos:], self.data)
}

func (self *ClientProtocol) UnPackAction(InData []byte) (pos int32, err error) {
	self.routepoint = binary.LittleEndian.Uint16(InData[pos:])
	pos += 2

	self.mainid = binary.LittleEndian.Uint16(InData[pos:])
	pos += 2

	self.subid = binary.LittleEndian.Uint16(InData[pos:])
	pos += 2

	self.length = binary.LittleEndian.Uint32(InData[pos:])
	pos += 4

	self.data = InData[pos:]
	return pos, nil
}

func (self *ClientProtocol) UnPackData() (msg proto.Message, cb reflect.Value, err error) {
	err = proto.Unmarshal(self.data, msg)
	return
}

func (self *ClientProtocol) PackData(msg proto.Message) (data []byte, err error) {
	data, err = proto.Marshal(msg)
	return
}

func (self *ClientProtocol) GetRouteID() (route uint16) {
	return self.routepoint
}

func (self *ClientProtocol) GetMessageID() (mainID uint16, subID uint16) {
	return self.mainid, self.subid
}

func (self *ClientProtocol) SetCmd(routepoint, mainid, subid uint16, data []byte) {
	self.routepoint = routepoint
	self.mainid = mainid
	self.subid = subid
	self.data = data
	self.length = uint32(len(data))
	Log.FmtPrintln("SetCmd data len: ", self.length)
}

func (self *ClientProtocol) Clean() {
	self.length = 0
	self.data = make([]byte, maxMessageSize)
	self.mainid = 0
	self.subid = 0
	self.routepoint = 0
}

func (self *ClientProtocol) PackMsg(routepoint, mainid, subid uint16, msg proto.Message) (out []byte) {
	data, err := proto.Marshal(msg)
	if err != nil {
		Log.FmtPrintln("proto marshal fail, data: ", err)
		return
	}

	self.SetCmd(routepoint, mainid, subid, data)
	out = make([]byte, len(data)+EnMessage_NoDataLen)
	self.PackAction(out)
	return
}
