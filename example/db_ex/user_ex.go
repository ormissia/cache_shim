package db_ex

import (
	"errors"
	"strconv"
)

// 模拟一个简单的持久层DB
var db = make(map[int]UserEx)

type UserEx struct {
	ID   int    `json:"id"`
	Age  int    `json:"age"`
	Name string `json:"name"`
}

func (t *UserEx) CacheKey() string {
	return strconv.Itoa(t.ID)
}

func (t *UserEx) Expiration() int {
	return 1
}

func (t *UserEx) Insert() error {
	db[t.ID] = *t
	return nil
}

func (t *UserEx) Delete() error {
	delete(db, t.ID)
	return nil
}

func (t *UserEx) Select() error {
	res, ok := db[t.ID]
	if !ok {
		return errors.New("not found")
	}

	// 直接修改指针会有问题，一般查数据库的结果会直接通过反射放到t中，无需将t指向新的地址
	// t=&res
	t.Age = res.Age
	t.Name = res.Name
	return nil
}

func (t *UserEx) Update() error {
	db[t.ID] = *t
	return nil
}
