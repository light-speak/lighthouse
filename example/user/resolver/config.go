package resolver

import (
	"user/graph/generate"

	"github.com/light-speak/lighthouse/db"
)

func LoadConfig() generate.Config {
	c := generate.Config{
		Resolvers: &Resolver{
			Db: db.GetDb(),
		},
	}
	return c
}
