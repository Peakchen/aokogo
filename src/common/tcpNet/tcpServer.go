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

import(
	"net"
	"os"
	"fmt"
	"sync"
	"context"
)

type TcpServer struct{
	sw  	sync.WaitGroup
	host   	string
	listener *net.TCPListener
	ctx 	context.Context
	cancel	context.CancelFunc
	srcSvr  int32
	dstSvr  int32
	cb 		MessageCb
}

func NewTcpServer(addr string, srcSvr, dstSvr int32, cb MessageCb)*TcpServer{
	return &TcpServer{
		host: 	addr,
		srcSvr: srcSvr,
		dstSvr: dstSvr,
		cb:		cb,
	}
}

func (self *TcpServer) StartTcpServer(){
	tcpAddr, err := net.ResolveTCPAddr("tcp4", self.host)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	self.listener = listener
	self.ctx, self.cancel = context.WithCancel(context.Background())

	self.sw.Add(1)
	go self.loop()
	self.sw.Wait()
}

func (self *TcpServer) loop(){
	defer self.sw.Done()
	for{
		select {
		case <-self.ctx.Done():
			self.Exit()
			return
		default:
			c, err := self.listener.AcceptTCP()
			if err != nil {
				fmt.Println("[TcpServer][acceptLoop] can not accept tcp .")
			}
	
			session := NewSession(self.host, c, self.ctx, self.srcSvr, self.dstSvr, self.cb)
			session.HandleSession()
		}
	}
}

func (self *TcpServer) Exit(){
	self.listener.Close()
	self.cancel()
	self.sw.Wait()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}