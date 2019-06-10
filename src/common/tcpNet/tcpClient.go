package tcpNet
// client connect server.
import (
	"net"
	"sync"
	//"time"
	"context"
	"fmt"
)

type TcpClient struct {
	wg  	sync.WaitGroup
	ctx 	context.Context
	cancel	context.CancelFunc
	host    string
	s 		*TcpSession
	srcSvr  int32
	dstSvr  int32
	sessAlive bool
	cb 		MessageCb
}

func NewClient(host string, srcSvr, dstSvr int32, cb MessageCb)*TcpClient{
	return &TcpClient{
		host: 	host,
		srcSvr:	srcSvr,
		dstSvr: dstSvr,
		cb: 	cb,
	}
}

func (self *TcpClient) Run(){
	self.ctx, self.cancel = context.WithCancel(context.Background())
	self.connect()
	self.wg.Add(1)
	go self.loopconn()
}

func (self *TcpClient) connect(){
	c, err := net.Dial("tcp", self.host)
	if err != nil {
		fmt.Println("net dial err: ", err)
		return
	}
	c.(*net.TCPConn).SetNoDelay(true)
	self.s = NewSession(self.host, c, self.ctx, self.srcSvr, self.dstSvr, self.cb)
	self.s.HandleSession()
}

func (self *TcpClient) loopconn(){
	for {
		select{
		case <-self.ctx.Done():
			self.Exit()
			return
		default:
			if !self.sessAlive {
				self.connect()
			}
		}
	}
}

func (self *TcpClient) Send(data []byte){
	self.s.SetSendCache(data)
}

func (self *TcpClient) Exit(){
	self.sessAlive = false
	self.cancel()
	self.s.exit()
	self.wg.Wait()
}