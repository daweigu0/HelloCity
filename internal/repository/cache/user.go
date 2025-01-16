package cache

import (
	"HelloCity/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, uid uint64) (domain.User, error)
	Set(ctx context.Context, du domain.User) error
	Del(ctx context.Context, id uint64) error
}

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (c *RedisUserCache) Del(ctx context.Context, id uint64) error {
	return c.cmd.Del(ctx, c.key(id)).Err()
}

func (c *RedisUserCache) Get(ctx context.Context, uid uint64) (domain.User, error) {
	key := c.key(uid)
	// 我假定这个地方用 JSON 来
	data, err := c.cmd.Get(ctx, key).Result()
	//data, err := c.cmd.Get(ctx, firstKey).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	//if err != nil {
	//	return domain.User{}, err
	//}
	//return u, nil
	return u, err
}

func (c *RedisUserCache) Set(ctx context.Context, du domain.User) error {
	key := c.key(du.ID)
	// 我假定这个地方用 JSON
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

func (c *RedisUserCache) key(uid uint64) string {
	return fmt.Sprintf("user:info:%d", uid)
}

func NewUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}
