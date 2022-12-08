package db_ex

import (
	"errors"
	"strconv"
)

// 模拟一个简单的DB
var db = make(map[int]UserEx)

type UserEx struct {
	ID   int    `json:"id"`
	Age  int    `json:"age"`
	Name string `json:"name"`
}

func (t *UserEx) CacheKey() string {
	return strconv.Itoa(t.ID)
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
	t = &res
	return nil
}

func (t *UserEx) Update() error {
	db[t.ID] = *t
	return nil
}

func (t *UserEx) Insert() error {
	db[t.ID] = *t
	return nil
}
