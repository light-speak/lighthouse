package net

import (
	"encoding/json"
	"net/http"

	"github.com/light-speak/lighthouse/graphql"
	"github.com/light-speak/lighthouse/log"
)

type GraphQLRequest struct {
	Query     string          `json:"query"`
	Variables json.RawMessage `json:"variables"`
}

type GraphQLResponse struct {
	Data   interface{}     `json:"data"`
	Errors []*GraphQLError `json:"errors"`
}

type GraphQLError struct {
	Message   string             `json:"message"`
	Locations []*GraphQLLocation `json:"locations"`
}

type GraphQLLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func graphQLHandler(w http.ResponseWriter, r *http.Request) {
	var request GraphQLRequest

	switch r.Method {
	case http.MethodPost:
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Error().Msgf("Failed to decode request body: %v", err)
			sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
			return
		}
	case http.MethodGet:
		request.Query = r.URL.Query().Get("query")
		if request.Query == "" {
			log.Error().Msg("Query is required")
			sendErrorResponse(w, "Query is required", http.StatusBadRequest)
			return
		}
		variables := r.URL.Query().Get("variables")
		if variables != "" {
			if err := json.Unmarshal([]byte(variables), &request.Variables); err != nil {
				log.Error().Msgf("Failed to parse variables: %v", err)
				sendErrorResponse(w, "Invalid variables", http.StatusBadRequest)
				return
			}
		}
	default:
		log.Error().Msgf("Method not allowed: %s", r.Method)
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := graphql.ExecuteQuery(request.Query, request.Variables)
	response := GraphQLResponse{
		Data: data,
	}
	if err != nil {
		log.Error().Msgf("Error executing query: %v", err)
		response.Errors = []*GraphQLError{
			{
				Message: err.Error(),
			},
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error().Msgf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := GraphQLResponse{
		Errors: []*GraphQLError{
			{
				Message: message,
			},
		},
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error().Msgf("Failed to encode error response: %v", err)
	}
}
