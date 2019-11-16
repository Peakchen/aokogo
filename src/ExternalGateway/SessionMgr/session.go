package SessionMgr

import (
	"common/Define"
	"common/Log"
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
	mf := func(k, v interface{}) bool {
		Log.FmtPrintf("client to server, key: %v.", k)
		return true
	}

	this.c2sSession.Range(mf)
	val, exist := this.c2sSession.Load(cmd)
	if exist {
		session = val.(*tcpNet.TcpSession)
	}
	return
}

func (this *TClient2ServerSession) GetSessionByID(sessionID uint64) (session *tcpNet.TcpSession) {
	val, exist := this.c2sSession.Load(sessionID)
	if exist {
		session = val.(*tcpNet.TcpSession)
	}
	return
}

func (this *TClient2ServerSession) AddSession(session *tcpNet.TcpSession) {
	this.c2sSession.Store(session.SessionID, session)
	this.c2sSession.Store(session.RegPoint, session)
}

func (this *TClient2ServerSession) GetSessionByType(RegPoint Define.ERouteId) (session *tcpNet.TcpSession) {
	val, exist := this.c2sSession.Load(RegPoint)
	if exist {
		session = val.(*tcpNet.TcpSession)
	}
	return
}

func (this *TClient2ServerSession) RemoveSessionByType(RegPoint Define.ERouteId) {
	this.c2sSession.Delete(RegPoint)
}

func init() {
	GClient2ServerSession = &TClient2ServerSession{}
}
