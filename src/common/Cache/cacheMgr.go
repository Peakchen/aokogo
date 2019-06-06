package Cache

/*
	cache data not exists forever, 
	which exist over deadline, 
	then delete it and to get new data form bigword.
*/

import (
	"sync"
	"time"
	"context"
	"net"
	"fmt"
	. "common/tcpNet"
	."common/S2SMessage"
	"github.com/golang/protobuf/proto"
)


type TCacheMgr struct {
	wg  	sync.WaitGroup
	c 		*TCache //some one cache
	ctx 	context.Context
	cancel	context.CancelFunc
	s       *TcpSession
	srcSvr  int32
	dstSvr  int32
	// add cache obj ...
}

func (self *TCacheMgr) connect(){
	c, err := net.Dial("tcp", ConstBigWordHost)
	if err != nil {
		fmt.Println("net dial err: ", err)
		return
	}
	c.(*net.TCPConn).SetNoDelay(true)
	self.s = NewSession(ConstBigWordHost, c, self.ctx, self.src, )
	self.s.HandleSession()
}

func (self *TCacheMgr) run(srcSev, dstSvr int32){
	self.ctx, self.cancel = context.WithCancel(context.Background())
	self.srcSvr = srcSev
	self.dstSvr = dstSvr
	self.wg.Add(1)
	self.c = &TCache{}
	self.c.init(ConstCacheOverTime, self.ctx)
	self.c.run()
	// add ...
}

func (self *TCacheMgr) exit(){
	self.cancel()
	self.wg.Wait()
}

func (self *TCacheMgr)loopc(){
	defer self.wg.Done()
	t := time.NewTicker(time.Duration(ConstCacheUpdateTime))
	for {
		select {
		case <-self.ctx.Done():
			self.exit()
			return
		case <-t.C:
			//begin request data from bigword server.
			s := &S2SBaseMessage{
				Srcid: self.srcSvr,
				Dstid: int32(ServerId_SID_BigWorld),
			}
			msg, err := proto.Marshal(s)
			if err == nil {
				fmt.Println("Marshal message fail.")
				return
			}
			self.s.SetSendCache(msg)
		}
	}
}