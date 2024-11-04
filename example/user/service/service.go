package service

import (
	_ "user/repo"
	"user/resolver"

	"github.com/light-speak/lighthouse/handler"
)

func StartService() {
	resolver := &resolver.Resolver{
		Www: "fuck",
	}
	handler.StartService(resolver)
}
