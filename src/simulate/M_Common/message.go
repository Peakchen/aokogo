package M_Common

import (
	"common/Log"
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

type TModuleCommon struct {
	host       string
	sw         sync.WaitGroup
	ctx        context.Context
	cancle     context.CancelFunc
	data       []byte
	module     string
	clientPack tcpNet.IMessagePack
}

var exitchan = make(chan os.Signal, 1)

func NewModule(host, module string) *TModuleCommon {
	return &TModuleCommon{
		host:       host,
		module:     module,
		clientPack: &tcpNet.ClientProtocol{},
		data:       make([]byte, 1024),
	}
}

func (this *TModuleCommon) PushMsg(dstpoint, mainid, subid uint16, msg proto.Message) {
	buff := this.clientPack.PackMsg(dstpoint,
		mainid,
		subid,
		msg)
	this.data = make([]byte, len(buff))
	copy(this.data, buff)
	Log.FmtPrintln("msg len: ", len(this.data))
}

func (this *TModuleCommon) Run() {
	this.dialSend()
}

func (this *TModuleCommon) RunEx() {
	this.sendDirectNoRecv()
}

//发送信息
func (this *TModuleCommon) sender(conn net.Conn) (succ bool) {
	if len(this.data) == 0 {
		succ = true
		return
	}
	n, err := conn.Write(this.data)
	if n == 0 || err != nil {
		Log.Error("Write fail, data: ", n, err)
		return false
	}
	Log.FmtPrintln("send over")
	succ = true
	return
}

func (this *TModuleCommon) readloop(conn net.Conn) {
	for {
		select {
		case <-this.ctx.Done():
			this.sw.Done()
			return
		default:
			//接收服务端反馈
			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil || n == 0 {
				Log.Error("waiting server back msg error: ", conn.RemoteAddr().String(), err)
				continue
			}

			_, err = this.clientPack.UnPackAction(buffer)
			if err != nil {
				Log.Error("unpack action err: ", err)
				return
			}

			route := this.clientPack.GetRouteID()
			Log.FmtPrintln("pack route: ", route)
			mainID, subID := this.clientPack.GetMessageID()
			Log.FmtPrintf("mainid: %v, subID: %v.", mainID, subID)
			Log.FmtPrintf("receive server back, ip: %v.", conn.RemoteAddr().String())
		}
	}

}

func (this *TModuleCommon) sendloop(conn net.Conn) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", this.host)
	if err != nil {
		Log.FmtPrintf("Fatal error: %s", err.Error())
		os.Exit(1)
	}

	for i := 0; i < 1; i++ {
		Log.FmtPrintln("time: ", i)
		if !this.sender(conn) {
			tick := time.NewTicker(time.Duration(3 * time.Second))
			for {
				select {
				case <-tick.C:
					conn, err = net.DialTCP("tcp", nil, tcpAddr)
					if err != nil {
						Log.FmtPrintf("dial to server, host: %v.", this.host)
						Log.Error("err: ", err.Error())
						continue
					}
					break
				default:

				}

			}
		}
		time.Sleep(time.Duration(2) * time.Second)
	}
}

func (this *TModuleCommon) dialSend() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", this.host)
	if err != nil {
		Log.Error("resolve error: %s", err.Error())
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		Log.Error("dial error: %s", err.Error())
		return
	}

	this.ctx, this.cancle = context.WithCancel(context.Background())
	Log.FmtPrintln("connection success")
	signal.Notify(exitchan, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGSEGV)

	this.sw.Add(3)
	go this.readloop(conn)
	go this.sendloop(conn)
	go this.exitloop()
	this.sw.Wait()
}

func (this *TModuleCommon) sendDirectNoRecv() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", this.host)
	if err != nil {
		Log.Error("resolve error: %s", err.Error())
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		Log.Error("dial error: %s", err.Error())
		return
	}

	this.ctx, this.cancle = context.WithCancel(context.Background())
	Log.FmtPrintln("connection success")
	signal.Notify(exitchan, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGSEGV)

	this.sendloop(conn)
}

func (this *TModuleCommon) exitloop() {
	for {
		//Block until a signal is received.
		if s, ok := <-exitchan; ok {
			fmt.Println("Got signal:", s)
		}
		os.Exit(1)
		select {
		case <-this.ctx.Done():
			this.sw.Done()
			return
		default:

		}
	}
}
