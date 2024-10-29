package errors

import "fmt"

type GraphQLError struct {
	Message   string             `json:"message"`
	Locations []*GraphqlLocation `json:"locations,omitempty"`
}

type GraphqlLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func (g *GraphQLError) GraphqlError() *GraphQLError {
	return g
}

func (e *GraphQLError) Error() string {
	return fmt.Sprintf("graphql error: %s", e.Message)
}

type GraphqlErrorInterface interface {
	GraphqlError() *GraphQLError
	Error() string
}
