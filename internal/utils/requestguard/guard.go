package requestguard

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Guard struct {
	rdb    *redis.Client
	window time.Duration
}

func New(rdb *redis.Client, window time.Duration) *Guard {
	return &Guard{
		rdb:    rdb,
		window: window,
	}
}

func (g *Guard) Allow(key string) bool {
	result, err := g.rdb.SetArgs(context.Background(), key, "1", redis.SetArgs{
		Mode: "NX",
		TTL:  g.window,
	}).Result()
	return err == nil && result == "OK"
}
