package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/light-speak/lighthouse/env"
	"github.com/light-speak/lighthouse/graphql"
	"github.com/light-speak/lighthouse/graphql/parser/lexer"
	"github.com/light-speak/lighthouse/log"
)

type GraphQLRequest struct {
	Query     string          `json:"query"`
	Variables json.RawMessage `json:"variables"`
}

type GraphQLResponse struct {
	Data   interface{}    `json:"data"`
	Errors []GraphQLError `json:"errors"`
}

type GraphQLError struct {
	Message   string            `json:"message"`
	Locations []GraphQLLocation `json:"locations"`
}

type GraphQLLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func handleGraphQL(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost && r.Method != http.MethodGet {
	// 	log.Error().Msgf("Method not allowed: %s", r.Method)
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }

	var request GraphQLRequest
	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response := GraphQLResponse{
				Errors: []GraphQLError{
					{
						Message:   "Invalid request body",
						Locations: []GraphQLLocation{{Line: 1, Column: 1}},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	} else if r.Method == http.MethodGet {
		request.Query = r.URL.Query().Get("query")
		variables := r.URL.Query().Get("variables")
		if variables != "" {
			request.Variables = json.RawMessage(variables)
		}
	}

	// log.Info().Msgf("Received GraphQL request: %s", request.Query)
	p := graphql.GetParser()

	qp := p.NewQueryParser(lexer.NewLexer([]*lexer.Content{
		{
			Content: request.Query,
		},
	}))
	qp.ParseSchema()
	err := qp.Validate(p.NodeStore)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}
	res, err := graphql.ExecuteQuery(qp)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	response := GraphQLResponse{
		Data:   res,
		Errors: nil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func StartService() {
	_, err := graphql.ParserSchema([]string{"graphql/base.graphql", "graphql/demo.graphql"})
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	log.Info().Msgf("Starting GraphQL service")
	port := env.GetEnv("PORT", "8000")
	r := chi.NewRouter()

	r.HandleFunc("/query", handleGraphQL)

	log.Info().Msgf("GraphQL service started on port %s, You can access the service at http://localhost:%s/query", port, port)
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), r)
}
