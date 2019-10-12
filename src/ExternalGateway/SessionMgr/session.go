package SessionMgr

import (
	"common/tcpNet"
	"sync"
)

var (
	GExternalGateWaySession *TExternalGateWaySession
)

type TExternalGateWaySession struct {
	client2SeverSession sync.Map
}

func (this *TExternalGateWaySession) AddSessionByID(session *tcpNet.TcpSession, cmd []int32) {
	this.client2SeverSession.Store(session.SessionID, cmd)
}

func (this *TExternalGateWaySession) AddSessionBycmd(session *tcpNet.TcpSession, cmds []int32) {
	for _, cmd := range cmds {
		this.client2SeverSession.Store(cmd, session)
	}
}

func (this *TExternalGateWaySession) RemoveByID(session *tcpNet.TcpSession) {
	this.client2SeverSession.Delete(session.SessionID)
}

func (this *TExternalGateWaySession) RemoveByCmd(cmd int32) {
	this.client2SeverSession.Delete(cmd)
}

func init() {
	GExternalGateWaySession = &TExternalGateWaySession{}
}
