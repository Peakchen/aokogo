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
	cb       MessageCb
	// person offline flag
	off chan *TcpSession
	// person online
	person     int32
	SvrType    Define.ERouteId
	Adacb      AfterDialAct
	mpobj      IMessagePack
	SessionMgr IProcessConnSession
}

func NewClient(host string, SvrType Define.ERouteId, cb MessageCb, Ada AfterDialAct, sessionMgr IProcessConnSession) *TcpClient {
	return &TcpClient{
		host:       host,
		cb:         cb,
		SvrType:    SvrType,
		Adacb:      Ada,
		SessionMgr: sessionMgr,
	}
}

func (this *TcpClient) Run() {
	this.ctx, this.cancel = context.WithCancel(context.Background())
	this.mpobj = &ClientProtocol{}
	this.connect(&this.wg)

	this.wg.Add(2)
	go this.loopconn(&this.wg)
	go this.loopoff(&this.wg)
	this.wg.Wait()
}

func (this *TcpClient) connect(sw *sync.WaitGroup) error {
	if this.dialsess != nil {
		if this.dialsess.isAlive {
			return nil
		}
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", this.host)
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
	this.dialsess = NewSession(this.host, c, this.ctx, this.SvrType, this.cb, this.off, this.mpobj, this)
	this.dialsess.HandleSession(sw)
	//this.AddSession(this.dialsess)
	this.afterDial()
	return nil
}

func (this *TcpClient) loopconn(sw *sync.WaitGroup) {
	defer sw.Done()
	defer this.Exit(sw)

	conntick := time.NewTicker(time.Duration(3 * time.Second))
	for {
		select {
		case <-this.ctx.Done():
			return
		case <-conntick.C:
			if err := this.connect(sw); err != nil {
				Log.FmtPrintf("dail to server fail, host: %v.", this.host)
				continue
			}
		default:

		}
	}
}

func (this *TcpClient) loopoff(sw *sync.WaitGroup) {
	defer this.Exit(sw)
	for {
		select {
		case os, ok := <-this.off:
			if !ok {
				return
			}
			this.offline(os)
		case <-this.ctx.Done():
			return
		}
	}
}

func (this *TcpClient) offline(os *TcpSession) {
	// process

}

func (this *TcpClient) Send(data []byte) {
	this.dialsess.SetSendCache(data)
}

func (this *TcpClient) SendMessage() {

}

func (this *TcpClient) PushCmdSession(session *TcpSession, cmds []uint32) {
	if this.SessionMgr == nil {
		return
	}
	//this.SessionMgr.AddSessionByCmd(session, cmds)
	this.AddSession(session)
}

func (this *TcpClient) GetSessionByCmd(cmd uint32) (session *TcpSession) {
	if this.SessionMgr == nil {
		return
	}
	return this.SessionMgr.GetByCmd(cmd)
}

func (this *TcpClient) AddSession(session *TcpSession) {
	if this.SessionMgr == nil {
		return
	}
	this.SessionMgr.AddSession(session)
}

func (this *TcpClient) GetSessionByID(sessionID uint64) (session *TcpSession) {
	return this.SessionMgr.GetSessionByID(sessionID)
}

func (this *TcpClient) Exit(sw *sync.WaitGroup) {
	this.dialsess = nil
	this.cancel()
	sw.Wait()
}

func (this *TcpClient) sendRegisterMsg() {
	Log.FmtPrintf("after dial, send point: %v register message to server.", this.SvrType)
	req := &MSG_Server.CS_ServerRegister_Req{}
	req.ServerType = int32(this.SvrType)
	req.Msgs = GetAllMessageIDs()
	Log.FmtPrintln("register context: ", req.Msgs)
	buff := this.mpobj.PackMsg(uint16(this.SvrType),
		uint16(MSG_MainModule.MAINMSG_SERVER),
		uint16(MSG_Server.SUBMSG_CS_ServerRegister),
		req)
	this.Send(buff)
}

func (this *TcpClient) afterDial() {
	if this.Adacb != nil {
		this.Adacb(this.dialsess)
	}
	this.sendRegisterMsg()
}

func (this *TcpClient) SessionType() (st ESessionType) {
	return ESessionType_Client
}

func (this *TcpClient) GetSessionByType(svrType Define.ERouteId) (session *TcpSession) {
	if this.SessionMgr == nil {
		return
	}
	return this.SessionMgr.GetSessionByType(svrType)
}

func (this *TcpClient) RemoveSession(session *TcpSession) {
	if this.SessionMgr == nil {
		return
	}

	if session.RegPoint != Define.ERouteId_ER_Invalid {
		this.SessionMgr.RemoveSessionByType(session.RegPoint)
	} else {
		this.SessionMgr.RemoveSessionByID(session)
	}
}
