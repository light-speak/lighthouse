package lighterr

import (
	"context"
	"errors"
	"runtime"

	"github.com/99designs/gqlgen/graphql"
	"github.com/light-speak/lighthouse/logs"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gorm.io/gorm"
)

func ErrorPresenter(ctx context.Context, e error) *gqlerror.Error {
	err := graphql.DefaultErrorPresenter(ctx, e)
	var myErr *GraphQLError

	// Check if error is gorm.ErrRecordNotFound
	if errors.Is(e, gorm.ErrRecordNotFound) {
		return nil
	}

	// Check if error is our custom GraphQLError type
	logs.Error().Err(e).Msg("error presenter")
	if errors.As(e, &myErr) {

		ext := map[string]interface{}{
			"code": myErr.Code,
			"info": GetCodeInfo(myErr.Code),
		}

		if config.Env != EnvProduction {
			// Capture stack trace with proper formatting
			stackTrace := captureStackTrace(5)
			ext["stack"] = stackTrace

			// Convert error to string if it's not nil
			if myErr.Err != nil {
				ext["err"] = myErr.Err.Error()
			}
		}

		return &gqlerror.Error{
			Message:    myErr.Message,
			Extensions: ext,
		}
	}

	// Add stack trace to other errors in development mode
	if config.Env != EnvProduction {
		if err.Extensions == nil {
			err.Extensions = map[string]interface{}{}
		}
		err.Extensions["stack"] = captureStackTrace(3)
		err.Extensions["originalError"] = e.Error()
	}

	return err
}

func captureStackTrace(skip int) []map[string]interface{} {
	// Maximum call depth
	const maxDepth = 10
	var pcs [maxDepth]uintptr

	// Get program counter for call stack
	n := runtime.Callers(skip, pcs[:])

	frames := runtime.CallersFrames(pcs[:n])
	var stack []map[string]interface{}

	// Iterate through each call frame to generate readable stack trace
	for frame, more := frames.Next(); ; frame, more = frames.Next() {
		stack = append(stack, map[string]interface{}{
			"function": frame.Function,
			"file":     frame.File,
			"line":     frame.Line,
		})

		if !more {
			break
		}
	}

	return stack
}
