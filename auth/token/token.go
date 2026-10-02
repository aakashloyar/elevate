package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Claims struct {
	Subject   string `json:"sub"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
}

var ErrInvalid = errors.New("invalid access token")

func Sign(secret string, claims Claims) (string, error) {
	if secret == "" || claims.Subject == "" {
		return "", ErrInvalid
	}
	h := encode(map[string]string{"alg": "HS256", "typ": "JWT"})
	b := encode(claims)
	input := h + "." + b
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(input))
	return input + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}
func Parse(secret, value string, now time.Time) (Claims, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 3 || secret == "" {
		return Claims{}, ErrInvalid
	}
	input := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(input))
	got, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || !hmac.Equal(got, mac.Sum(nil)) {
		return Claims{}, ErrInvalid
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, ErrInvalid
	}
	var c Claims
	if json.Unmarshal(raw, &c) != nil || c.Subject == "" || c.ExpiresAt <= now.Unix() {
		return Claims{}, ErrInvalid
	}
	return c, nil
}
func Bearer(header string) (string, error) {
	const p = "Bearer "
	if !strings.HasPrefix(header, p) {
		return "", fmt.Errorf("missing bearer token")
	}
	return strings.TrimSpace(strings.TrimPrefix(header, p)), nil
}
func encode(v any) string {
	raw, _ := json.Marshal(v)
	return base64.RawURLEncoding.EncodeToString(raw)
}
