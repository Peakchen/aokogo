package logic

import "common/tcpNet"

/*
	by stefan
	date: 20191111 16:08
*/
type ILogicReady interface {
	EnterReady(session tcpNet.TcpSession)
	LeaveReady(session tcpNet.TcpSession)
	ReconnectReady(session tcpNet.TcpSession)
}

var (
	GEnterReadyModule map[string]ILogicReady = map[string]ILogicReady{}
)

var (
	GLeaveReadyModule map[string]ILogicReady = map[string]ILogicReady{}
)

var (
	GReconnReadyModule map[string]ILogicReady = map[string]ILogicReady{}
)

func RegisterEnterReadyModule(module string, data ILogicReady) {
	GEnterReadyModule[module] = data
}

func RegisterReconnReadyModule(module string, data ILogicReady) {
	GReconnReadyModule[module] = data
}

func RegisterLeaveReadyModule(module string, data ILogicReady) {
	GLeaveReadyModule[module] = data
}
