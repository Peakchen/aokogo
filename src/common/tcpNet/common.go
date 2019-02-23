package tcpNet

/*
* CopyRight(C) StefanChen e-mail:2572915286@qq.com
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/


import (
	"encoding/binary"
)

/*
	model: ServerProtocol
	server to server, message  
*/
type ServerProtocol struct{
	mainid uint16
	subid  uint16
	messagelength uint16
	messageData []byte
}

func (self *ServerProtocol)S2SPack(Output []byte){
	var pos int32 = 0
	binary.LittleEndian.PutUint16(Output[pos:], self.mainid)
	pos += 2

	binary.LittleEndian.PutUint16(Output[pos:], self.subid)
	pos += 2

	binary.LittleEndian.PutUint16(Output[pos:], self.messagelength)
	pos += 2

	copy(Output[pos:], self.messageData)
}

func (self *ServerProtocol)S2SUnPack(InData []byte)(int32){
	var pos int32 = 0
	self.mainid = binary.LittleEndian.Uint16(InData[pos:])
	pos += 2

	self.subid = binary.LittleEndian.Uint16(InData[pos:])
	pos += 2

	self.messagelength = binary.LittleEndian.Uint16(InData[pos:])
	pos += 2

	pos += int32(self.messagelength)
	return pos
}

/*
	func: EncodeCmd
	purpose: Encode message mainid and subid to cmd.
*/
func EncodeCmd(mainID, subID uint16)uint32{
	return (uint32(mainID) << 16 ) | uint32(subID)
}

/*
	func: DecodeCmd
	purpose: DecodeCmd message cmd to mainid and subid.
*/
func DecodeCmd(cmd uint16)(uint16, uint16){
	return uint16(cmd >> 16 ), uint16(cmd)
}
 