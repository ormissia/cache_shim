package cache_shim

type CacheType interface {
	CacheKey() string

	Delete() error
	Select() error
	Update() error
}
