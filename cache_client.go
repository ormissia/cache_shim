package cache_shim

import "sync"

type CacheClint interface {
	Del(key string) (int64, error)
	SetString(key, value string, expire int) error
	GetString(k string) (string, error)
}

var (
	client CacheClint
	once   sync.Once
)

func InitCacheClient(cacheClient CacheClint) {
	once.Do(func() {
		client = cacheClient
	})
}

func CacheClient() CacheClint {
	return client
}
