package SessionMgr

import (
	"common/tcpNet"
	"sync"
)

var (
	GServer2ServerSession *TServer2ServerSession
)

type TServer2ServerSession struct {
	s2sSession sync.Map
}

func (this *TServer2ServerSession) AddSessionByID(session *tcpNet.TcpSession, cmd []uint32) {
	this.s2sSession.Store(session.SessionID, cmd)
}

func (this *TServer2ServerSession) AddSessionByCmd(session *tcpNet.TcpSession, cmds []uint32) {
	for _, cmd := range cmds {
		this.s2sSession.Store(cmd, session)
	}
}

func (this *TServer2ServerSession) RemoveByID(session *tcpNet.TcpSession) {
	this.s2sSession.Delete(session.SessionID)
}

func (this *TServer2ServerSession) RemoveByCmd(cmd uint32) {
	this.s2sSession.Delete(cmd)
}

func (this *TServer2ServerSession) GetByCmd(cmd uint32) (session *tcpNet.TcpSession) {
	val, exist := this.s2sSession.Load(cmd)
	if exist {
		session = val.(*tcpNet.TcpSession)
	}
	return
}

func (this *TServer2ServerSession) GetBySessionID(sessionID uint64) (session *tcpNet.TcpSession) {
	val, exist := this.s2sSession.Load(sessionID)
	if exist {
		session = val.(*tcpNet.TcpSession)
	}
	return
}

func init() {
	GServer2ServerSession = &TServer2ServerSession{}
}
