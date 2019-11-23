package tcpNet

// client connect server.
import (
	"common/Define"
	"common/Log"
	"common/msgProto/MSG_MainModule"
	"common/msgProto/MSG_Server"
	"common/pprof"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	//"time"
	"context"
)

type TcpClient struct {
	sync.Mutex

	cancel    context.CancelFunc
	host      string
	pprofAddr string
	dialsess  *TcpSession
	cb        MessageCb
	// person offline flag
	off chan *TcpSession
	// person online
	person     int32
	SvrType    Define.ERouteId
	Adacb      AfterDialAct
	mpobj      IMessagePack
	SessionMgr IProcessConnSession
}

func NewClient(host, pprofAddr string, SvrType Define.ERouteId, cb MessageCb, Ada AfterDialAct, sessionMgr IProcessConnSession) *TcpClient {
	return &TcpClient{
		host:       host,
		pprofAddr:  pprofAddr,
		cb:         cb,
		SvrType:    SvrType,
		Adacb:      Ada,
		SessionMgr: sessionMgr,
	}
}

func (this *TcpClient) Run() {
	os.Setenv("GOTRACEBACK", "crash")

	var (
		ctx context.Context
		sw  = sync.WaitGroup{}
	)

	ctx, this.cancel = context.WithCancel(context.Background())
	this.mpobj = &ClientProtocol{}
	this.connect(ctx, &sw)
	pprof.Run(ctx)
	sw.Add(3)
	go this.loopconn(ctx, &sw)
	go this.loopoff(ctx, &sw)
	go func() {
		Log.FmtPrintln("[client] run http server, host: ", this.pprofAddr)
		http.ListenAndServe(this.pprofAddr, nil)
	}()
	sw.Wait()
}

func (this *TcpClient) connect(ctx context.Context, sw *sync.WaitGroup) (err error) {
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

	Log.FmtPrintf("[----------client-----------], svrtype: %v.", this.SvrType)
	c.SetNoDelay(true)
	this.dialsess = NewSession(this.host, c, ctx, this.SvrType, this.cb, this.off, this.mpobj, this)
	this.dialsess.HandleSession(sw)
	this.afterDial()
	return nil
}

func (this *TcpClient) loopconn(ctx context.Context, sw *sync.WaitGroup) {
	defer func() {
		sw.Done()
		this.Exit(sw)
	}()

	ticker := time.NewTicker(time.Duration(EClientSessionCheckInterval) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if this.dialsess == nil || false == this.dialsess.isAlive {
				if err := this.connect(ctx, sw); err != nil {
					Log.FmtPrintf("dail to server fail, host: %v.", this.host)
				}
			}
			//default:
			// if this.dialsess == nil || false == this.dialsess.isAlive {
			// 	if err := this.connect(sw); err != nil {
			// 		Log.FmtPrintf("dail to server fail, host: %v.", this.host)
			// 	}
			// }
		}
	}
}

func (this *TcpClient) loopoff(ctx context.Context, sw *sync.WaitGroup) {
	defer func() {
		sw.Done()
		this.Exit(sw)
	}()

	for {
		select {
		case os, ok := <-this.off:
			if ok {
				this.offline(os)
			}

		case <-ctx.Done():
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
	Log.FmtPrintf("[client] push session, SvrType: %v, RegPoint: %v.", session.SvrType, session.RegPoint)
	if session.RegPoint == Define.ERouteId_ER_ESG {
		GServer2ServerSession.AddSession(session.RegPoint, session)
	} else {
		if this.SessionMgr == nil {
			return
		}
		this.SessionMgr.AddSession(session.RegPoint, session)
	}
}

func (this *TcpClient) Exit(sw *sync.WaitGroup) {
	this.dialsess = nil
	this.cancel()
	pprof.Exit()
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
