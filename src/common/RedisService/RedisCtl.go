package RedisService

/*
	ctrl some redis server change channl
*/

type TRedisCtrl struct {
	redisCns    []*TRedisConn
	curRedisIdx int
}

func LazyInit() (redctl *TRedisCtrl) {
	redctl = &TRedisCtrl{}
	redctl.redisCns = []*TRedisConn{}
	return
}

func (this *TRedisCtrl) AddRedisConn(c *TRedisConn) {
	this.redisCns = append(this.redisCns, c)
}

func (this *TRedisCtrl) ChangeNextChannl() (c *TRedisConn) {
	c = nil
	if this.curRedisIdx == len(this.redisCns) {
		return
	}
	mod := (this.curRedisIdx + 1) % len(this.redisCns)
	if mod == 0 {
		this.curRedisIdx = len(this.redisCns) - 1
	} else {
		this.curRedisIdx++
	}
	return this.redisCns[this.curRedisIdx]
}
