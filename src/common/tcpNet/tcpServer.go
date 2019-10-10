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
	"context"
	"fmt"
	"net"
	"os"
	"sync"
)

type TcpServer struct {
	sw       *sync.WaitGroup
	host     string
	listener *net.TCPListener
	ctx      context.Context
	cancel   context.CancelFunc
	mapSvr   map[int32][]int32
	cb       MessageCb
	off      chan *TcpSession
	on       *TcpSession
	// person online
	person     int32
	SvrType    Define.ERouteId
	pack       IMessagePack
	sessionMgr TMessageSession
}

func NewTcpServer(addr string, SvrType Define.ERouteId, mapSvr *map[int32][]int32, cb MessageCb, sessionMgr TMessageSession) *TcpServer {
	return &TcpServer{
		host:       addr,
		mapSvr:     *mapSvr,
		cb:         cb,
		SvrType:    SvrType,
		sessionMgr: sessionMgr,
	}
}

func (self *TcpServer) StartTcpServer(sw *sync.WaitGroup, ctx context.Context, cancle context.CancelFunc) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", self.host)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	self.listener = listener
	self.ctx, self.cancel = ctx, cancle
	self.pack = &ServerProtocol{}
	sw.Add(2)
	go self.loop(sw)
	go self.loopoff(sw)
	sw.Wait()
}

func (self *TcpServer) loop(sw *sync.WaitGroup) {
	defer self.Exit(sw)
	for {
		select {
		case <-self.ctx.Done():
			return
		default:
			c, err := self.listener.AcceptTCP()
			if err != nil || c == nil {
				Log.FmtPrintf("can not accept tcp.")
				continue
			}

			Log.FmtPrintf("connect here addr: %v.", c.RemoteAddr())
			c.SetNoDelay(true)
			c.SetKeepAlive(true)
			self.on = NewSession(self.host, c, self.ctx, &self.mapSvr, self.cb, self.off, self.pack, self.sessionMgr)
			self.on.HandleSession(sw)
			self.online()
		}
	}
}

func (self *TcpServer) loopoff(sw *sync.WaitGroup) {
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

func (self *TcpServer) online() {
	self.person++
}

func (self *TcpServer) offline(os *TcpSession) {
	// process
	self.person--
}

func (self *TcpServer) SendMessage() {

}

func (self *TcpServer) Exit(sw *sync.WaitGroup) {
	self.listener.Close()
	self.cancel()
	sw.Wait()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
