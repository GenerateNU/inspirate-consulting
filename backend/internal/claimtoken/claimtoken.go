package claimtoken

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid claim token")

type Claims struct {
	Email        string `json:"email"`
	TempPassword string `json:"temp_password"`
	jwt.RegisteredClaims
}

func Sign(secret string, email string, tempPassword string, ttl time.Duration) (string, error) {
	claims := Claims{
		Email:        email,
		TempPassword: tempPassword,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func Verify(secret string, tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
