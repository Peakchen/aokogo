package main

import (
	"common/Define"
	"common/msgProto/MSG_Server"
	"common/tcpNet"
	"fmt"
	"net"
	"testing"

	"github.com/golang/protobuf/proto"
)

func Test1Tcp(t *testing.T) {
	t.Log("[Test_tcp_1] start.")
	cmd := 65546
	c := uint16(cmd)
	a := uint16(cmd >> 16)
	b := uint16(cmd)

	fmt.Printf("a: %v, b: %v, c: %v.\n", a, b, c)
}

func Test2Tcp(t *testing.T) {
	t.Log("TestTcp2...")
	var (
		mapsvr map[int32][]int32 = map[int32][]int32{
			int32(Define.ERouteId_ER_Client): []int32{int32(Define.ERouteId_ER_ISG), int32(Define.ERouteId_ER_ESG)},
		}
	)
	client := tcpNet.NewClient(Define.ExternalServerHost,
		&mapsvr,
		ClientMessageCallBack)

	client.Run()
	req := &MSG_Server.CS_EnterServer_Req{}
	req.Enter = 1
	data, err := proto.Marshal(req)
	if err != nil {
		t.Log("proto marshal fail, data: ", err)
		return
	}
	client.Send(data)
}

func ClientMessageCallBack(c net.Conn, data []byte, len int) {
	fmt.Println("exec client message call back.", c.RemoteAddr(), c.LocalAddr())
}

func init() {

}
