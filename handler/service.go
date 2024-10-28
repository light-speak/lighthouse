package handler

import (
	"fmt"
	"net/http"

	"github.com/light-speak/lighthouse/env"
	"github.com/light-speak/lighthouse/graphql"
	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/net"
)

func StartService() {
	err := graphql.LoadSchema()
	if err != nil {
		log.Error().Msgf("Failed to load schema: %v", err)
		return
	}

	log.Info().Msgf("Starting GraphQL service")
	port := env.LighthouseConfig.Server.Port
	r := net.New()

	log.Info().Msgf("GraphQL service started on port %s, You can access the service at http://localhost:%s/query", port, port)
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), r)
	if err != nil {
		log.Error().Msgf("Failed to start GraphQL service: %v", err)
	}
}
