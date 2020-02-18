// add by stefan

package Player

import (
	"GameServer/dbo"
	"GameServer/logic/LogicDef"
	"common/Log"
	"common/ado"
	"common/msgProto/MSG_Player"
)

/*
	player struct.
*/
type TPlayer struct {
	ado.IDBModule

	BaseInfo  map[MSG_Player.EmBaseInfo]interface{}  //基础玩家信息
	BaseMoney map[MSG_Player.EmBaseMoney]interface{} //基础金币类信息
}

func (this *TPlayer) Identify() string {
	return this.StrIdentify
}

func (this *TPlayer) MainModel() string {
	return LogicDef.CstUsrDataCenter
}

func (this *TPlayer) SubModel() string {
	return cstPlayerSubModule
}

func GetPlayer(Identify string) (player *TPlayer) {
	player = &TPlayer{}
	err, exist := dbo.A_DBRead(Identify, player)
	if err != nil {
		Log.Error("can not read player data, err: ", err)
		return
	}

	if !exist {
		// ... data init, then insert cache and db.
		player.StrIdentify = Identify
		player.BaseInfo = map[MSG_Player.EmBaseInfo]interface{}{}
		player.BaseMoney = map[MSG_Player.EmBaseMoney]interface{}{}
		player.initdata()
		err = dbo.A_DBInsert(Identify, player)
		if err != nil {
			Log.Error("can not insert player data, err: ", err)
			return
		}
	}
	return
}

func (this *TPlayer) initdata() {
	// base info
	this.BaseInfo[MSG_Player.EmBaseInfo_Name] = "嘿嘿嘿"
	this.BaseInfo[MSG_Player.EmBaseInfo_Level] = 1
	this.BaseInfo[MSG_Player.EmBaseInfo_HeadIcon] = 101
	this.BaseInfo[MSG_Player.EmBaseInfo_DBID] = 1001

	// base money
	this.BaseMoney[MSG_Player.EmBaseMoney_Coin] = 3000
}
