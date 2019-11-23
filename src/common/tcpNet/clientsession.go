package tcpNet

import (
	"common/Define"
	"sync"
)

var (
	GClient2ServerSession *TClient2ServerSession
)

type TClient2ServerSession struct {
	sync.Mutex

	c2sSession sync.Map
}

func (this *TClient2ServerSession) AddSessionByCmd(session *TcpSession, cmds []uint32) {
	this.Lock()
	defer this.Unlock()
	for _, cmd := range cmds {
		this.c2sSession.Store(cmd, session)
	}
}

func (this *TClient2ServerSession) RemoveSessionByID(session *TcpSession) {
	this.Lock()
	defer this.Unlock()
	this.c2sSession.Delete(session.SessionID)
}

func (this *TClient2ServerSession) AddSession(key interface{}, session *TcpSession) {
	this.Lock()
	defer this.Unlock()

	this.c2sSession.Store(key, session)
}

func (this *TClient2ServerSession) GetSessionByType(RegPoint Define.ERouteId) (session *TcpSession) {
	this.Lock()
	defer this.Unlock()
	val, exist := this.c2sSession.Load(RegPoint)
	if exist {
		session = val.(*TcpSession)
	}
	return
}

func (this *TClient2ServerSession) RemoveSessionByType(RegPoint Define.ERouteId) {
	this.Lock()
	defer this.Unlock()
	this.c2sSession.Delete(RegPoint)
}

func (this *TClient2ServerSession) AddSessionByModuleID(moduleID uint16, session *TcpSession) {
	this.Lock()
	defer this.Unlock()
	this.c2sSession.Store(moduleID, session)
}

func (this *TClient2ServerSession) GetSessionByModuleID(moduleID uint16) (session *TcpSession) {
	this.Lock()
	defer this.Unlock()
	val, exist := this.c2sSession.Load(moduleID)
	if exist {
		session = val.(*TcpSession)
	}
	return
}

func init() {
	GClient2ServerSession = &TClient2ServerSession{}
}
