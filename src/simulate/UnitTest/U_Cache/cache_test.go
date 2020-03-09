package U_Cache

import (
	"common/Cache"
	"common/Log"
	"fmt"
	"golang.org/x/sync/singleflight"
	"sync"
	"sync/atomic"
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

func TestSingleLight(t *testing.T) {
	var g singleflight.Group
	v, err, _ := g.Do("key", func() (interface{}, error) {
		return "bar", nil
	})

	if err != nil {
		Log.Error("Do error = %v", err)
		return
	}

	got := fmt.Sprintf("%v (%T)", v, v)
	want := "bar (string)"
	Log.FmtPrintf("Do = %v; want %v", got, want)
}

func TestDoDupSuppress(t *testing.T) {
	var g singleflight.Group
	var wg1, wg2 sync.WaitGroup
	c := make(chan string, 1)
	var calls int32
	fn := func() (interface{}, error) {
		Log.FmtPrintf("fn calls = %T %v.", calls, calls)
		if atomic.AddInt32(&calls, 1) == 1 {
			// First invocation.
			Log.FmtPrintf("wg1.Done in fn.")
			wg1.Done()
		}
		v := <-c
		c <- v // pump; make available for any future calls
		Log.FmtPrintf("fn v = %T %v.", v, v)
		time.Sleep(10 * time.Millisecond) // let more goroutines enter Do

		return v, nil
	}

	const n = 10
	wg1.Add(1)
	for i := 0; i < n; i++ {
		wg1.Add(1)
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			wg1.Done()
			Log.FmtPrintf("wg1.Done before Do.")
			v, err, _ := g.Do("key", fn)
			if err != nil {
				Log.Error("Do error: %v", err)
				return
			}

			Log.FmtPrintf("Do = %T %v; want %q", v, v, "bar")
		}()
	}
	wg1.Wait()
	// At least one goroutine is in fn now and all of them have at
	// least reached the line before the Do.
	Log.FmtPrintf("begin bar -> c.")
	c <- "bar"
	// wg1.Wait()
	// Log.FmtPrintf("begin bar -> c.")
	// c <- "bar"
	wg2.Wait()
	Log.FmtPrintf("calls = %T %v.", calls, calls)
	if got := atomic.LoadInt32(&calls); got <= 0 || got >= n {
		t.Errorf("number of calls = %d; want over 0 and less than %d", got, n)
	}
}
