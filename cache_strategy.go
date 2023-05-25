package cache_shim

import (
	"encoding/json"
	"golang.org/x/sync/singleflight"
	"time"
)

func InitCacheClient(cacheClient CacheClint) *CacheShim {
	return &CacheShim{
		CacheClient: cacheClient,
		sg:          singleflight.Group{},
	}
}

type CacheShim struct {
	CacheClient CacheClint
	sg          singleflight.Group
}

// Insert Insert
func Insert(_ *CacheShim, entity CacheType) error {
	return entity.Insert()
}

// Delete Delete
func Delete(cacheShim *CacheShim, entity CacheType) error {
	if _, err := cacheShim.CacheClient.Del(entity.CacheKey()); err != nil {
		return err
	}

	defer func() {
		go func() {
			time.Sleep(time.Second)
			_, _ = cacheShim.CacheClient.Del(entity.CacheKey())
		}()
	}()

	if err := entity.Delete(); err != nil {
		return err
	}

	return nil
}

// Select by primary key
func Select[T CacheType](cacheShim *CacheShim, entity CacheType) (res T, err error) {
	jsonStr, err := cacheShim.CacheClient.GetString(entity.CacheKey())
	if err != nil {
		// TODO 添加查询失败的情况，判断是否需要缓存数据不存在的标记，防止缓存穿透

		err := entity.Select()
		if err != nil {
			return res, err
		}
		// 缓存到redis
		// TODO 修改redis前加锁，防止出现多实例查询造成的缓存一致性问题
		// 如果发现缓存已锁定，自旋查询
		go func() {
			_ = cacheShim.CacheClient.SetString(entity.CacheKey(), string(entityMarshal(entity)), entity.Expiration())
		}()
		return entity.(T), nil
	}

	return entityUnMarshal[T]([]byte(jsonStr))
}

// Update 延时双删策略
func Update(cacheShim *CacheShim, entity CacheType) error {
	if _, err := cacheShim.CacheClient.Del(entity.CacheKey()); err != nil {
		return err
	}

	defer func() {
		go func() {
			time.Sleep(time.Second)
			_, _ = cacheShim.CacheClient.Del(entity.CacheKey())
		}()
	}()

	if err := entity.Update(); err != nil {
		return err
	}

	return nil
}

func entityMarshal(entity CacheType) []byte {
	res, _ := json.Marshal(entity)
	return res
}

func entityUnMarshal[T CacheType](source []byte) (res T, err error) {
	if err := json.Unmarshal(source, &res); err != nil {
		return res, err
	}
	return res, nil
}
