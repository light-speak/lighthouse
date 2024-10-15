package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/light-speak/lighthouse/env"
	"github.com/light-speak/lighthouse/log"
)

type GraphQLRequest struct {
	Query     string          `json:"query"`
	Variables json.RawMessage `json:"variables"`
}

type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []GraphQLError  `json:"errors"`
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
	var request GraphQLRequest
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

	response := GraphQLResponse{
		Data:   json.RawMessage("{}"),
		Errors: []GraphQLError{},
	}

	log.Info().Msgf("Received GraphQL request: %s", request.Query)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func StartService() {
	log.Info().Msgf("Starting GraphQL service")
	port := env.GetEnv("PORT", "8000")
	r := chi.NewRouter()

	r.Post("/query", handleGraphQL)

	log.Info().Msgf("GraphQL service started on port %s, You can access the service at http://localhost:%s/query", port, port)
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), r)
}
