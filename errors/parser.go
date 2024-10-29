package errors

import (
	"fmt"
)

type LexerError struct {
	Path         *string
	Line         int
	LinePosition int
	Message      string
}

func (e *LexerError) Error() string {
	if e.Path != nil {
		return fmt.Sprintf("lexer error: %s, path: %s, line: %d, line position: %d", e.Message, *e.Path, e.Line, e.LinePosition)
	}
	return fmt.Sprintf("lexer error: %s, line: %d, line position: %d", e.Message, e.Line, e.LinePosition)
}

func (e *LexerError) GraphqlError() *GraphQLError {
	return &GraphQLError{
		Message:   fmt.Sprintf("[lexer error]: %s", e.Message),
		Locations: []*GraphqlLocation{{Line: e.Line, Column: e.LinePosition}},
	}
}

type ParserError struct {
	Message   string
	Locations *GraphqlLocation
}

func (e *ParserError) Error() string {
	return fmt.Sprintf("parser error: %s", e.Message)
}

func (e *ParserError) GraphqlError() *GraphQLError {
	return &GraphQLError{
		Message:   fmt.Sprintf("[parser error]: %s", e.Message),
		Locations: []*GraphqlLocation{e.Locations},
	}
}
