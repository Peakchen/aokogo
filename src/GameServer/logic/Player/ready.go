package Player

import (
	"GameServer/logic"
	"common/Log"
	"common/akNet"
)

type TPlayerReady struct {
}

func (this *TPlayerReady) EnterReady(session akNet.TcpSession) {
	Log.FmtPrintln("enter ready.")
	player := GetPlayer(session.GetIdentify())
	if player == nil {
		Log.Error("can not find ")
		return
	}

	//for test
	//RunModuleRpc4GetPlayerInfoTest(session, cstRpcModule_GetPlayerInfo, cstRpcFunc_GetPlayerInfo)
	//RunRpc4GetPlayerInfoTest(session, cstRpcFunc_GetPlayerInfo)
}

func (this *TPlayerReady) LeaveReady(session akNet.TcpSession) {

}

func (this *TPlayerReady) ReconnectReady(session akNet.TcpSession) {

}

func init() {
	logic.RegisterEnterReadyModule(cstPlayerSubModule, &TPlayerReady{})
	logic.RegisterReconnReadyModule(cstPlayerSubModule, &TPlayerReady{})
	logic.RegisterLeaveReadyModule(cstPlayerSubModule, &TPlayerReady{})
}
