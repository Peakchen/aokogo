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

package websockNet

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	Session map[*Session]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Session

	// Unregister requests from clients.
	unregister chan *Session
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Session),
		unregister: make(chan *Session),
		Session:    make(map[*Session]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case Session := <-h.register:
			h.Session[Session] = true
		case Session := <-h.unregister:
			if _, ok := h.Session[Session]; ok {
				delete(h.Session, Session)
				close(Session.send)
			}
		case message := <-h.broadcast:
			for Session := range h.Session {
				select {
				case Session.send <- message:
				default:
					close(Session.send)
					delete(h.Session, Session)
				}
			}
		}
	}
}
