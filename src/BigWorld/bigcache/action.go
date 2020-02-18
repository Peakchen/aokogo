package bigcache

import (
	. "common/S2SMessage"
	"log"
	//"reflect"
)

type CacheOperCB func(key string, c *CacheOperation, out interface{})

var CacheOperCall map[ECacheOper]CacheOperCB = map[ECacheOper]CacheOperCB{
	ECacheOper_Add:    AddBigC,
	ECacheOper_Delete: DeleteBigC,
	ECacheOper_Update: UpdateBigC,
	ECacheOper_Select: SelectBigC,
}

func SelectOper(key string, c *CacheOperation, out interface{}) {
	cb, ok := CacheOperCall[c.Oper]
	if !ok || cb == nil {
		log.Fatal("can not find cache cb, oper: ", c.Oper)
		return
	}
	cb(key, c, out)
}

func AddBigC(key string, c *CacheOperation, out interface{}) {
	UpdateBigC(key, c, out)
}

func DeleteBigC(key string, c *CacheOperation, out interface{}) {
	for _, cache := range c.Caches {
		cname, ok := GCacheTab[cache.NameId]
		if !ok {
			log.Fatal("can not find cache name, NameId: ", cache.NameId)
			return
		}
		BigCacheRemove(key, cname)
	}
}

func UpdateBigC(key string, c *CacheOperation, out interface{}) {
	for _, cache := range c.Caches {
		cname, ok := GCacheTab[cache.NameId]
		if !ok {
			log.Fatal("can not find cache name, NameId: ", cache.NameId)
		}
		BigCacheSet(key, cname, cache.Data)
	}
}

func SelectBigC(key string, c *CacheOperation, out interface{}) {
	querycolls := []interface{}{}
	for _, cache := range c.Caches {
		cname, ok := GCacheTab[cache.NameId]
		if !ok {
			log.Fatal("can not find cache name, NameId: ", cache.NameId)
			return
		}
		outc := &ArrayCache_Repeat{}
		outc.NameId = cache.NameId
		outc.Data = BigCacheGet(key, cname).(string)
		querycolls = append(querycolls, outc)
	}
	out = querycolls
}
