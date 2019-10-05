package tcpNet

// client connect server.
import (
	"common/Define"
	"common/Log"
	"common/msgProto/MSG_MainModule"
	"common/msgProto/MSG_Server"
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
	person  int32
	SvrType Define.ERouteId
	Adacb   AfterDialAct
	mpobj   IMessagePack
}

func NewClient(host string, SvrType Define.ERouteId, mapSvr *map[int32][]int32, cb MessageCb, Ada AfterDialAct) *TcpClient {
	return &TcpClient{
		host:    host,
		mapSvr:  *mapSvr,
		cb:      cb,
		SvrType: SvrType,
		Adacb:   Ada,
	}
}

func (self *TcpClient) Run() {
	self.ctx, self.cancel = context.WithCancel(context.Background())
	sw := &sync.WaitGroup{}
	self.mpobj = &ClientProtocol{}
	self.connect(sw)
	self.wg.Add(2)
	go self.loopconn(sw)
	go self.loopoff(sw)
	self.wg.Wait()
}

func (self *TcpClient) connect(sw *sync.WaitGroup) error {
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

	self.sessAlive = true
	c.SetNoDelay(true)
	self.s = NewSession(self.host, c, self.ctx, &self.mapSvr, self.cb, self.off, self.mpobj)
	self.s.HandleSession(sw)
	if self.Adacb != nil {
		self.Adacb()
	}
	self.sendRegisterMsg()
	return nil
}

func (self *TcpClient) loopconn(sw *sync.WaitGroup) {
	defer self.Exit(sw)
	for {
		select {
		case <-self.ctx.Done():
			return
		default:
			if self.sessAlive {
				continue
			}
			if err := self.connect(sw); err != nil {
				Log.FmtPrintf("dail to server fail, host: ", self.host)
				return
			}
		}
	}
}

func (self *TcpClient) loopoff(sw *sync.WaitGroup) {
	defer self.Exit(sw)
	for {
		select {
		case os, ok := <-self.off:
			if !ok {
				return
			}
			self.offline(os)
		case <-self.ctx.Done():
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

func (self *TcpClient) Exit(sw *sync.WaitGroup) {
	self.sessAlive = false
	self.cancel()
	self.s.exit(sw)
	sw.Wait()
}

func (self *TcpClient) sendRegisterMsg() {
	Log.FmtPrintf("after dial, send point: %v register message to server.", self.SvrType)
	req := &MSG_Server.CS_ServerRegister_Req{}
	req.ServerType = int32(self.SvrType)
	req.Msgs = GetAllMessageIDs()
	Log.FmtPrintln(req.Msgs)
	buff := self.mpobj.GetSendPackMsg(uint16(MSG_MainModule.MAINMSG_SERVER),
		uint16(MSG_Server.SUBMSG_CS_ServerRegister),
		req)
	self.Send(buff)
}
