package M_tcp

import (
	"common/Log"
	"common/msgProto/MSG_Server"
	"common/tcpNet"
	"net"
	"testing"

	"common/Define"

	"github.com/golang/protobuf/proto"
)

func TestMsg(t *testing.T) {
	Log.FmtPrintf("msg test.")
	var (
		mapsvr map[int32][]int32 = map[int32][]int32{
			int32(Define.ERouteId_ER_Client): []int32{int32(Define.ERouteId_ER_ISG), int32(Define.ERouteId_ER_ESG)},
		}
	)
	client := tcpNet.NewClient(Define.GameServerHost,
		&mapsvr,
		ClientMessageCallBack)

	client.Run()
	data, err := getServerData()
	if err != nil {
		return
	}
	client.Send(data)
}

func getServerData() (data []byte, err error) {
	req := &MSG_Server.CS_EnterServer_Req{}
	req.Enter = 1
	data, err = proto.Marshal(req)
	if err != nil {
		Log.Error("proto marshal fail, data: ", err)
	}
	return
}

func ClientMessageCallBack(c net.Conn, mainID int32, subID int32, msg proto.Message) {
	Log.FmtPrintf("exec client message call back.", c.RemoteAddr(), c.LocalAddr())
}
