package cache_shim

// CacheClint 缓存客户端接口定义
type CacheClint interface {
	Del(key string) (int64, error)
	SetString(key, value string, expire int) error
	GetString(k string) (string, error)
}
