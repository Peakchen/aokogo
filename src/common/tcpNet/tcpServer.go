/*
Copyright (c) <year> <copyright holders>

"Anti 996" License Version 1.0 (Draft)

Permission is hereby granted to any individual or legal entity
obtaining a copy of this licensed work (including the source code,
documentation and/or related items, hereinafter collectively referred
to as the "licensed work"), free of charge, to deal with the licensed
work for any purpose, including without limitation, the rights to use,
reproduce, modify, prepare derivative works of, distribute, publish
and sublicense the licensed work, subject to the following conditions:

1. The individual or the legal entity must conspicuously display,
without modification, this License and the notice on each redistributed
or derivative copy of the Licensed Work.

2. The individual or the legal entity must strictly comply with all
applicable laws, regulations, rules and standards of the jurisdiction
relating to labor and employment where the individual is physically
located or where the individual was born or naturalized; or where the
legal entity is registered or is operating (whichever is stricter). In
case that the jurisdiction has no such laws, regulations, rules and
standards or its laws, regulations, rules and standards are
unenforceable, the individual or the legal entity are required to
comply with Core International Labor Standards.

3. The individual or the legal entity shall not induce, metaphor or force
its employee(s), whether full-time or part-time, or its independent
contractor(s), in any methods, to agree in oral or written form, to
directly or indirectly restrict, weaken or relinquish his or her
rights or remedies under such laws, regulations, rules and standards
relating to labor and employment as mentioned above, no matter whether
such written or oral agreement are enforceable under the laws of the
said jurisdiction, nor shall such individual or the legal entity
limit, in any methods, the rights of its employee(s) or independent
contractor(s) from reporting or complaining to the copyright holder or
relevant authorities monitoring the compliance of the license about
its violation(s) of the said license.

THE LICENSED WORK IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE COPYRIGHT HOLDER BE LIABLE FOR ANY CLAIM,
DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
OTHERWISE, ARISING FROM, OUT OF OR IN ANY WAY CONNECTION WITH THE
LICENSED WORK OR THE USE OR OTHER DEALINGS IN THE LICENSED WORK.
*/

package tcpNet

import (
	"common/Define"
	"common/Log"
	"common/pprof"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
)

type TcpServer struct {
	sync.Mutex

	host      string
	pprofAddr string
	listener  *net.TCPListener
	cancel    context.CancelFunc
	cb        MessageCb
	off       chan *SvrTcpSession
	session   *SvrTcpSession
	// person online
	person     int32
	SvrType    Define.ERouteId
	pack       IMessagePack
	SessionMgr IProcessConnSession
	// session id
	SessionID uint64
}

func NewTcpServer(listenAddr, pprofAddr string, SvrType Define.ERouteId, cb MessageCb, sessionMgr IProcessConnSession) *TcpServer {
	return &TcpServer{
		host:       listenAddr,
		pprofAddr:  pprofAddr,
		cb:         cb,
		SvrType:    SvrType,
		SessionMgr: sessionMgr,
		SessionID:  ESessionBeginNum,
	}
}

func (this *TcpServer) Run() {
	os.Setenv("GOTRACEBACK", "crash")
	tcpAddr, err := net.ResolveTCPAddr("tcp4", this.host)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	this.listener = listener

	var (
		ctx context.Context
		sw  = sync.WaitGroup{}
	)

	ctx, this.cancel = context.WithCancel(context.Background())
	pprof.Run(ctx)

	this.pack = &ServerProtocol{}
	sw.Add(3)
	go this.loop(ctx, &sw)
	go this.loopoff(ctx, &sw)
	go func() {
		Log.FmtPrintln("[server] run http server, host: ", this.pprofAddr)
		http.ListenAndServe(this.pprofAddr, nil)
	}()
	sw.Wait()
}

func (this *TcpServer) loop(ctx context.Context, sw *sync.WaitGroup) {
	defer func() {
		sw.Done()
		this.Exit(sw)
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			c, err := this.listener.AcceptTCP()
			if err != nil || c == nil {
				Log.FmtPrintf("can not accept tcp.")
				continue
			}

			c.SetNoDelay(true)
			c.SetKeepAlive(true)
			atomic.AddUint64(&this.SessionID, 1)
			Log.FmtPrintf("[server] accept connect here addr: %v, SessionID: %v.", c.RemoteAddr(), this.SessionID)
			this.session = NewSvrSession(c.RemoteAddr().String(), c, ctx, this.SvrType, this.cb, this.off, this.pack)
			this.session.HandleSession(sw)
			this.online()
		}
	}
}

func (this *TcpServer) loopoff(ctx context.Context, sw *sync.WaitGroup) {
	defer func() {
		sw.Done()
		this.Exit(sw)
	}()
	for {
		select {
		case os, ok := <-this.off:
			if !ok {
				return
			}
			this.offline(os)
		case <-ctx.Done():
			return
		}
	}
}

func (this *TcpServer) online() {
	this.person++
}

func (this *TcpServer) offline(os *SvrTcpSession) {
	// process
	this.person--
}

func (this *TcpServer) SendMessage() {

}

func (this *TcpServer) Exit(sw *sync.WaitGroup) {
	this.listener.Close()
	this.cancel()
	pprof.Exit()
	sw.Wait()
}

func (this *TcpServer) SessionType() (st ESessionType) {
	return ESessionType_Server
}

func (this *TcpServer) RemoveSession(session *SvrTcpSession) {
	if this.SessionMgr == nil {
		return
	}

	if session.RegPoint != Define.ERouteId_ER_Invalid {
		this.SessionMgr.RemoveSession(session.RemoteAddr)
	} else {
		this.SessionMgr.RemoveSession(session.StrIdentify)
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
