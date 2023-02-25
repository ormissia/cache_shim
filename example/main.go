package main

import (
	"fmt"
	"github.com/ormissia/cache_shim"
	"github.com/ormissia/cache_shim/example/db_ex"
	"github.com/ormissia/cache_shim/example/redis_ex"
)

func main() {
	redis_ex.Init()
	cache_shim.InitCacheClient(&redis_ex.RDB)

	t1 := &db_ex.UserEx{
		ID:   1,
		Age:  123,
		Name: "ormissia",
	}

	// 插入之后再通过相同的ID去查询
	_ = cache_shim.Insert(t1)

	t2 := &db_ex.UserEx{
		ID: 1,
	}
	t, err := cache_shim.Select[*db_ex.UserEx](t2)

	fmt.Println()
	fmt.Printf("t.type: %T\tt: %v\terr: %v", t, t, err)
	fmt.Println()
}
