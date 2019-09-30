package main

import (
	"common/Log"
	"common/msgProto/MSG_Server"
	"common/tcpNet"
	"net"
	"os"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
)

//发送信息
func sender(conn net.Conn) bool {
	//words := "Hello Server!"
	cp := tcpNet.ClientProtocol{}
	req := &MSG_Server.CS_EnterServer_Req{}
	req.Enter = 2
	data, err := proto.Marshal(req)
	if err != nil {
		Log.FmtPrintln("proto marshal fail, data: ", err)
		return false
	}
	cp.SetCmd(1, 1, data)
	buff := make([]byte, len(data)+8)
	cp.PackAction(buff)
	Log.FmtPrintln("send buff len: ", len(buff))
	n, err := conn.Write(buff)
	if n == 0 || err != nil {
		Log.FmtPrintln("Write fail, data: ", n, err)
		return false
	}
	Log.FmtPrintln("send over")
	return true
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
		Log.FmtPrintf("receive server back, ip: %v, msg: %v.", conn.RemoteAddr().String(), string(buffer[:n]))
	}

}

var sw sync.WaitGroup

func main() {
	dialsend()
}

var server string = "0.0.0.0:51001"

func sendloop(conn net.Conn) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		Log.FmtPrintf("Fatal error: %s", err.Error())
		os.Exit(1)
	}
	for i := 0; i < 10; i++ {
		Log.FmtPrintln("time: ", i)
		if !sender(conn) {
			conn, err = net.DialTCP("tcp", nil, tcpAddr)
			if err != nil {
				Log.FmtPrintf("Fatal error: %s", err.Error())
				os.Exit(1)
			}
		}
		time.Sleep(time.Duration(2) * time.Second)
	}
}

func dialsend() {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		Log.FmtPrintf("Fatal error: %s", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		Log.FmtPrintf("Fatal error: %s", err.Error())
		os.Exit(1)
	}

	Log.FmtPrintln("connection success")
	sw.Add(2)
	go readloop(conn)
	go sendloop(conn)
	sw.Wait()
}
