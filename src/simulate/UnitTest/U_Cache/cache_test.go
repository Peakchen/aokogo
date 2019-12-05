package U_Cache

import (
	"common/Cache"
	"common/Log"
	"testing"
	"time"
)

var (
	ck  string = "1"
	ckd string = "aaa"
)

func TestCacheNormal(t *testing.T) {
	Cache.Init()

	Cache.SetTempData(ck, ckd)
	Log.FmtPrintln("[Normal] cache temp data: ", Cache.GetTempData(ck))
}

func TestCacheDealLine(t *testing.T) {
	TestCacheNormal(t)
	time.Sleep(time.Duration(Cache.ConstCacheOverTime) * time.Second)
	Log.FmtPrintln("[DealLine] cache temp data: ", Cache.GetTempData(ck))
}
