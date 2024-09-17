package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/light-speak/lighthouse/db"
	"github.com/light-speak/lighthouse/env"
	"github.com/light-speak/lighthouse/graphql/middleware"
	"github.com/light-speak/lighthouse/log"
	"net/http"
)

func StartServer() error {
	port := env.Getenv("PORT", "4000")
	c := Config{
		Resolvers: &Resolver{
			Db: db.GetDb(),
		},
	}
	SetDirective(&c.Directives)
	router := middleware.GetRouter()

	srv := handler.NewDefaultServer(NewExecutableSchema(c))

	router.Handle("/", playground.ApolloSandboxHandler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Info("成功运行， 点击进入 GraphQL playground ： http://localhost:%s/ ", port)
	return http.ListenAndServe(":"+port, Middleware(db.GetDb(), router))
}
