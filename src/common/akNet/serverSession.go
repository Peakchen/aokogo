package akNet

// add by stefan

import (
	"common/define"
	"common/utls"
	"sync"
)

var (
	GServer2ServerSession *TSvr2SvrSession
)

type TSvr2SvrSession struct {
	sync.Mutex

	s2sSession sync.Map
}

func (this *TSvr2SvrSession) RemoveSession(key interface{}) {
	this.Lock()
	defer this.Unlock()

	this.s2sSession.Delete(key)
}

func (this *TSvr2SvrSession) AddSession(key interface{}, session TcpSession) {
	this.Lock()
	defer this.Unlock()

	this.s2sSession.Store(key, session)
}

func (this *TSvr2SvrSession) GetSession(key interface{}) (session TcpSession) {
	this.Lock()
	defer this.Unlock()

	var (
		sessions = []TcpSession{}
		slen     int32
		randIdx  int32
	)
	this.s2sSession.Range(func(k, v interface{}) bool {
		cs := v.(TcpSession)
		if cs.GetRegPoint() == key.(define.ERouteId) && cs.Alive() {
			sessions = append(sessions, cs)
		}
		return true
	})

	slen = int32(len(sessions))
	if slen > 1 {
		randIdx = utls.RandInt32FromZero(slen)
	} else if slen == 0 {
		return
	}

	session = sessions[randIdx]
	return
}

func (this *TSvr2SvrSession) GetSessionByIdentify(key interface{}) (session TcpSession) {
	this.Lock()
	defer this.Unlock()

	val, exist := this.s2sSession.Load(key)
	if exist {
		session = val.(TcpSession)
	}
	return
}

func (this *TSvr2SvrSession) GetAllSession() (sessions sync.Map) {
	return this.s2sSession
}

func init() {
	GServer2ServerSession = &TSvr2SvrSession{}
}
