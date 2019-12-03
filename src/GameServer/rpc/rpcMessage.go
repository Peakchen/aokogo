package rpc

/*
	rpc process message module
	date: 20191203
	author: stefan
	version: 1.0
*/

import (
	"common/Log"
	"common/msgProto/MSG_MainModule"
	"common/msgProto/MSG_Rpc"
	"common/tcpNet"
)

func Init() {

}

/*
	@func: onRpcProcess 接收处理rpc消息
	@param1: session obj
	@param2: req content (module, func, data)
*/
func onRpcProcess(session tcpNet.TcpSession, req *MSG_Rpc.CS_Rpc_Req) (succ bool, err error) {
	Log.FmtPrintf("rpc process, rpc module: %v, func: %v.", req.Rpcmodule, req.Rpcfunc)
	if len(req.Rpcmodule) == 0 {
		succ, err = onSingleRpc(req.Rpcfunc, req.Data)
	} else {
		succ, err = onModuleRpcProcess(req.Rpcmodule, req.Rpcfunc, req.Data)
	}
	return
}

func init() {
	tcpNet.RegisterMessage(uint16(MSG_MainModule.MAINMSG_RPC), uint16(MSG_Rpc.SUBMSG_CS_Rpc), onRpcProcess)
}
