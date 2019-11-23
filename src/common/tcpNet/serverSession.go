package tcpNet

import (
	"common/Define"
	"sync"
)

var (
	GServer2ServerSession *TServer2ServerSession
)

type TServer2ServerSession struct {
	sync.Mutex

	s2sSession sync.Map
}

func (this *TServer2ServerSession) RemoveSessionByID(session *TcpSession) {
	this.Lock()
	defer this.Unlock()

	this.s2sSession.Delete(session.SessionID)
}

func (this *TServer2ServerSession) AddSession(key interface{}, session *TcpSession) {
	this.Lock()
	defer this.Unlock()

	this.s2sSession.Store(key, session)
}

func (this *TServer2ServerSession) GetSessionByType(RegPoint Define.ERouteId) (session *TcpSession) {
	this.Lock()
	defer this.Unlock()

	val, exist := this.s2sSession.Load(RegPoint)
	if exist {
		session = val.(*TcpSession)
	}
	return
}

func (this *TServer2ServerSession) RemoveSessionByType(RegPoint Define.ERouteId) {
	this.Lock()
	defer this.Unlock()

	this.s2sSession.Delete(RegPoint)
}

func (this *TServer2ServerSession) GetSessionByModuleID(moduleID uint16) (session *TcpSession) {
	this.Lock()
	defer this.Unlock()

	val, exist := this.s2sSession.Load(moduleID)
	if exist {
		session = val.(*TcpSession)
	}
	return
}

func init() {
	GServer2ServerSession = &TServer2ServerSession{}
}
