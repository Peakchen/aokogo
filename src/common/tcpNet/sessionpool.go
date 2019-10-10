package tcpNet

import "sync"

// client session pool

var (
	GClientSessionPool *TClientSession
	GServerSessionPool *TServerSession
)

func init() {
	GClientSessionPool = &TClientSession{
		sessionpool: &sync.Pool{
			New: func() interface{} {
				return new(TcpSession)
			},
		},
	}

	GServerSessionPool = &TServerSession{
		sessionpool: &sync.Pool{
			New: func() interface{} {
				return new(TcpSession)
			},
		},
	}
}

type TClientSession struct {
	sessionpool *sync.Pool
}

func (this *TClientSession) Push(s *TcpSession) {
	this.sessionpool.Put(s)
}

func (this *TClientSession) Get() (s *TcpSession) {
	s = this.sessionpool.Get().(*TcpSession)
	return
}

// server session pool

type TServerSession struct {
	sessionpool *sync.Pool
}

func (this *TServerSession) Push(s *TcpSession) {
	this.sessionpool.Put(s)
}

func (this *TServerSession) Get() (s *TcpSession) {
	s = this.sessionpool.Get().(*TcpSession)
	return
}
