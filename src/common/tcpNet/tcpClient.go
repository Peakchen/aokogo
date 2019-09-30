package tcpNet

// client connect server.
import (
	"common/Log"
	"net"
	"os"
	"sync"

	//"time"
	"context"
)

type TcpClient struct {
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	host      string
	s         *TcpSession
	mapSvr    map[int32][]int32
	sessAlive bool
	cb        MessageCb
	// person offline flag
	off chan *TcpSession
	// person online
	person int32
}

func NewClient(host string, mapSvr *map[int32][]int32, cb MessageCb) *TcpClient {
	return &TcpClient{
		host:   host,
		mapSvr: *mapSvr,
		cb:     cb,
	}
}

func (self *TcpClient) Run() {
	self.ctx, self.cancel = context.WithCancel(context.Background())
	self.connect()
	self.wg.Add(1)
	go self.loopconn()
	self.wg.Wait()
}

func (self *TcpClient) connect() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", self.host)
	if err != nil {
		Log.FmtPrintf("Fatal error: %s", err.Error())
		os.Exit(1)
	}

	c, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		Log.FmtPrintf("net dial err: ", err)
		return err
	}

	c.SetNoDelay(true)
	self.s = NewSession(self.host, c, self.ctx, &self.mapSvr, self.cb, self.off, &ClientProtocol{})
	self.s.HandleSession()
	return nil
}

func (self *TcpClient) loopconn() {
	for {
		select {
		case <-self.ctx.Done():
			self.Exit()
			return
		default:
			if self.sessAlive {
				continue
			}
			if err := self.connect(); err != nil {
				Log.FmtPrintf("dail to server fail, host: ", self.host)
			}
		}
	}
}

func (self *TcpClient) loopoff() {
	defer self.wg.Done()
	for {
		select {
		case os, ok := <-self.off:
			if !ok {
				return
			}
			self.offline(os)
		case <-self.ctx.Done():
			self.Exit()
			return
		}
	}
}

func (self *TcpClient) offline(os *TcpSession) {
	// process

}

func (self *TcpClient) Send(data []byte) {
	self.s.SetSendCache(data)
}

func (self *TcpClient) SendMessage() {

}

func (self *TcpClient) Exit() {
	self.sessAlive = false
	self.cancel()
	self.s.exit()
	self.wg.Wait()
}
