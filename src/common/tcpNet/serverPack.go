package tcpNet

import (
	"encoding/binary"
	"reflect"

	"github.com/golang/protobuf/proto"
)

/*
	model: ServerProtocol
	server to server, message
*/
type ServerProtocol struct {
	mainid uint16
	subid  uint16
	length uint16
	data   []byte
}

func (self *ServerProtocol) PackAction(Output []byte) {
	var pos int32 = 0
	binary.LittleEndian.PutUint16(Output[pos:], self.mainid)
	pos += 2

	binary.LittleEndian.PutUint16(Output[pos:], self.subid)
	pos += 2

	binary.LittleEndian.PutUint16(Output[pos:], self.length)
	pos += 4

	copy(Output[pos:], self.data)
}

func (self *ServerProtocol) UnPackAction(InData []byte) int32 {
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

func (self *ServerProtocol) UnPackData() (msg proto.Message, err error) {
	msgsrc := reflect.New(reflect.TypeOf((*proto.Message)(nil)).Elem()).Type()
	dst := reflect.New(msgsrc.Elem()).Interface()
	err = proto.Unmarshal(self.data, dst.(proto.Message))
	msg = dst.(proto.Message)
	return
}

func (self *ServerProtocol) PackData(msg proto.Message) (data []byte, err error) {
	data, err = proto.Marshal(msg)
	return
}

func (self *ServerProtocol) GetMessageID() (mainID int32, subID int32) {
	return int32(self.mainid), int32(self.subid)
}
