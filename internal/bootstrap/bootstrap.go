package bootstrap

import (
	"context"
	"log"
	"time"

	"feed_system/internal/config"
	"feed_system/internal/model"
	"feed_system/internal/repository"
	"feed_system/internal/repository/memory"
	mysqlrepo "feed_system/internal/repository/mysql"
	"feed_system/internal/service"
	"feed_system/internal/utils/jwt"
	"feed_system/internal/utils/requestguard"
	"feed_system/internal/utils/snowflake"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct {
	Config *config.Config
	DB     *gorm.DB
	RDB    *redis.Client
	Bloom  *memory.BloomRepository
	Auth   *service.AuthService
}

func New() *App {
	cfg := config.Load()
	if cfg.JWT.AccessSecret == "" || cfg.JWT.RefreshSecret == "" ||
		cfg.JWT.AccessSecret == "replace-with-your-access-secret" ||
		cfg.JWT.RefreshSecret == "replace-with-your-refresh-secret" {
		log.Fatal("valid jwt access_secret and refresh_secret must be set by configs/dev_local.yaml, configs/config.local.yaml, or environment variables")
	}

	db, err := gorm.Open(mysql.Open(cfg.MySQL.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("open mysql failed: %v", err)
	}

	snow, err := snowflake.New(cfg.Snowflake.WorkerID, cfg.Snowflake.DatacenterID)
	if err != nil {
		log.Fatalf("create snowflake failed: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("connect redis failed: %v", err)
	}

	tokenMgr := jwt.New(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret, cfg.JWT.Issuer, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL)
	bloomRepo := memory.NewBloomRepository(rdb, memory.BloomKey("users"), cfg.Bloom.Capacity, cfg.Bloom.ErrorRate)
	if err := bloomRepo.Init(context.Background()); err != nil {
		log.Fatalf("init redis bloom failed: %v", err)
	}
	guard := requestguard.New(rdb, 2*time.Second)
	userRepo := mysqlrepo.NewUserRepository(db)
	loadBloomFromUsers(userRepo, bloomRepo)
	authSvc := service.NewAuthService(userRepo, tokenMgr, bloomRepo, snow, guard)

	return &App{
		Config: &cfg,
		DB:     db,
		RDB:    rdb,
		Bloom:  bloomRepo,
		Auth:   authSvc,
	}
}

func loadBloomFromUsers(repo repository.UserRepository, bloom *memory.BloomRepository) {
	users, err := repo.ListAll()
	if err != nil {
		log.Printf("load bloom skipped: %v", err)
		return
	}
	for _, user := range users {
		bloom.Add("login:" + user.Username)
	}
}
