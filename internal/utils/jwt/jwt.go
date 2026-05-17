package jwt

import (
	"errors"
	"strconv"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type Claims struct {
	UserID int64     `json:"user_id"`
	Typ    TokenType `json:"typ"`
	gojwt.RegisteredClaims
}

type Manager struct {
	accessSecret  []byte
	refreshSecret []byte
	issuer        string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func New(accessSecret, refreshSecret, issuer string, accessTTL, refreshTTL time.Duration) *Manager {
	return &Manager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		issuer:        issuer,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (m *Manager) BuildPair(userID int64) (string, string, error) {
	access, err := m.build(userID, TokenTypeAccess, m.accessTTL, m.accessSecret)
	if err != nil {
		return "", "", err
	}
	refresh, err := m.build(userID, TokenTypeRefresh, m.refreshTTL, m.refreshSecret)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func (m *Manager) RefreshPair(refreshToken string) (string, string, int64, error) {
	claims, err := m.Parse(refreshToken, TokenTypeRefresh)
	if err != nil {
		return "", "", 0, err
	}
	access, refresh, err := m.BuildPair(claims.UserID)
	if err != nil {
		return "", "", 0, err
	}
	return access, refresh, claims.UserID, nil
}

func (m *Manager) Parse(token string, typ TokenType) (*Claims, error) {
	claims := &Claims{}
	parsed, err := gojwt.ParseWithClaims(token, claims, func(t *gojwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		if typ == TokenTypeRefresh {
			return m.refreshSecret, nil
		}
		return m.accessSecret, nil
	}, gojwt.WithIssuer(m.issuer))
	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.Typ != typ {
		return nil, errors.New("token type mismatch")
	}
	return claims, nil
}

func (m *Manager) build(userID int64, typ TokenType, ttl time.Duration, secret []byte) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Typ:    typ,
		RegisteredClaims: gojwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   strconv.FormatInt(userID, 10),
			IssuedAt:  gojwt.NewNumericDate(now),
			NotBefore: gojwt.NewNumericDate(now),
			ExpiresAt: gojwt.NewNumericDate(now.Add(ttl)),
		},
	}
	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (m *Manager) AccessTTL() time.Duration {
	return m.accessTTL
}

func (m *Manager) RefreshTTL() time.Duration {
	return m.refreshTTL
}
