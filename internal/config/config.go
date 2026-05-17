package config

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"go.yaml.in/yaml/v3"
)

const (
	exampleConfigPath  = "configs/config.example.yaml"
	localConfigPath    = "configs/config.local.yaml"
	devLocalConfigPath = "configs/dev_local.yaml"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	MySQL     MySQLConfig     `yaml:"mysql"`
	Redis     RedisConfig     `yaml:"redis"`
	JWT       JWTConfig       `yaml:"jwt"`
	Bloom     BloomConfig     `yaml:"bloom"`
	Snowflake SnowflakeConfig `yaml:"snowflake"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
}

type MySQLConfig struct {
	DSN string `yaml:"dsn"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type JWTConfig struct {
	AccessSecret      string        `yaml:"access_secret"`
	RefreshSecret     string        `yaml:"refresh_secret"`
	AccessTokenTTL    time.Duration `yaml:"-"`
	RefreshTokenTTL   time.Duration `yaml:"-"`
	AccessTTLSeconds  int           `yaml:"access_ttl_seconds"`
	RefreshTTLSeconds int           `yaml:"refresh_ttl_seconds"`
	Issuer            string        `yaml:"issuer"`
}

type BloomConfig struct {
	Capacity  uint    `yaml:"capacity"`
	ErrorRate float64 `yaml:"error_rate"`
}

type SnowflakeConfig struct {
	WorkerID     int64 `yaml:"worker_id"`
	DatacenterID int64 `yaml:"datacenter_id"`
}

func Load() Config {
	cfg := defaultConfig()
	root := projectRoot()
	_ = mergeYAMLFile(filepath.Join(root, exampleConfigPath), &cfg)
	_ = mergeYAMLFile(filepath.Join(root, localConfigPath), &cfg)
	_ = mergeYAMLFile(filepath.Join(root, devLocalConfigPath), &cfg)
	applyEnv(&cfg)
	normalize(&cfg)
	return cfg
}

func defaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Addr: ":8080",
		},
		MySQL: MySQLConfig{
			DSN: "root:root@tcp(127.0.0.1:3306)/feed_system?charset=utf8mb4&parseTime=True&loc=Local",
		},
		Redis: RedisConfig{
			Addr: "127.0.0.1:6379",
			DB:   0,
		},
		JWT: JWTConfig{
			AccessTTLSeconds:  15 * 60,
			RefreshTTLSeconds: 7 * 24 * 60 * 60,
			Issuer:            "feed_system",
		},
		Bloom: BloomConfig{
			Capacity:  100000,
			ErrorRate: 0.001,
		},
		Snowflake: SnowflakeConfig{
			WorkerID:     1,
			DatacenterID: 1,
		},
	}
}

func mergeYAMLFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	return yaml.Unmarshal(data, cfg)
}

func projectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			return "."
		}
		wd = parent
	}
}

func applyEnv(cfg *Config) {
	cfg.Server.Addr = getEnv("APP_ADDR", cfg.Server.Addr)
	cfg.MySQL.DSN = getEnv("MYSQL_DSN", cfg.MySQL.DSN)
	cfg.Redis.Addr = getEnv("REDIS_ADDR", cfg.Redis.Addr)
	cfg.Redis.Password = getEnv("REDIS_PASSWORD", cfg.Redis.Password)
	cfg.Redis.DB = getIntEnv("REDIS_DB", cfg.Redis.DB)

	cfg.JWT.AccessSecret = getEnv("JWT_ACCESS_SECRET", cfg.JWT.AccessSecret)
	cfg.JWT.RefreshSecret = getEnv("JWT_REFRESH_SECRET", cfg.JWT.RefreshSecret)
	cfg.JWT.AccessTTLSeconds = getIntEnv("JWT_ACCESS_TTL_SECONDS", cfg.JWT.AccessTTLSeconds)
	cfg.JWT.RefreshTTLSeconds = getIntEnv("JWT_REFRESH_TTL_SECONDS", cfg.JWT.RefreshTTLSeconds)
	cfg.JWT.Issuer = getEnv("JWT_ISSUER", cfg.JWT.Issuer)

	cfg.Bloom.Capacity = uint(getIntEnv("BLOOM_CAPACITY", int(cfg.Bloom.Capacity)))
	cfg.Bloom.ErrorRate = getFloatEnv("BLOOM_ERROR_RATE", cfg.Bloom.ErrorRate)
	cfg.Snowflake.WorkerID = int64(getIntEnv("SNOWFLAKE_WORKER_ID", int(cfg.Snowflake.WorkerID)))
	cfg.Snowflake.DatacenterID = int64(getIntEnv("SNOWFLAKE_DATACENTER_ID", int(cfg.Snowflake.DatacenterID)))
}

func normalize(cfg *Config) {
	cfg.JWT.AccessTokenTTL = time.Duration(cfg.JWT.AccessTTLSeconds) * time.Second
	cfg.JWT.RefreshTokenTTL = time.Duration(cfg.JWT.RefreshTTLSeconds) * time.Second
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

func getFloatEnv(key string, fallback float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fallback
	}
	return n
}
