package GameLog

import (
	"common/Log"
)

const (
	Const_GameServerLog string = "GameServer"
)

func init(){
	Log.NewLog(Const_GameServerLog)
}