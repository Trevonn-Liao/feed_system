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
	ok, err := g.rdb.SetNX(context.Background(), key, "1", g.window).Result()
	return err == nil && ok
}
