/*
Copyright (this) <year> <copyright holders>

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

package tcpNet

import (
	"common/Log"
	"encoding/binary"
	"io"
	"net"
	"reflect"
	"time"

	//"common/S2SMessage"
	"context"
	"sync"

	"github.com/golang/protobuf/proto"
	//. "common/Define"
)

// session, data, data len
type MessageCb func(this net.Conn, mainID int32, subID int32, msg proto.Message)

type TcpSession struct {
	host    string
	isAlive bool
	// The net connection.
	conn *net.TCPConn
	// Buffered channel of outbound messages.
	send chan []byte
	// send/recv
	sw  sync.WaitGroup
	ctx context.Context
	// source server or client/destination  server or client.
	mapSvr map[int32][]int32
	// receive message call back
	recvCb MessageCb
	// person offline flag
	off chan *TcpSession
	//message pack
	pack IMessagePack
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 4096
)

func (this *TcpSession) Connect() {
	if !this.isAlive {
		tcpAddr, err := net.ResolveTCPAddr("tcp", this.host)
		if err != nil {
			Log.FmtPrintln("session failed: ", err)
			return
		}

		this.conn, err = net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			return
		}

		this.isAlive = true
	}

}

func NewSession(addr string,
	this *net.TCPConn,
	ctx context.Context,
	mapSvr *map[int32][]int32,
	newcb MessageCb,
	off chan *TcpSession,
	pack IMessagePack) *TcpSession {
	return &TcpSession{
		host:    addr,
		conn:    this,
		send:    make(chan []byte, maxMessageSize),
		isAlive: false,
		ctx:     ctx,
		mapSvr:  *mapSvr,
		recvCb:  newcb,
		pack:    pack,
	}
}

func (this *TcpSession) exit() {
	this.off <- this
	close(this.send)
	this.conn.CloseRead()
	this.conn.CloseWrite()
	this.conn.Close()
	this.sw.Wait()
}

func (this *TcpSession) SetSendCache(data []byte) {
	this.send <- data
}

func (this *TcpSession) Sendmessage(sw *sync.WaitGroup) {
	defer func() {
		this.exit()
	}()

	for {
		select {
		case <-this.ctx.Done():
			this.exit()
			return
		case data, ok := <-this.send:
			if !ok {
				// The hub closed the channel.
				this.isAlive = false
				this.conn.Close()
				return
			}
			this.writeMessage(data)
		}
	}
}

func (this *TcpSession) Recvmessage(sw *sync.WaitGroup) {
	defer func() {
		this.exit()
	}()

	for {
		select {
		case <-this.ctx.Done():
			return
		default:
			this.readMessage()
		}
	}
}

func (this *TcpSession) writeMessage(data []byte) {
	this.conn.SetWriteDeadline(time.Now().Add(writeWait))
	//pack message then send.

	//send...
	var err error
	Log.FmtPrintln("begin send response message to client.")
	_, err = this.conn.Write(data)
	if err != nil {
		this.isAlive = false
		this.conn.Close()
	}
}

func (this *TcpSession) readMessage() {
	this.conn.SetReadDeadline(time.Now().Add(pongWait))
	packLenBuf := make([]byte, EnMessage_NoDataLen)
	readn, err := io.ReadFull(this.conn, packLenBuf)
	if err != nil || readn < EnMessage_NoDataLen {
		Log.FmtPrintln("read err:", err)
		return
	}

	packlen := binary.LittleEndian.Uint32(packLenBuf[EnMessage_DataPackLen:EnMessage_NoDataLen])
	if packlen > maxMessageSize {
		Log.FmtPrintln("error receiving packLen:", packlen)
		return
	}

	data := make([]byte, EnMessage_NoDataLen+packlen)
	readn, err = io.ReadFull(this.conn, data[EnMessage_NoDataLen:])
	if err != nil || readn < int(packlen) {
		Log.FmtPrintln("error receiving msg, readn:", readn, "packLen:", packlen, "reason:", err)
		return
	}

	//todo: unpack message then read real date.
	copy(data[:EnMessage_NoDataLen], packLenBuf[:])
	_, err = this.pack.UnPackAction(data)
	if err != nil {
		Log.FmtPrintln("unpack action err: ", err)
		return
	}

	msg, cb, err := this.pack.UnPackData()
	if err != nil {
		Log.FmtPrintln("unpack data err: ", err)
		return
	}

	mainID, subID := this.pack.GetMessageID()
	Log.FmtPrintf("mainid: %v, subID: %v.", mainID, subID)
	//read...
	params := []reflect.Value{
		reflect.ValueOf("1"),
		reflect.ValueOf(msg),
	}
	cb.Call(params)
	this.SetSendCache([]byte("respone client.\n"))
	this.recvCb(this.conn, mainID, subID, msg)
}

func (this *TcpSession) HandleSession() {
	this.sw.Add(1)
	go this.Recvmessage(&this.sw)
	this.sw.Add(1)
	go this.Sendmessage(&this.sw)
	//this.sw.Wait()
}
