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
	model: ClientProtocol
	Client to Server, message
*/
type ClientProtocol struct {
	routepoint uint16
	mainid     uint16
	subid      uint16
	length     uint32
	data       []byte
	srcdata    []byte
}

func (this *ClientProtocol) PackAction(Output []byte) {
	var pos int32
	binary.LittleEndian.PutUint16(Output[pos:], this.routepoint)
	pos += 2

	binary.LittleEndian.PutUint16(Output[pos:], this.mainid)
	pos += 2

	binary.LittleEndian.PutUint16(Output[pos:], this.subid)
	pos += 2

	binary.LittleEndian.PutUint32(Output[pos:], this.length)
	pos += 4
	Log.FmtPrintln("client PackAction-> data len: ", this.length)
	copy(Output[pos:], this.data)
}

func (this *ClientProtocol) UnPackAction(InData []byte) (pos int32, err error) {
	this.routepoint = binary.LittleEndian.Uint16(InData[pos:])
	pos += 2

	this.mainid = binary.LittleEndian.Uint16(InData[pos:])
	pos += 2

	this.subid = binary.LittleEndian.Uint16(InData[pos:])
	pos += 2

	this.length = binary.LittleEndian.Uint32(InData[pos:])
	pos += 4

	if utls.SliceBytesLength(InData) < int(pos+int32(this.length)) {
		err = fmt.Errorf("client routepoint: %v, mainid: %v, subid: %v; err: InData len: %v, pos: %v, data len: %v.", this.routepoint, this.mainid, this.subid, len(InData), pos, this.length)
		return
	}

	this.data = InData[pos : pos+int32(this.length)]
	this.srcdata = InData
	return pos, nil
}

func (this *ClientProtocol) UnPackData() (msg proto.Message, cb reflect.Value, err error, exist bool) {
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

func (this *ClientProtocol) PackData(msg proto.Message) (data []byte, err error) {
	data, err = proto.Marshal(msg)
	return
}

func (this *ClientProtocol) GetRouteID() (route uint16) {
	return this.routepoint
}

func (this *ClientProtocol) GetMessageID() (mainID uint16, subID uint16) {
	return this.mainid, this.subid
}

func (this *ClientProtocol) SetCmd(routepoint, mainid, subid uint16, data []byte) {
	this.routepoint = routepoint
	this.mainid = mainid
	this.subid = subid
	this.data = data
	this.length = uint32(len(data))
	Log.FmtPrintln("SetCmd data len: ", this.length)
}

func (this *ClientProtocol) Clean() {
	this.length = 0
	this.data = make([]byte, maxMessageSize)
	this.mainid = 0
	this.subid = 0
	this.routepoint = 0
}

func (this *ClientProtocol) PackMsg(routepoint, mainid, subid uint16, msg proto.Message) (out []byte) {
	data, err := proto.Marshal(msg)
	if err != nil {
		Log.FmtPrintln("proto marshal fail, data: ", err)
		return
	}

	this.SetCmd(routepoint, mainid, subid, data)
	out = make([]byte, len(data)+EnMessage_NoDataLen)
	this.PackAction(out)
	return
}

func (this *ClientProtocol) GetSrcMsg() (data []byte) {
	return this.srcdata
}
