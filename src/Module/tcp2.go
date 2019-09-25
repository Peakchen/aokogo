package main

import (
	"common/msgProto/MSG_Server"
	"common/tcpNet"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
)

//发送信息
func sender(conn net.Conn) {
	//words := "Hello Server!"
	cp := tcpNet.ClientProtocol{}

	req := &MSG_Server.CS_EnterServer_Req{}
	req.Enter = 2
	data, err := proto.Marshal(req)
	if err != nil {
		fmt.Println("proto marshal fail, data: ", err)
		return
	}
	cp.SetCmd(1, 1, data)
	buff := make([]byte, 512)
	cp.PackAction(buff)
	conn.Write(buff) //[]byte(words)
	fmt.Println("send over")

	//接收服务端反馈
	buffer := make([]byte, 2048)

	n, err := conn.Read(buffer)
	if err != nil {
		Log(conn.RemoteAddr().String(), "waiting server back msg error: ", err)
		return
	}
	Log(conn.RemoteAddr().String(), "receive server back msg: ", string(buffer[:n]))

}

//日志
func Log(v ...interface{}) {
	log.Println(v...)
}

func main() {
	for i := 0; i < 100; i++ {
		time.Sleep(time.Duration(2))
		dialsend()
	}
}

func dialsend() {
	server := "0.0.0.0:51001"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("connection success")
	sender(conn)
}
