package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/light-speak/lighthouse/routers"
)

type claim struct {
	UserId uint `json:"user_id"`
	jwt.RegisteredClaims
}

var key []byte

func init() {
	key = []byte(routers.Config.JWT_SECRET)
}
