package AsyncLock

import (
	"gopkg.in/redsync.v1"
	//"github.com/garyburd/redigo/redis"
	//"fmt"
	"time"
)

var GAsynclock map[string]*redsync.Mutex = map[string]*redsync.Mutex{}
var GRedsyncObj *redsync.Redsync

func NewAsyncLock(pools []redsync.Pool) {
	GRedsyncObj = redsync.New(pools)
}

func AddAsyncLock(key, Name string) {
	lockid := key + ":" + Name
	if _, ok := GAsynclock[lockid]; !ok {
		GAsynclock[lockid] = GRedsyncObj.NewMutex(lockid,
			redsync.SetExpiry(time.Duration(10*time.Second)),
			redsync.SetRetryDelay(time.Duration(1*time.Second)))
	}
	GAsynclock[lockid].Lock()
	return
}

func ReleaseAsyncLock(key, Name string) {
	lockid := key + ":" + Name
	if _, ok := GAsynclock[lockid]; !ok {
		return
	}
	GAsynclock[lockid].Unlock()
	return
}
