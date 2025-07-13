package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var ErrInvalidToken = errors.New("invalid token")
var ErrExpiredToken = errors.New("expired token")

type Claims struct {
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (c *Claims) Valid() error {
	if time.Now().After(c.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

type JWTProcessor struct {
	secretKey string
}

func (j *JWTProcessor) VerifyToken(token string) (*Claims, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(j.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Claims{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := jwtToken.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func NewJWTProcessor(secretKey string) *JWTProcessor {
	return &JWTProcessor{secretKey: secretKey}
}
