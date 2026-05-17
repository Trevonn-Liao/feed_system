package service

import (
	"log"
	"net/http"
	"strings"
	"time"

	"feed_system/internal/errs"
	"feed_system/internal/model"
	"feed_system/internal/repository"
	"feed_system/internal/response"
	"feed_system/internal/utils/jwt"
	"feed_system/internal/utils/password"
	"feed_system/internal/utils/requestguard"
	"feed_system/internal/utils/snowflake"

	"github.com/gin-gonic/gin"
)

type AuthService struct {
	users  repository.UserRepository
	tokens *jwt.Manager
	bloom  interface {
		AddWithError(string) error
		MightContain(string) (bool, error)
	}
	snowflake *snowflake.Generator
	guard     *requestguard.Guard
}

func NewAuthService(users repository.UserRepository, tokens *jwt.Manager, bloom interface {
	AddWithError(string) error
	MightContain(string) (bool, error)
}, snowflake *snowflake.Generator, guard *requestguard.Guard) *AuthService {
	return &AuthService{
		users:     users,
		tokens:    tokens,
		bloom:     bloom,
		snowflake: snowflake,
		guard:     guard,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RegisterResponse struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

type LoginResponse struct {
	UserID         int64  `json:"user_id"`
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	AccessExpires  int64  `json:"access_expires_in"`
	RefreshExpires int64  `json:"refresh_expires_in"`
}

func (s *AuthService) Register(c *gin.Context) error {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return errs.BadRequest("invalid register request")
	}

	username := strings.TrimSpace(req.Username)
	if username == "" {
		return errs.BadRequest("username required")
	}

	registerKey := "register:" + username
	if !s.guard.Allow(registerKey) {
		return errs.New(42900, "too many repeated requests", http.StatusTooManyRequests)
	}

	existing, err := s.users.FindByUsername(username)
	if err != nil {
		return errs.Wrap(err, errs.CodeInternal, "query user failed", http.StatusInternalServerError)
	}
	if existing != nil {
		return errs.New(40901, "username already exists", http.StatusConflict)
	}

	hash, err := password.Hash(req.Password)
	if err != nil {
		return errs.Wrap(err, errs.CodeInternal, "hash password failed", http.StatusInternalServerError)
	}

	user := &model.User{
		ID:           s.snowflake.NextID(),
		Username:     username,
		PasswordHash: hash,
	}
	if err := s.users.Create(user); err != nil {
		return errs.Wrap(err, errs.CodeInternal, "create user failed", http.StatusInternalServerError)
	}

	if err := s.bloom.AddWithError("login:" + username); err != nil {
		return errs.Wrap(err, errs.CodeInternal, "sync bloom failed", http.StatusInternalServerError)
	}

	response.Success(c, RegisterResponse{
		UserID:   user.ID,
		Username: user.Username,
	})
	return nil
}

func (s *AuthService) Login(c *gin.Context) error {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return errs.BadRequest("invalid login request")
	}

	username := strings.TrimSpace(req.Username)
	key := "login:" + username
	if !s.guard.Allow(key) {
		return errs.New(42900, "too many repeated requests", http.StatusTooManyRequests)
	}

	mightContain, err := s.bloom.MightContain(key)
	if err != nil {
		log.Printf("redis bloom query failed, fallback to mysql: %v", err)
	} else if !mightContain {
		return errs.New(40101, "user not found", http.StatusUnauthorized)
	}

	user, err := s.users.FindByUsername(username)
	if err != nil {
		return errs.Wrap(err, errs.CodeInternal, "query user failed", http.StatusInternalServerError)
	}
	if user == nil {
		return errs.New(40101, "user not found", http.StatusUnauthorized)
	}

	if err := password.Compare(user.PasswordHash, req.Password); err != nil {
		return errs.New(40102, "password incorrect", http.StatusUnauthorized)
	}

	access, refresh, err := s.tokens.BuildPair(user.ID)
	if err != nil {
		return errs.Wrap(err, errs.CodeInternal, "build token failed", http.StatusInternalServerError)
	}

	response.Success(c, LoginResponse{
		UserID:         user.ID,
		AccessToken:    access,
		RefreshToken:   refresh,
		AccessExpires:  int64(s.tokensAccessTTLSeconds()),
		RefreshExpires: int64(s.tokensRefreshTTLSeconds()),
	})
	return nil
}

func (s *AuthService) Refresh(c *gin.Context) error {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return errs.BadRequest("invalid refresh request")
	}

	access, refresh, userID, err := s.tokens.RefreshPair(req.RefreshToken)
	if err != nil {
		return errs.New(40103, "refresh token invalid", http.StatusUnauthorized)
	}

	response.Success(c, LoginResponse{
		UserID:         userID,
		AccessToken:    access,
		RefreshToken:   refresh,
		AccessExpires:  int64(s.tokensAccessTTLSeconds()),
		RefreshExpires: int64(s.tokensRefreshTTLSeconds()),
	})
	return nil
}

func (s *AuthService) tokensAccessTTL() time.Duration {
	return s.tokens.AccessTTL()
}

func (s *AuthService) tokensRefreshTTL() time.Duration {
	return s.tokens.RefreshTTL()
}

func MustSeedDemoUser(repo repository.UserRepository, bloom interface{ Add(string) }, snowflake *snowflake.Generator) {
	if repo == nil || bloom == nil || snowflake == nil {
		return
	}
}

func (s *AuthService) tokensAccessTTLSeconds() int64 {
	return int64(s.tokensAccessTTL().Seconds())
}

func (s *AuthService) tokensRefreshTTLSeconds() int64 {
	return int64(s.tokensRefreshTTL().Seconds())
}
