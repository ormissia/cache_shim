package main

import (
	"fmt"
	"time"

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
		Name: "xiaoming",
	}

	_ = t1.Insert()
	t, err := cache_shim.Select[*db_ex.UserEx](t1)

	fmt.Println()
	fmt.Printf("t.type: %T\tt: %v\terr: %v", t, t, err)
	fmt.Println()

	time.Sleep(time.Second * 3)
}
