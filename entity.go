package cache_shim

type CacheType interface {
	CacheKey() string
	Expiration() int

	Delete() error
	Select() error
	Update() error
}
