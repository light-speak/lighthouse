package net

import (
	"encoding/json"
	"net/http"

	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/excute"
	"github.com/light-speak/lighthouse/log"
)

type GraphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

type GraphQLResponse struct {
	Data   interface{}            `json:"data"`
	Errors []*errors.GraphQLError `json:"errors,omitempty"`
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

	data, err := excute.ExecuteQuery(request.Query, request.Variables)
	response := GraphQLResponse{
		Data: data,
	}
	if err != nil {
		log.Error().Msgf("Error executing query: %v", err)
		response.Errors = []*errors.GraphQLError{
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
		Errors: []*errors.GraphQLError{
			{
				Message: message,
			},
		},
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error().Msgf("Failed to encode error response: %v", err)
	}
}
