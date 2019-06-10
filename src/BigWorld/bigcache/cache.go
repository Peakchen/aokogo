package bigcache

import (
	. "common/cache"
	"context"
	//"strings"
)

type TBigCache struct {
	c 		*TCache
	ctx 	context.Context
	cancel 	context.CancelFunc
}

func (self *TBigCache) StartCache(){
	self.ctx, self.cancel = context.WithCancel(context.Background())
	self.c = &TCache{}
	self.c.Init(ConstCacheOverTime, self.ctx)
	self.c.Run()
}

func (self *TBigCache) Exit(){
	self.cancel()
}

func (self *TBigCache) Set(key, name, data string){
	cid := key+":"+name
	self.c.Set(cid, data)
}

func (self *TBigCache) Get(key, name string) interface{}{
	cid := key+":"+name
	return self.c.Get(cid)
}