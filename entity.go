package cache_shim

// CacheType 需要缓存的实体接口定义
type CacheType interface {
	CacheKey() string
	Expiration() int

	Insert() error
	Delete() error
	Select() error
	Update() error
}
