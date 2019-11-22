package tcpNet

import (
	"common/Log"
	"common/utls"
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
	data       []byte //消息体
	srcdata    []byte //源消息（未解包）
}

func (this *ServerProtocol) PackAction(Output []byte) {
	var pos int32
	binary.LittleEndian.PutUint16(Output[pos:], this.routepoint)
	pos += 2

	binary.LittleEndian.PutUint16(Output[pos:], this.mainid)
	pos += 2

	binary.LittleEndian.PutUint16(Output[pos:], this.subid)
	pos += 2

	binary.LittleEndian.PutUint32(Output[pos:], this.length)
	pos += 4

	Log.FmtPrintln("server PackAction-> data len: ", this.length)
	copy(Output[pos:], this.data)
}

func (this *ServerProtocol) UnPackAction(InData []byte) (pos int32, err error) {
	this.routepoint = binary.LittleEndian.Uint16(InData[pos:])
	pos += 2

	this.mainid = binary.LittleEndian.Uint16(InData[pos:])
	pos += 2

	this.subid = binary.LittleEndian.Uint16(InData[pos:])
	pos += 2

	Log.FmtPrintln("server UnPackAction-> len: ", this.length)
	this.length = binary.LittleEndian.Uint32(InData[pos:])
	pos += 4

	if utls.SliceBytesLength(InData) < int(pos+int32(this.length)) {
		err = fmt.Errorf("server routepoint: %v, mainid: %v, subid: %v; err: InData len: %v, pos: %v, data len: %v.", this.routepoint, this.mainid, this.subid, len(InData), pos, this.length)
		return
	}

	Log.FmtPrintf("normal: InData len: %v, pos: %v, data len: %v.", len(InData), pos, this.length)
	Log.FmtPrintf("message head: routepoint: %v, mainid: %v, subid: %v.", this.routepoint, this.mainid, this.subid)
	this.data = InData[pos : pos+int32(this.length)]
	this.srcdata = InData
	return pos, nil
}

func (this *ServerProtocol) UnPackData() (msg proto.Message, cb reflect.Value, err error, exist bool) {
	err = nil
	mt, finded := GetMessageInfo(this.mainid, this.subid)
	if !finded {
		err = fmt.Errorf("can not regist message, mainid: %v, subid: %v.", this.mainid, this.subid)
		return
	}

	exist = true
	dst := reflect.New(mt.paramTypes[1].Elem()).Interface()
	err = proto.Unmarshal(this.data, dst.(proto.Message))
	if err != nil {
		return
	}
	msg = dst.(proto.Message)
	cb = mt.proc
	return
}

func (this *ServerProtocol) PackData(msg proto.Message) (data []byte, err error) {
	data, err = proto.Marshal(msg)
	return
}

func (this *ServerProtocol) GetRouteID() (route uint16) {
	return this.routepoint
}

func (this *ServerProtocol) GetMessageID() (mainID uint16, subID uint16) {
	return this.mainid, this.subid
}

func (this *ServerProtocol) Clean() {
	this.length = 0
	this.data = make([]byte, maxMessageSize)
	this.mainid = 0
	this.subid = 0
}

func (this *ServerProtocol) SetCmd(routepoint, mainid, subid uint16, data []byte) {
	this.mainid = mainid
	this.subid = subid
	this.data = data
	this.length = uint32(len(data))
	Log.FmtPrintln("SetCmd data len: ", this.length)
}

func (this *ServerProtocol) PackMsg(routepoint, mainid, subid uint16, msg proto.Message) (out []byte) {
	data, err := proto.Marshal(msg)
	if err != nil {
		Log.FmtPrintln("server proto marshal fail, data: ", err)
		return
	}

	this.SetCmd(routepoint, mainid, subid, data)
	out = make([]byte, len(data)+EnMessage_NoDataLen)
	this.PackAction(out)
	return
}

func (this *ServerProtocol) GetSrcMsg() (data []byte) {
	return this.srcdata
}
