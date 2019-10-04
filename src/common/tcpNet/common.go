package tcpNet

import (
	"net"
	"reflect"

	"github.com/golang/protobuf/proto"
)

/*
Copyright (c) <year> <copyright holders>

"Anti 996" License Version 1.0 (Draft)

Permission is hereby granted to any individual or legal entity
obtaining a copy of this licensed work (including the source code,
documentation and/or related items, hereinafter collectively referred
to as the "licensed work"), free of charge, to deal with the licensed
work for any purpose, including without limitation, the rights to use,
reproduce, modify, prepare derivative works of, distribute, publish
and sublicense the licensed work, subject to the following conditions:

1. The individual or the legal entity must conspicuously display,
without modification, this License and the notice on each redistributed
or derivative copy of the Licensed Work.

2. The individual or the legal entity must strictly comply with all
applicable laws, regulations, rules and standards of the jurisdiction
relating to labor and employment where the individual is physically
located or where the individual was born or naturalized; or where the
legal entity is registered or is operating (whichever is stricter). In
case that the jurisdiction has no such laws, regulations, rules and
standards or its laws, regulations, rules and standards are
unenforceable, the individual or the legal entity are required to
comply with Core International Labor Standards.

3. The individual or the legal entity shall not induce, metaphor or force
its employee(s), whether full-time or part-time, or its independent
contractor(s), in any methods, to agree in oral or written form, to
directly or indirectly restrict, weaken or relinquish his or her
rights or remedies under such laws, regulations, rules and standards
relating to labor and employment as mentioned above, no matter whether
such written or oral agreement are enforceable under the laws of the
said jurisdiction, nor shall such individual or the legal entity
limit, in any methods, the rights of its employee(s) or independent
contractor(s) from reporting or complaining to the copyright holder or
relevant authorities monitoring the compliance of the license about
its violation(s) of the said license.

THE LICENSED WORK IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE COPYRIGHT HOLDER BE LIABLE FOR ANY CLAIM,
DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
OTHERWISE, ARISING FROM, OUT OF OR IN ANY WAY CONNECTION WITH THE
LICENSED WORK OR THE USE OR OTHER DEALINGS IN THE LICENSED WORK.
*/

type IMessagePack interface {
	PackAction(Output []byte)
	PackData(msg proto.Message) (data []byte, err error)
	UnPackAction(InData []byte) (pos int32, err error)
	UnPackData() (msg proto.Message, cb reflect.Value, err error)
	GetMessageID() (mainID int32, subID int32)
	Clean()
}

/*
	func: EncodeCmd
	purpose: Encode message mainid and subid to cmd.
*/
func EncodeCmd(mainID, subID uint16) uint32 {
	return (uint32(mainID) << 16) | uint32(subID)
}

/*
	func: DecodeCmd
	purpose: DecodeCmd message cmd to mainid and subid.
*/
func DecodeCmd(cmd uint16) (uint16, uint16) {
	return uint16(cmd >> 16), uint16(cmd)
}

/*

- 协议格式，小端字节序

主命令号 | 次命令号 | 包长度 | 包体
--------| --------| -------- | ----
2字节   |  2字节  | 4字节  | 包体
*/
const (
	EnMessage_MainIDPackLen = 2 //主命令
	EnMessage_SubIDPackLen  = 2 //次命令
	EnMessage_DataPackLen   = 4 //真实数据长度
	EnMessage_NoDataLen     = 8 //非data数据长度(包体之前的)
)

// session, data, data len
type MessageCb func(c net.Conn, mainID int32, subID int32, msg proto.Message)
