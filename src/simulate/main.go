package main

import (
	"common/Log"
	"common/msgProto/MSG_Server"
	"common/tcpNet"
	"fmt"
	"net"
	"os"
	"sync"
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

}

func readloop(conn net.Conn) {
	for {
		//接收服务端反馈
		buffer := make([]byte, 2048)
		n, err := conn.Read(buffer)
		if err != nil {
			Log.FmtPrintln("waiting server back msg error: ", conn.RemoteAddr().String(), err)
			continue
		}
		Log.FmtPrintln("receive server back msg: ", conn.RemoteAddr().String(), string(buffer[:n]))
	}

}

var sw sync.WaitGroup

func main() {
	for i := 0; i < 100; i++ {
		Log.FmtPrintln("time: ", i)
		time.Sleep(time.Duration(1))
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
	sw.Add(1)
	go readloop(conn)
	sw.Wait()
}
