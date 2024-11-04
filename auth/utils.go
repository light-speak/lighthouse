package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetToken(userId int64) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, buildClaims(userId)).SignedString(key)
}

func buildClaims(userId int64) *claim {
	now := time.Now()
	return &claim{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * 30)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "lighthouse",
		},
	}
}

func GetUserId(token string) (int64, error) {
	t, err := jwt.ParseWithClaims(token, &claim{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := t.Claims.(*claim); ok && t.Valid {
		return claims.UserId, nil
	}
	return 0, errors.New("invalid token")
}
