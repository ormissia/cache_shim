package cache_shim

import (
	"encoding/json"
	"time"
)

// Delete Delete
func Delete(entity CacheType) error {
	_, err := CacheClient().Del(entity.CacheKey())

	if err := entity.Delete(); err != nil {
		return err
	}

	go func() {
		time.Sleep(time.Second)
		_, _ = CacheClient().Del(entity.CacheKey())
	}()
	return err
}

// Select by primary key
func Select[T CacheType](entity CacheType) (res T, err error) {
	jsonStr, err := CacheClient().GetString(entity.CacheKey())
	if err != nil {
		err := entity.Select()
		if err != nil {
			return res, err
		}
		go func() {
			// 缓存
			_ = CacheClient().SetString(entity.CacheKey(), string(entityMarshal(entity)))
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

	if err := entity.Update(); err != nil {
		return err
	}

	go func() {
		time.Sleep(time.Second)
		_, _ = CacheClient().Del(entity.CacheKey())
	}()

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
