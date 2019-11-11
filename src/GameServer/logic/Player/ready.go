package Player

import (
	"GameServer/logic"
	"common/Log"
)

type TPlayerReady struct {
}

func (this *TPlayerReady) EnterReady(identify string) {
	player := GetPlayer(identify)
	if player == nil {
		Log.Error("can not find ")
		return
	}
}

func (this *TPlayerReady) LeaveReady(identify string) {

}

func (this *TPlayerReady) ReconnectReady(identify string) {

}

func init() {
	logic.RegisterEnterReadyModule(cstPlayerSubModule, &TPlayerReady{})
	logic.RegisterReconnReadyModule(cstPlayerSubModule, &TPlayerReady{})
	logic.RegisterLeaveReadyModule(cstPlayerSubModule, &TPlayerReady{})
}
