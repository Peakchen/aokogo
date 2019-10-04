package main

import (
	"common/Log"
	"common/msgProto/MSG_Server"
	"common/tcpNet"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

func readloop(conn net.Conn, ctx context.Context, sw *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			sw.Done()
			return
		default:
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

}

var sw sync.WaitGroup

func main() {
	dialsend()
}

func init() {

}

var server string = "0.0.0.0:51001"

func sendloop(conn net.Conn, ctx context.Context, sw *sync.WaitGroup) {
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

var exitchan = make(chan os.Signal, 1)

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

	ctx, _ := context.WithCancel(context.Background())
	Log.FmtPrintln("connection success")
	signal.Notify(exitchan, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGSEGV)

	sw.Add(3)
	go readloop(conn, ctx, &sw)
	go sendloop(conn, ctx, &sw)
	go exitloop(ctx, &sw)
	sw.Wait()
}

func exitloop(ctx context.Context, sw *sync.WaitGroup) {
	for {
		//Block until a signal is received.
		if s, ok := <-exitchan; ok {
			fmt.Println("Got signal:", s)
		}
		os.Exit(1)
		select {
		case <-ctx.Done():
			sw.Done()
			return
		default:

		}
	}
}
