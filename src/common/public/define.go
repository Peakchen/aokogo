package public

type IDBCache interface {
	CacheKey() string
	MainModel() string
	SubModel() string
}
