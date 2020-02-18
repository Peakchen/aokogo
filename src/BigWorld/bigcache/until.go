package bigcache

// add by stefan 20190610 20:10
import ()

var GBigCache *TBigCache = &TBigCache{}

func BigCacheSet(key, name string, data interface{}) {
	GBigCache.Set(key, name, data)
}

func BigCacheGet(key, name string) interface{} {
	return GBigCache.Get(key, name)
}

func BigCacheRemove(key, name string) {
	GBigCache.Del(key, name)
}
