package tcpNet

import (
	"common/Define"
	"sync"
)

var (
	GServer2ServerSession *TServer2ServerSession
)

type TServer2ServerSession struct {
	s2sSession sync.Map
}

func (this *TServer2ServerSession) AddSessionByID(session *TcpSession, cmd []uint32) {
	this.s2sSession.Store(session.SessionID, cmd)
}

func (this *TServer2ServerSession) AddSessionByCmd(session *TcpSession, cmds []uint32) {
	for _, cmd := range cmds {
		this.s2sSession.Store(cmd, session)
	}
}

func (this *TServer2ServerSession) RemoveSessionByID(session *TcpSession) {
	this.s2sSession.Delete(session.SessionID)
}

func (this *TServer2ServerSession) RemoveByCmd(cmd uint32) {
	this.s2sSession.Delete(cmd)
}

func (this *TServer2ServerSession) GetByCmd(cmd uint32) (session *TcpSession) {
	val, exist := this.s2sSession.Load(cmd)
	if exist {
		session = val.(*TcpSession)
	}
	return
}

func (this *TServer2ServerSession) GetSessionByID(sessionID uint64) (session *TcpSession) {
	val, exist := this.s2sSession.Load(sessionID)
	if exist {
		session = val.(*TcpSession)
	}
	return
}

func (this *TServer2ServerSession) AddSession(session *TcpSession) {
	this.s2sSession.Store(session.SessionID, session)
	this.s2sSession.Store(session.RegPoint, session)
}

func (this *TServer2ServerSession) GetSessionByType(RegPoint Define.ERouteId) (session *TcpSession) {
	val, exist := this.s2sSession.Load(RegPoint)
	if exist {
		session = val.(*TcpSession)
	}
	return
}

func (this *TServer2ServerSession) RemoveSessionByType(RegPoint Define.ERouteId) {
	this.s2sSession.Delete(RegPoint)
}

func init() {
	GServer2ServerSession = &TServer2ServerSession{}
}
