package U_Log

import (
	"common/Log"
	"sync"
	"testing"
)

func TestLogNormal(t *testing.T) {
	Log.FmtPrintln("test log: ", "yes")
	Log.Error("test error.")
}

func TestLogLoopWrite(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i := 0; i < 100; i++ {
			Log.Info("info, idx: %v.", i)
		}
	}()

	//time.Sleep(time.Duration(30 * time.Second))
	wg.Wait()
}
