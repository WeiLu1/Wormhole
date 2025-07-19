package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

const testSecretKey = "123456789012345678901234456789012"

func createToken() (string, *Claims) {
	claims := &Claims{
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(10 * time.Minute),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, _ := jwtToken.SignedString([]byte(testSecretKey))

	return token, claims
}

func TestVerifyToken(t *testing.T) {
	token, _ := createToken()

	jwtProcessor := NewJWTProcessor(testSecretKey)
	claims, err := jwtProcessor.VerifyToken(token)

	if err != nil {
		t.Errorf("verify token did not work")
	}

	if claims == nil {
		t.Errorf("claims is empty")
	}
}
