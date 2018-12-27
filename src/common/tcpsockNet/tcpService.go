/*
* CopyRight(C) StefanChen e-mail:2572915286@qq.com
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

package tcpsockNet

import(
	"net"
	"os"
	"fmt"
	"sync"
)

type TcpService struct{
	sw  	sync.WaitGroup
	host   	string
	listener *net.TCPListener
}

func NewTcpServer(addr string)*TcpService{
	return &TcpService{
		host: addr,
	}
}

func (self *TcpService) StartTcpServer(){
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	self.listener = listener

	self.sw.Add(1)
	go self.acceptLoop()
	self.sw.Wait()
}

func (self *TcpService) acceptLoop(){
	defer self.sw.Done()
	for{
		c, err := self.listener.AcceptTCP()
		if err != nil {
			fmt.Errorf("[TcpService][acceptLoop] can not accept tcp .")
		}

		session := NewSession(self.host, c)
		self.handleSession(session)
	}
}

func (self *TcpService) handleSession(s *TcpSession){
	self.sw.Add(1)
	go s.Recvmessage(&self.sw)
	self.sw.Add(1)
	go s.Sendmessage(&self.sw)
	self.sw.Wait()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}