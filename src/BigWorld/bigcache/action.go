package bigcache

import (
	. "common/S2SMessage"
	"log"
)

type CacheOperCB func(c *CacheOperation)
var CacheOperCall map[ECacheOper]CacheOperCB = map[ECacheOper]CacheOperCB{
	ECacheOper_Add: AddBigC,
	ECacheOper_Delete: DeleteBigC,
	ECacheOper_Update: UpdateBigC,
	ECacheOper_Select: SelectBigC,
}

func SelectOper(c *CacheOperation) {
	cb, ok := CacheOperCall[c.Oper]
	if !ok || cb == nil {
		log.Fatal("can not find cache cb, oper: ", c.Oper)
		return
	}
	cb(c)
}

func AddBigC(c *CacheOperation) {
	
}

func DeleteBigC(c *CacheOperation) {

}

func UpdateBigC(c *CacheOperation) {

}

func SelectBigC(c *CacheOperation) {

}