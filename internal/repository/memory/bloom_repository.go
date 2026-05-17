package memory

import (
	"context"
	"errors"
	"strings"

	"github.com/redis/go-redis/v9"
)

type BloomRepository struct {
	rdb         *redis.Client
	key         string
	errorRate   float64
	capacity    uint
	initialized bool
}

func NewBloomRepository(rdb *redis.Client, key string, capacity uint, errorRate float64) *BloomRepository {
	if capacity == 0 {
		capacity = 100000
	}
	if errorRate <= 0 {
		errorRate = 0.001
	}
	return &BloomRepository{
		rdb:       rdb,
		key:       key,
		capacity:  capacity,
		errorRate: errorRate,
	}
}

func (r *BloomRepository) Init(ctx context.Context) error {
	err := r.rdb.Do(ctx, "BF.RESERVE", r.key, r.errorRate, r.capacity).Err()
	if err != nil && !isRedisBloomExistsErr(err) {
		return err
	}
	r.initialized = true
	return nil
}

func (r *BloomRepository) Add(value string) {
	_ = r.AddWithError(value)
}

func (r *BloomRepository) AddWithError(value string) error {
	if err := r.ensureInitialized(context.Background()); err != nil {
		return err
	}
	return r.rdb.Do(context.Background(), "BF.ADD", r.key, value).Err()
}

func (r *BloomRepository) MightContain(value string) (bool, error) {
	if err := r.ensureInitialized(context.Background()); err != nil {
		return false, err
	}
	v, err := r.rdb.Do(context.Background(), "BF.EXISTS", r.key, value).Int()
	if err != nil {
		return false, err
	}
	return v == 1, nil
}

func BloomKey(namespace string) string {
	return "bloom:" + namespace
}

func (r *BloomRepository) ensureInitialized(ctx context.Context) error {
	if r.initialized {
		return nil
	}
	return r.Init(ctx)
}

func isRedisBloomExistsErr(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return errors.Is(err, redis.Nil) ||
		strings.Contains(msg, "item exists") ||
		strings.Contains(msg, "key exists") ||
		strings.Contains(msg, "exists")
}
