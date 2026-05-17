package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type Claims struct {
	UserID int64     `json:"user_id"`
	Typ    TokenType `json:"typ"`
	Sub    string    `json:"sub"`
	Iss    string    `json:"iss"`
	Iat    int64     `json:"iat"`
	Exp    int64     `json:"exp"`
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
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}
	if claims.Typ != typ {
		return nil, errors.New("token type mismatch")
	}
	if claims.Iss != m.issuer {
		return nil, errors.New("issuer mismatch")
	}
	if time.Now().Unix() > claims.Exp {
		return nil, errors.New("token expired")
	}

	secret := m.accessSecret
	if typ == TokenTypeRefresh {
		secret = m.refreshSecret
	}
	expected, err := m.sign(parts[0]+"."+parts[1], secret)
	if err != nil {
		return nil, err
	}
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return nil, errors.New("signature mismatch")
	}
	return &claims, nil
}

func (m *Manager) build(userID int64, typ TokenType, ttl time.Duration, secret []byte) (string, error) {
	now := time.Now().Unix()
	claims := Claims{
		UserID: userID,
		Typ:    typ,
		Sub:    strconv.FormatInt(userID, 10),
		Iss:    m.issuer,
		Iat:    now,
		Exp:    now + int64(ttl.Seconds()),
	}
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	headerJSON, _ := json.Marshal(header)
	payloadJSON, _ := json.Marshal(claims)
	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadJSON)
	signature, err := m.sign(encodedHeader+"."+encodedPayload, secret)
	if err != nil {
		return "", err
	}
	return encodedHeader + "." + encodedPayload + "." + signature, nil
}

func (m *Manager) sign(data string, secret []byte) (string, error) {
	mac := hmac.New(sha256.New, secret)
	if _, err := mac.Write([]byte(data)); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}

func (m *Manager) AccessTTL() time.Duration {
	return m.accessTTL
}

func (m *Manager) RefreshTTL() time.Duration {
	return m.refreshTTL
}
