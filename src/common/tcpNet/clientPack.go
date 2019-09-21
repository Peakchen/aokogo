package tcpNet

import (
	"encoding/binary"

	"github.com/golang/protobuf/proto"
)

/*
	model: ClientProtocol
	Client to Server, message
*/
type ClientProtocol struct {
	mainid uint16
	subid  uint16
	length uint16
	data   []byte
}

func (self *ClientProtocol) PackAction(Output []byte) {
	var pos int32 = 0
	binary.LittleEndian.PutUint16(Output[pos:], self.mainid)
	pos += 2

	binary.LittleEndian.PutUint16(Output[pos:], self.subid)
	pos += 2

	binary.LittleEndian.PutUint16(Output[pos:], self.length)
	pos += 4

	copy(Output[pos:], self.data)
}

func (self *ClientProtocol) UnPackAction(InData []byte) int32 {
	var pos int32 = 0
	self.mainid = binary.LittleEndian.Uint16(InData[pos:])
	pos += 2

	self.subid = binary.LittleEndian.Uint16(InData[pos:])
	pos += 2

	self.length = binary.LittleEndian.Uint16(InData[pos:])
	pos += 4

	self.data = InData[pos:]
	return pos
}

func (self *ClientProtocol) UnPackData() (msg proto.Message, err error) {
	err = proto.Unmarshal(self.data, msg)
	return
}

func (self *ClientProtocol) PackData(msg proto.Message) (data []byte, err error) {
	data, err = proto.Marshal(msg)
	return
}

func (self *ClientProtocol) GetMessageID() (mainID int32, subID int32) {
	return int32(self.mainid), int32(self.subid)
}
