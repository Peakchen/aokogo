package SessionMgr

import (
	"common/tcpNet"
	"sync"
)

var (
	GClient2ServerSession *TClient2ServerSession
)

type TClient2ServerSession struct {
	c2sSession sync.Map
}

func (this *TClient2ServerSession) AddSessionByID(session *tcpNet.TcpSession, cmd []uint32) {
	this.c2sSession.Store(session.SessionID, cmd)
}

func (this *TClient2ServerSession) AddSessionByCmd(session *tcpNet.TcpSession, cmds []uint32) {
	for _, cmd := range cmds {
		this.c2sSession.Store(cmd, session)
	}
}

func (this *TClient2ServerSession) RemoveByID(session *tcpNet.TcpSession) {
	this.c2sSession.Delete(session.SessionID)
}

func (this *TClient2ServerSession) RemoveByCmd(cmd uint32) {
	this.c2sSession.Delete(cmd)
}

func (this *TClient2ServerSession) GetByCmd(cmd uint32) (session *tcpNet.TcpSession) {
	val, exist := this.c2sSession.Load(cmd)
	if exist {
		session = val.(*tcpNet.TcpSession)
	}
	return
}

func (this *TClient2ServerSession) GetBySessionID(sessionID uint64) (session *tcpNet.TcpSession) {
	val, exist := this.c2sSession.Load(sessionID)
	if exist {
		session = val.(*tcpNet.TcpSession)
	}
	return
}

func init() {
	GClient2ServerSession = &TClient2ServerSession{}
}
