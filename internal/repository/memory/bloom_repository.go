package memory

import (
	"context"
	"hash/fnv"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type BloomRepository struct {
	rdb     *redis.Client
	key     string
	bitSize uint
	hashes  uint
}

func NewBloomRepository(rdb *redis.Client, key string, bitSize, hashes uint) *BloomRepository {
	if bitSize == 0 {
		bitSize = 1
	}
	if hashes == 0 {
		hashes = 1
	}
	return &BloomRepository{
		rdb:     rdb,
		key:     key,
		bitSize: bitSize,
		hashes:  hashes,
	}
}

func (r *BloomRepository) Add(value string) {
	ctx := context.Background()
	h1, h2 := bloomHashes(value)
	for i := uint(0); i < r.hashes; i++ {
		index := (h1 + uint64(i)*h2) % uint64(r.bitSize)
		if err := r.rdb.SetBit(ctx, r.key, int64(index), 1).Err(); err != nil {
			return
		}
	}
}

func (r *BloomRepository) AddWithError(value string) error {
	ctx := context.Background()
	h1, h2 := bloomHashes(value)
	for i := uint(0); i < r.hashes; i++ {
		index := (h1 + uint64(i)*h2) % uint64(r.bitSize)
		if err := r.rdb.SetBit(ctx, r.key, int64(index), 1).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (r *BloomRepository) MightContain(value string) (bool, error) {
	ctx := context.Background()
	h1, h2 := bloomHashes(value)
	for i := uint(0); i < r.hashes; i++ {
		index := (h1 + uint64(i)*h2) % uint64(r.bitSize)
		v, err := r.rdb.GetBit(ctx, r.key, int64(index)).Result()
		if err != nil {
			return false, err
		}
		if v == 0 {
			return false, nil
		}
	}
	return true, nil
}

func (r *BloomRepository) Load(values []string) {
	for _, v := range values {
		r.Add(v)
	}
}

func bloomHashes(value string) (uint64, uint64) {
	h1 := fnv.New64a()
	_, _ = h1.Write([]byte(value))
	sum1 := h1.Sum64()

	h2 := fnv.New64()
	_, _ = h2.Write([]byte(value + "#salt"))
	sum2 := h2.Sum64()
	if sum2 == 0 {
		sum2 = 1
	}
	return sum1, sum2
}

func BloomKey(namespace string) string {
	return "bloom:" + namespace
}

func Uint64ToString(v uint64) string {
	return strconv.FormatUint(v, 10)
}
