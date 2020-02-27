package U_Redis

import (
	//"fmt"
	"testing"
	"time"

	//"strings"
	//"flag"
	//"os"
	//"github.com/globalsign/mgo"
	//"github.com/globalsign/mgo/bson"
	"github.com/gomodule/redigo/redis"
	//"encoding/json"
	"common/Log"
	"strconv"
)

/*
	cmd: go test -v filename
*/

const (
	Addr string = ":6379"
)

func DialDefaultServer() (redis.Conn, error) {

	c, err := redis.Dial("tcp", Addr, redis.DialReadTimeout(1*time.Second), redis.DialWriteTimeout(1*time.Second))
	if err != nil {
		return nil, err
	}
	c.Do("FLUSHDB")
	return c, nil
}

func Test1(t *testing.T) {
	t.Log("[Test_1] start.")

	c, err := DialDefaultServer()
	if err != nil {
		t.Errorf("connect database err: %v.", err)
		return
	}

	defer c.Close()
	c.Do("SET", "Key", "test001", "EX", "100")
	_, err = c.Do("HSET", "myh", "testval", "test004")
	if err != nil {
		t.Errorf("Expected err for HSET on string key.")
		return
	}

	if c.Err() != nil {
		t.Errorf("Conn has Err()=%v, expect nil.", c.Err())
		return
	}

	_, err = c.Do("SET", "key", "test003")
	if err != nil {
		t.Errorf("Do(SET, key, test003) returned errror %v, expected nil.", err)
		return
	}

	outdata1, err := redis.String(c.Do("GET", "key"))
	t.Log("Do(GET, key) data: ", outdata1)

	outdata2, err := redis.String(c.Do("HGET", "myh", "testval"))
	t.Log("Do(HGET, key) data: ", outdata2)
	return
}

func Test2(t *testing.T) {
	t.Log("[Test2] start.")

	c, err := DialDefaultServer()
	if err != nil {
		t.Errorf("connect database err: %v.", err)
		return
	}

	defer c.Close()
	data, err := c.Do("TIME")
	if err != nil {
		t.Errorf("redis get time err: %v.", err)
		return
	}

	var (
		t1 int64
		t2 int64
	)
	for idx, item := range data.([]interface{}) {
		t, err := strconv.Atoi(string(item.([]byte)))
		if err != nil {
			continue
		}

		if idx == 0 {
			t1 += int64(t)
		} else {
			t2 += int64(t) * 1e3
		}
	}

	Log.FmtPrintf("now time: %v.", time.Unix(int64(t1), t2))
	return
}

func init() {

}
