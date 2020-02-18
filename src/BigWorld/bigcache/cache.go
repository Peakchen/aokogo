package bigcache

import (
	. "common/cache"
	"context"
	//"strings"
)

type TBigCache struct {
	c      *TCache
	ctx    context.Context
	cancel context.CancelFunc
}

func (this *TBigCache) StartCache() {
	this.ctx, this.cancel = context.WithCancel(context.Background())
	this.c = &TCache{}
	this.c.Init(ConstCacheOverTime, this.ctx)
	this.c.Run()
}

func (this *TBigCache) Exit() {
	this.cancel()
}

func (this *TBigCache) Set(key, name string, data interface{}) {
	cid := key + ":" + name
	this.c.Set(cid, data)
}

func (this *TBigCache) Get(key, name string) interface{} {
	cid := key + ":" + name
	return this.c.Get(cid)
}

func (this *TBigCache) Del(key, name string) {
	cid := key + ":" + name
	this.c.Remove(cid)
}
