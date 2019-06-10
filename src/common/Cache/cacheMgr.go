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
	. "common/Define"
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
	sessAlive bool
	cb 		MessageCb
	// add cache obj ...
}

func (self *TCacheMgr) connect()error{
	c, err := net.Dial("tcp", ConstBigWordHost)
	if err != nil {
		fmt.Println("net dial err: ", err)
		return err
	}
	c.(*net.TCPConn).SetNoDelay(true)
	self.s = NewSession(ConstBigWordHost, c, self.ctx, self.srcSvr, self.dstSvr, self.cb)
	self.s.HandleSession()
	self.sessAlive = true
	return nil
}

func (self *TCacheMgr) run(srcSev, dstSvr int32, cb MessageCb){
	self.ctx, self.cancel = context.WithCancel(context.Background())
	self.srcSvr = srcSev
	self.dstSvr = dstSvr
	self.cb = cb
	
	self.c = &TCache{}
	self.c.Init(ConstCacheOverTime, self.ctx)
	self.c.Run()
	self.connect()

	self.wg.Add(1)
	go self.loopc()
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
			msg := &SS_BaseMessage_Req{
				Srcid: self.srcSvr,
				Dstid: int32(ERouteId_ER_BigWorld),
			}
			data, err := proto.Marshal(msg)
			if err == nil {
				fmt.Println("Marshal message fail.")
				return
			}
			self.s.SetSendCache(data)
		default:
			if !self.sessAlive {
				self.connect()
			}
		}
	}
}