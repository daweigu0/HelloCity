package ioc

import (
	"HelloCity/internal/utils"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	addr := fmt.Sprintf("%s:%d", utils.Config.GetString("redis.host"),
		utils.Config.GetInt("redis.port"))
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: utils.Config.GetString("redis.password"),
		Username: utils.Config.GetString("redis.username"),
		DB:       utils.Config.GetInt("redis.database"),
	})
}
