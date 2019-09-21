package tcpMessage

import (
	"common/Log"
	"common/msgProto/MSG_Server"
	"common/tcpNet"
	"net"
	"testing"

	"common/Define"

	"github.com/golang/protobuf/proto"
)

func main() {
	Log.FmtPrintf("msg main.")
}

func MsgTest(t *testing.T) {
	Log.FmtPrintf("msg test.")
	var (
		mapsvr map[int32][]int32 = map[int32][]int32{
			int32(Define.ERouteId_ER_Client): []int32{int32(Define.ERouteId_ER_ISG), int32(Define.ERouteId_ER_ESG)},
		}
	)
	client := tcpNet.NewClient("0.0.0.1:19000",
		&mapsvr,
		ClientMessageCallBack)

	client.Run()
	req := &MSG_Server.CS_EnterServer_Req{}
	req.Enter = 1
	data, err := proto.Marshal(req)
	if err != nil {
		Log.FmtPrintf("proto marshal fail, data: ", err)
		return
	}
	client.Send(data)
}

func ClientMessageCallBack(c net.Conn, data []byte, len int) {
	Log.FmtPrintf("exec client message call back.", c.RemoteAddr(), c.LocalAddr())
}
