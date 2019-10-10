package tcpNet

// client connect server.
import (
	"common/Define"
	"common/Log"
	"common/msgProto/MSG_MainModule"
	"common/msgProto/MSG_Server"
	"net"
	"sync"
	"time"

	//"time"
	"context"
)

type TcpClient struct {
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	host     string
	dialsess *TcpSession
	mapSvr   map[int32][]int32
	cb       MessageCb
	// person offline flag
	off chan *TcpSession
	// person online
	person     int32
	SvrType    Define.ERouteId
	Adacb      AfterDialAct
	mpobj      IMessagePack
	sessionMgr TMessageSession
}

func NewClient(host string, SvrType Define.ERouteId, mapSvr *map[int32][]int32, cb MessageCb, Ada AfterDialAct, sessionMgr TMessageSession) *TcpClient {
	return &TcpClient{
		host:       host,
		mapSvr:     *mapSvr,
		cb:         cb,
		SvrType:    SvrType,
		Adacb:      Ada,
		sessionMgr: sessionMgr,
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
	if self.dialsess != nil {
		if self.dialsess.isAlive {
			return nil
		}
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", self.host)
	if err != nil {
		Log.Error("resolve tcp error: %v.", err.Error())
		return err
	}

	c, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		Log.Error("net dial err: %v.", err)
		return err
	}

	c.SetNoDelay(true)
	self.dialsess = NewSession(self.host, c, self.ctx, &self.mapSvr, self.cb, self.off, self.mpobj, self.sessionMgr)
	self.dialsess.HandleSession(sw)
	self.afterDial()
	return nil
}

func (self *TcpClient) loopconn(sw *sync.WaitGroup) {
	defer sw.Done()
	defer self.Exit(sw)

	conntick := time.NewTicker(time.Duration(3 * time.Second))
	for {
		select {
		case <-self.ctx.Done():
			return
		case <-conntick.C:
			if err := self.connect(sw); err != nil {
				Log.FmtPrintf("dail to server fail, host: %v.", self.host)
				continue
			}
		default:

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
	self.dialsess.SetSendCache(data)
}

func (self *TcpClient) SendMessage() {

}

func (self *TcpClient) Exit(sw *sync.WaitGroup) {
	self.dialsess = nil
	self.cancel()
	sw.Wait()
}

func (self *TcpClient) sendRegisterMsg() {
	Log.FmtPrintf("after dial, send point: %v register message to server.", self.SvrType)
	req := &MSG_Server.CS_ServerRegister_Req{}
	req.ServerType = int32(self.SvrType)
	req.Msgs = GetAllMessageIDs()
	Log.FmtPrintln("register context: ", req.Msgs)
	buff := self.mpobj.PackMsg(uint16(MSG_MainModule.MAINMSG_SERVER),
		uint16(MSG_Server.SUBMSG_CS_ServerRegister),
		req)
	self.Send(buff)
}

func (self *TcpClient) afterDial() {
	if self.Adacb != nil {
		self.Adacb()
	}
	self.sendRegisterMsg()
}
