package bigcache

// enum cache number
const (
	ConstI32PlayerName int32 = 1
	// ...
)

// enum cache name
const (
	ConstStrPlayerName string = "PlayerName"
	// ...
)

// cache table
var GCacheTab map[int32]string = map[int32]string{
	ConstI32PlayerName: ConstStrPlayerName,
}
