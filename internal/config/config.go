package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server    ServerConfig
	MySQL     MySQLConfig
	Redis     RedisConfig
	JWT       JWTConfig
	Bloom     BloomConfig
	Snowflake SnowflakeConfig
}

type ServerConfig struct {
	Addr string
}

type MySQLConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	AccessSecret    string
	RefreshSecret   string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
}

type BloomConfig struct {
	Capacity uint
	Bits     uint
	Hashes   uint
}

type SnowflakeConfig struct {
	WorkerID     int64
	DatacenterID int64
}

func Load() Config {
	return Config{
		Server: ServerConfig{
			Addr: getEnv("APP_ADDR", ":8080"),
		},
		MySQL: MySQLConfig{
			DSN: getEnv("MYSQL_DSN", "root:root@tcp(127.0.0.1:3306)/feed_system?charset=utf8mb4&parseTime=True&loc=Local"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			AccessSecret:    getEnv("JWT_ACCESS_SECRET", "change-me-access-secret"),
			RefreshSecret:   getEnv("JWT_REFRESH_SECRET", "change-me-refresh-secret"),
			AccessTokenTTL:  getDurationEnv("JWT_ACCESS_TTL_SECONDS", 15*60),
			RefreshTokenTTL: getDurationEnv("JWT_REFRESH_TTL_SECONDS", 7*24*60*60),
			Issuer:          getEnv("JWT_ISSUER", "feed_system"),
		},
		Bloom: BloomConfig{
			Capacity: uint(getIntEnv("BLOOM_CAPACITY", 100000)),
			Bits:     uint(getIntEnv("BLOOM_BITS", 1000000)),
			Hashes:   uint(getIntEnv("BLOOM_HASHES", 7)),
		},
		Snowflake: SnowflakeConfig{
			WorkerID:     int64(getIntEnv("SNOWFLAKE_WORKER_ID", 1)),
			DatacenterID: int64(getIntEnv("SNOWFLAKE_DATACENTER_ID", 1)),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func getDurationEnv(key string, fallbackSeconds int) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return time.Duration(fallbackSeconds) * time.Second
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return time.Duration(fallbackSeconds) * time.Second
	}
	return time.Duration(n) * time.Second
}
