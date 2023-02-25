package cache_shim

import (
	"encoding/json"
	"time"
)

// Insert Insert
func Insert(entity CacheType) error {
	// TODO 将数据添加到布隆过滤器
	// TODO 此处如果数据如果后续被删除，布隆过滤器中状态仍为存在
	// TODO 或者使用缓存空对象解决穿透问题
	return entity.Insert()
}

// Delete Delete
func Delete(entity CacheType) error {
	if _, err := CacheClient().Del(entity.CacheKey()); err != nil {
		return err
	}

	defer func() {
		go func() {
			time.Sleep(time.Second)
			_, _ = CacheClient().Del(entity.CacheKey())
		}()
	}()

	if err := entity.Delete(); err != nil {
		return err
	}

	return nil
}

// Select by primary key
func Select[T CacheType](entity CacheType) (res T, err error) {
	jsonStr, err := CacheClient().GetString(entity.CacheKey())
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
			_ = CacheClient().SetString(entity.CacheKey(), string(entityMarshal(entity)), entity.Expiration())
		}()
		return entity.(T), nil
	}

	return entityUnMarshal[T]([]byte(jsonStr))
}

// Update 延时双删策略
func Update(entity CacheType) error {
	if _, err := CacheClient().Del(entity.CacheKey()); err != nil {
		return err
	}

	defer func() {
		go func() {
			time.Sleep(time.Second)
			_, _ = CacheClient().Del(entity.CacheKey())
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
