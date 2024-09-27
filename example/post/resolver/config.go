package resolver

import (
	"post/graph/generate"
	"search/kitex_gen/search/searchservice"

	"github.com/cloudwego/kitex/client"
	"github.com/light-speak/lighthouse/db"
	"github.com/light-speak/lighthouse/log"
)

func LoadConfig() generate.Config {
	searchClient, err := searchservice.NewClient("SearchService", client.WithHostPorts("127.0.0.1:8888"))
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	c := generate.Config{
		Resolvers: &Resolver{
			Db:           db.GetDb(),
			SearchClient: &searchClient,
		},
	}
	return c
}
