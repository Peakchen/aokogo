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

package tcpNet

import(
	"net"
	"time"
	"log"
	"common/S2SMessage"
	"sync"
)

type TcpSession struct{
	host 	string
	isAlive bool
	// The net connection.
	conn 	net.Conn
	// Buffered channel of outbound messages.
	send 	chan []byte
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

func (c* TcpSession) connect(){
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

func NewSession(addr string, c net.Conn)*TcpSession{
	return &TcpSession{
		host: addr,
		conn: c,
		send: make(chan []byte, 4096),
	}
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
			c.connect()
		}

		select {
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
		S2SMessage.DispatchMessage(buff[0:len], c.conn)
	}
}
