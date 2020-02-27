// add by stefan

package tcpWebNet

import (
	//"bytes"

	"common/aktime"
	"log"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/gorilla/websocket"
)

// TWebsocketSession is a middleman between the websocket connection and the hub.
type TWebsocketSession struct {
	hub *TWebSockHub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

/*
	@purpose: receive user message loop.
	@param1: no
*/
func (c *TWebsocketSession) recvloop() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(aktime.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(aktime.Now().Add(pongWait)); return nil })

	for {
		mt, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			continue
		}

		if mt == websocket.TextMessage {
			if !utf8.Valid(message) {
				c.conn.WriteControl(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseInvalidFramePayloadData, ""),
					time.Time{})
				log.Println("ReadAll: invalid utf8")
			}
		} else if mt == websocket.BinaryMessage {
			// then dispatch message to msg deal with module
			//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
			//C2SMessage.DispatchMessage(message, c.conn)
		} else if mt == websocket.PingMessage {
			// todo:
		} else if mt == websocket.PongMessage {
			// todo:
		}
	}
}

/*
	@purpose: send user message loop.
	@param1: no
*/
func (c *TWebsocketSession) sendloop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(aktime.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(aktime.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
