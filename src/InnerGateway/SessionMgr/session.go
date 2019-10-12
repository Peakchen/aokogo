package SessionMgr

import (
	"common/tcpNet"
	"sync"
)

var (
	GInnerGateWaySession *TInnerGateWaySession
)

type TInnerGateWaySession struct {
	client2SeverSession sync.Map
}

func (this *TInnerGateWaySession) AddSessionByID(session *tcpNet.TcpSession, cmd []int32) {
	this.client2SeverSession.Store(session.SessionID, cmd)
}

func (this *TInnerGateWaySession) AddSessionBycmd(session *tcpNet.TcpSession, cmds []int32) {
	for _, cmd := range cmds {
		this.client2SeverSession.Store(cmd, session)
	}
}

func (this *TInnerGateWaySession) RemoveByID(session *tcpNet.TcpSession) {
	this.client2SeverSession.Delete(session.SessionID)
}

func (this *TInnerGateWaySession) RemoveByCmd(cmd int32) {
	this.client2SeverSession.Delete(cmd)
}

func init() {
	GInnerGateWaySession = &TInnerGateWaySession{}
}
