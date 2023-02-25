package redis_ex

import (
	"errors"
)

type CacheStorage map[string]interface{}

// RDB 模拟一个简单的Redis客户端
var RDB CacheStorage

func Init() {
	RDB = make(map[string]interface{})
}

func (c *CacheStorage) Del(key string) (int64, error) {
	delete(RDB, key)
	return 0, nil
}

func (c *CacheStorage) SetString(key, value string, expire int) (err error) {
	RDB[key] = value
	return nil
}

func (c *CacheStorage) GetString(key string) (string, error) {
	res, ok := RDB[key]
	if !ok {
		return "", errors.New("key: " + key + " not found")
	}
	return res.(string), nil
}
