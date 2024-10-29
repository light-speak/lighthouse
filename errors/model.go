package errors

import "fmt"

type DataloaderError struct {
	Msg string
}

func (e *DataloaderError) Error() string {
	return e.Msg
}

func (e *DataloaderError) GraphqlError() *GraphQLError {
	return &GraphQLError{
		Message:   fmt.Sprintf("[dataloader error]: %s", e.Msg),
		Locations: []*GraphqlLocation{{Line: 1, Column: 1}},
	}
}
