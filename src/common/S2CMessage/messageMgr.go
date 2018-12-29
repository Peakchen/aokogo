package S2CMessage

import (
	"strconv"
	"github.com/protobuf/proto"
)

var tMessageMap map[string]MessageCallBack = map[string]MessageCallBack{}

type MessageCallBack func(session string, msg *proto.Message)(error, bool)

func BuildMessageKey(){
	
}

func SendMessage2c(Session string, MainID, SubID uint32, msg *proto.Message){
	mkey := strconv.Itoa(int(MainID)) + "." + strconv.Itoa(int(SubID))
	_, ok := tMessageMap[mkey]
	if !ok{
		panic("[SendMessage2c] can not find message key from messagemap, MainID: %v, SubID: %v.", MainID, SubID)
		return
	}

	
}

func RegisterMessageCallBack(MainID, SubID uint32, cb MessageCallBack){ 
	mkey := strconv.Itoa(int(MainID)) + "." + strconv.Itoa(int(SubID))
	tMessageMap[mkey] = cb
}