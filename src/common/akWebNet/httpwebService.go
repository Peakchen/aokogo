// add by stefan

package tcpWebNet

import (
	//"fmt"
	"common/Log"
	"log"
	"net/http"
)

func StartWebSockSvr(addr string) bool {
	var (
		hub = newHub()
	)

	go hub.run()
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		Log.Error("ListenAndServe websock listen fail, addr: ", addr, "err: ", err)
		return false
	}
	return true
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *TWebSockHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	session := &TWebsocketSession{hub: hub, conn: conn, send: make(chan []byte, maxMessageSize)}
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go session.sendloop()
	go session.recvloop()

	session.hub.register <- session
}
