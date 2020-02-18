// add by stefan

package tcpWebNet

var (
	mapSession map[int64]*TWebsocketSession
)

func register(sessionid int64, s *TWebsocketSession) {
	mapSession[sessionid] = s
}

func deleteSession(sessionid int64) {
	delete(mapSession, sessionid)
}

func getSession(sessionid int64) *TWebsocketSession {
	var s, ok = mapSession[sessionid]
	if !ok {
		return nil
	}

	return s
}
