package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/light-speak/lighthouse/env"
)

type claim struct {
	UserId int64 `json:"user_id"`
	jwt.RegisteredClaims
}

var key []byte

func init() {
	key = []byte(env.GetEnv("JWT_SECRET", "IWY@*3JUI#d309HhefzX2WpLtPKtD!hn"))
}
