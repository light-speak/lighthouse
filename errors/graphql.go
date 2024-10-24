package errors

import "fmt"

type GraphQLError struct {
	Message   string            `json:"message"`
	Locations []GraphqlLocation `json:"locations"`
}

type GraphqlLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func (e *GraphQLError) Error() string {
	return fmt.Sprintf("graphql error: %s", e.Message)
}
