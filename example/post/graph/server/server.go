package server

import (
	"context"
	"net/http"
	"post/graph/generate"
	"post/graph/models"
	"post/resolver"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/light-speak/lighthouse/db"
	"github.com/light-speak/lighthouse/env"
	"github.com/light-speak/lighthouse/graphql/middleware"
	"github.com/light-speak/lighthouse/log"
)

func StartServer() error {
	port := env.Getenv("PORT", "4000")
	c := resolver.LoadConfig()
	err := LoadSearchService(c.Resolvers.(*resolver.Resolver))
	if err != nil {
		log.Error(err.Error())
		return err
	}

	SetDirective(&c.Directives)
	router := middleware.GetRouter()

	srv := handler.NewDefaultServer(generate.NewExecutableSchema(c))

	router.Handle("/", playground.ApolloSandboxHandler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Info("成功运行， 点击进入 GraphQL playground ： http://localhost:%s/ ", port)
	return http.ListenAndServe(":"+port, generate.Middleware(db.GetDb(), router))
}

func LoadSearchService(r *resolver.Resolver) error {
	searchClient := *r.SearchClient
	var post models.Post
	mapping := post.Mapping()
	return searchClient.CreateIndex(context.Background(), post.IndexName(), mapping)
}
