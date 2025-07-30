package sysinit

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

var (
	RedisCli       *redis.Client
	redisOnce      sync.Once
	redisInitError error
)

func InitRedis() error {
	redisOnce.Do(func() {
		if GCF == nil {
			redisInitError = fmt.Errorf("configuration system not initialized")
			return
		}

		addr := GCF.UString("redis.addr")
		if addr == "" {
			redisInitError = fmt.Errorf("redis.addr configuration is required")
			return
		}

		RedisCli = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: GCF.UString("redis.password"),
			DB:       GCF.UInt("redis.database"),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if _, err := RedisCli.Ping(ctx).Result(); err != nil {
			redisInitError = fmt.Errorf("failed to connect to Redis: %w", err)
			RedisCli = nil
		}
	})

	return redisInitError
}
