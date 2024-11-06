package net

import (
	"encoding/json"
	"net/http"

	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/env"
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/excute"
	"github.com/light-speak/lighthouse/log"
)

var requestLimit = make(chan struct{}, env.LighthouseConfig.Server.Throttle/8)

type GraphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

type GraphQLResponse struct {
	Data   interface{}            `json:"data"`
	Errors []*errors.GraphQLError `json:"errors,omitempty"`
}

func graphQLHandler(w http.ResponseWriter, r *http.Request) {
	requestLimit <- struct{}{}
	defer func() {
		<-requestLimit
	}()

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
	ctx := r.Context().(*context.Context)
	data := excute.ExecuteQuery(ctx, request.Query, request.Variables)
	response := GraphQLResponse{
		Data: data,
	}
	if len(ctx.Errors) > 0 {
		response.Errors = make([]*errors.GraphQLError, 0)
		for _, err := range ctx.Errors {
			log.Error().Interface("ctx", ctx).Msgf("error: %v", err)
			e := err.GraphqlError()
			if env.LighthouseConfig.App.Environment != env.Development {
				e.Locations = nil
			}
			response.Errors = append(response.Errors, e)
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
