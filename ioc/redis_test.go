package ioc

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestInitRedis(t *testing.T) {
	rdb := InitRedis()
	ctx := context.Background()
	result, err := rdb.Set(ctx, "key", "value", 5*time.Minute).Result()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(result)
	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(val)
}
