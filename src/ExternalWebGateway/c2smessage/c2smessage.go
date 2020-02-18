// add by stefan

package c2smessage

import (
	"common/C2SMessage"
	"log"

	"github.com/gorilla/websocket"
)

/* login game by account and passwd
 */
func OnC2SDispatchMessage(basemsg *C2SMessage.CS_BaseMessage_Req, c *websocket.Conn) {
	// todo:
	log.Printf("OnC2SDispatchMessage: player Sid[%d].", basemsg.Sid)
	// switch S2SMessage.ServerId(basemsg.Sid) {
	// case S2SMessage.ServerId_SID_Login:
	// 	//send to loginserver
	// 	S2SMessage.DispatchClientMessage(basemsg.Data, c)
	// case S2SMessage.ServerId_SID_BigWorld:

	// case S2SMessage.ServerId_SID_Game:

	// case S2SMessage.ServerId_SID_SmallWorld:

	// }
}

func init() {
	//C2SMessage.Register(int32(C2SMessage.MessageRoute_R_SID_ESG), OnC2SDispatchMessage)
}
