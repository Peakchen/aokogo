// add by stefan

package tcpWebNet

import (
	"context"
)

// TWebSockHub maintains the set of active clients and broadcasts messages to the
// clients.
type TWebSockHub struct {
	// Registered clients.
	wsocketSession map[*TWebsocketSession]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *TWebsocketSession

	// Unregister requests from clients.
	unregister chan *TWebsocketSession

	ctx    context.Context
	cancel context.CancelFunc
}

func newHub() *TWebSockHub {
	return &TWebSockHub{
		broadcast:      make(chan []byte),
		register:       make(chan *TWebsocketSession),
		unregister:     make(chan *TWebsocketSession),
		wsocketSession: make(map[*TWebsocketSession]bool),
	}
}

func (this *TWebSockHub) run() {
	go func() {
		this.exit()
	}()

	this.ctx, this.cancel = context.WithCancel(context.Background())

	for {
		select {
		case <-this.ctx.Done():
			return
		case session := <-this.register:
			this.wsocketSession[session] = true
		case session := <-this.unregister:
			if _, ok := this.wsocketSession[session]; ok {
				delete(this.wsocketSession, session)
				close(session.send)
			}
		case message := <-this.broadcast:
			for session := range this.wsocketSession {
				select {
				case session.send <- message:
					continue
				default:
					close(session.send)
					delete(this.wsocketSession, session)
				}
			}
		}
	}
}

func (this *TWebSockHub) exit() {
	for session := range this.wsocketSession {
		close(session.send)
		delete(this.wsocketSession, session)
	}
	this.cancel()
}
