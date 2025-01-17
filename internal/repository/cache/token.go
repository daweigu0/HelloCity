package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type TokenCache interface {
	Set(ctx context.Context, prefix, key, value string) error
	Get(ctx context.Context, prefix, key string) (string, error)
	Verify(ctx context.Context, prefix, key, value string) (bool, error)
	Del(ctx context.Context, prefix, key string) error
}
type RedisTokenCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (c *RedisTokenCache) Get(ctx context.Context, prefix, key string) (string, error) {
	res, err := c.cmd.Get(ctx, c.key(prefix, key)).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}

func (c *RedisTokenCache) Del(ctx context.Context, prefix, key string) error {
	_, err := c.cmd.Del(ctx, prefix, key).Result()
	if err != nil {
		return err
	}
	return nil
}

func NewTokenCache(cmd redis.Cmdable, expiration time.Duration) TokenCache {
	return &RedisTokenCache{
		cmd:        cmd,
		expiration: expiration,
	}
}

func (c *RedisTokenCache) Set(ctx context.Context, prefix, key, value string) error {
	res, err := c.cmd.Set(ctx, c.key(prefix, key), value, c.expiration).Result()
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func (c *RedisTokenCache) Verify(ctx context.Context, prefix, key, value string) (bool, error) {
	res, err := c.cmd.Get(ctx, c.key(prefix, key)).Result()
	if err != nil {
		return false, err
	}
	fmt.Println(res)
	if res == value {
		return true, nil
	} else {
		return false, nil
	}
}

func (c *RedisTokenCache) key(prefix, key string) string {
	return fmt.Sprintf("token_cache:%s:%s", prefix, key)
}
