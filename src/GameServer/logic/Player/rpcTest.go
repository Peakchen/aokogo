package Player

/*
	rpc for test about module and single func.
	date: 20191203 by stefan
*/
import (
	"GameServer/rpc"
	"common/Log"
	"common/tcpNet"
)

const (
	cstRpcModule_GetPlayerInfo = string("PlayerInfo")
	cstRpcFunc_GetPlayerInfo   = string("GetPlayerInfo")
)

type TPlayerInfoTest struct {
	Content string
}

/*
	module rpc
*/
type TPlayerUpdateRpc struct {
}

func (this *TPlayerUpdateRpc) GetPlayerInfo(info *TPlayerInfoTest) {
	Log.FmtPrintln("recv module rpc msg, info content: ", info.Content)
}

func RunModuleRpc4GetPlayerInfoTest(session tcpNet.TcpSession, module, funcName string) {
	info := &TPlayerInfoTest{
		Content: "hi，stefan.",
	}

	rpc.SendRpcMsg(session, module, funcName, info)
}

func RunRpc4GetPlayerInfoTest(session tcpNet.TcpSession, funcName string) {
	info := &TPlayerInfoTest{
		Content: "hi，stefan.",
	}

	rpc.SendRpcMsg(session, "", funcName, info)
}

func GetPlayerInfo(info *TPlayerInfoTest) {
	Log.FmtPrintln("recv single rpc msg, info content: ", info.Content)
}

func init() {
	rpc.RegisterRpc(cstRpcFunc_GetPlayerInfo, GetPlayerInfo)
	rpc.RegisterModuleRpc(cstRpcModule_GetPlayerInfo, &TPlayerUpdateRpc{})
}
