package public

import "common/ado"

/*

 */
type IDBCache interface {
	Identify() string
	MainModel() string
	SubModel() string
}

type TCommonRedisCache struct {
	ado.IDBModule
}

func (this *TCommonRedisCache) Identify() string {
	return ""
}

func (this *TCommonRedisCache) MainModel() string {
	return ""
}

func (this *TCommonRedisCache) SubModel() string {
	return ""
}

type UpdateDBCacheCallBack func(string, string, []byte) bool
