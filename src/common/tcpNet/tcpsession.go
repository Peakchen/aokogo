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

package tcpNet

import(
	"net"
	"time"
	"log"
	"common/S2SMessage"
	"sync"
	"fmt"
	"context"
)

type TcpSession struct{
	host 	string
	isAlive bool
	// The net connection.
	conn 	net.Conn
	// Buffered channel of outbound messages.
	send 	chan []byte
	// send/recv 
	sw  	sync.WaitGroup
	ctx 	context.Context
	// source server or client.
	srcSvr  int32	
	// destination  server or client.
	dstSvr  int32
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

func (c* TcpSession) Connect(){
	if !c.isAlive {
		var err error
		var server_host string = c.host
		c.conn, err = net.Dial("tcp", server_host)
		if err != nil {
			return
		}

		c.isAlive = true
	}

}

func NewSession(addr string, c net.Conn, ctx context.Context, srcSvr, dstSvr int32)*TcpSession{
	return &TcpSession{
		host: addr,
		conn: c,
		send: make(chan []byte, 4096),
		isAlive: false,
		ctx: ctx,
		srcSvr:	srcSvr,
		dstSvr: dstSvr,
	}
}

func (c* TcpSession) exit(){
	c.conn.Close()
	c.sw.Wait()
}

func (c* TcpSession) SetSendCache(data []byte) {
	c.send <- data
}

func (c* TcpSession) Sendmessage(sw *sync.WaitGroup){
	//ticker := time.NewTicker(pingPeriod)
	defer func() {
		//ticker.Stop()
		c.conn.Close()
		sw.Done()
	}()

	for {
		if !c.isAlive {
			c.Connect()
		}

		select {
		case <-c.ctx.Done():
			c.exit()
			return
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.isAlive = false
				c.conn.Close()
				return
			}
			var err error
			_, err = c.conn.Write(message)
			if err != nil {
				c.isAlive = false
				c.conn.Close()
				continue
			}
		}
	}
}

func (c* TcpSession) Recvmessage(sw *sync.WaitGroup){
	defer func() {
		c.conn.Close()
		sw.Done()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	buff := make([]byte, maxMessageSize)
	for {
		len, err := c.conn.Read(buff)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		fmt.Printf("recv mesg: %v.\n", string(buff))
		if self.dstSvr == ER_Client {
			
		}else{
			S2SMessage.DispatchMessage(buff[:len], c.conn, self.srcSvr, self.dstSvr)
		}
	}
}

func (c *TcpSession) HandleSession(){
	c.sw.Add(1)
	go c.Recvmessage(&c.sw)
	c.sw.Add(1)
	go c.Sendmessage(&c.sw)
}