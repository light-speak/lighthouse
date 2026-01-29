package lighterr

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/light-speak/lighthouse/logs"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gorm.io/gorm"
)

func ErrorPresenter(ctx context.Context, e error) *gqlerror.Error {
	err := graphql.DefaultErrorPresenter(ctx, e)
	var myErr *GraphQLError

	// Check if error is gorm.ErrRecordNotFound, convert to NotFound error
	if errors.Is(e, gorm.ErrRecordNotFound) {
		myErr = NewNotFoundError("record not found")
	}

	// Check if error is context deadline exceeded or canceled
	if myErr == nil && errors.Is(e, context.DeadlineExceeded) {
		myErr = NewRequestTimeoutError("request timeout: database query took too long", e)
	}
	if myErr == nil && errors.Is(e, context.Canceled) {
		myErr = NewRequestTimeoutError("request canceled", e)
	}

	// Check if error message contains common timeout patterns
	if myErr == nil && e != nil {
		errMsg := e.Error()
		if strings.Contains(errMsg, "context deadline exceeded") ||
			strings.Contains(errMsg, "context canceled") ||
			strings.Contains(errMsg, "i/o timeout") {
			myErr = NewRequestTimeoutError("request timeout", e)
		}
	}

	// Check if error is our custom GraphQLError type
	if myErr == nil {
		errors.As(e, &myErr)
	}

	if myErr != nil {
		// Capture stack trace for logging (always, for debugging)
		stackTrace := captureStackTrace(5)

		// Log based on error type: client errors (4xx) use Warn, server errors (5xx) use Error
		if IsClientError(myErr.Code) {
			logs.Warn().
				Str("key", GetCodeKey(myErr.Code)).
				Msg(myErr.Message)
		} else {
			// For server errors, log full details including stack
			logEvent := logs.Error().
				Str("key", GetCodeKey(myErr.Code)).
				Str("code_info", GetCodeInfo(myErr.Code))

			if myErr.Err != nil {
				logEvent = logEvent.Str("cause", myErr.Err.Error())
			}

			// Log stack trace as structured data
			logEvent.Interface("stack", stackTrace).Msg(myErr.Message)
		}

		ext := map[string]interface{}{
			"code": myErr.Code,
			"key":  GetCodeKey(myErr.Code),
		}

		// 非生产环境：总是返回详细信息
		// 生产环境：只有服务端错误（非客户端错误）才返回详细信息
		shouldIncludeDetails := config.Env != EnvProduction || !IsClientError(myErr.Code)

		if shouldIncludeDetails {
			ext["info"] = GetCodeInfo(myErr.Code)
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

	// Unknown error - capture full context for logging
	stackTrace := captureStackTrace(3)
	logs.Error().
		Err(e).
		Str("error_type", getErrorType(e)).
		Interface("stack", stackTrace).
		Msg("unhandled error")

	// 未知错误视为服务端错误，总是返回详细信息帮助调试
	if err.Extensions == nil {
		err.Extensions = map[string]interface{}{}
	}
	err.Extensions["stack"] = stackTrace
	err.Extensions["originalError"] = e.Error()
	err.Extensions["errorType"] = getErrorType(e)

	return err
}

// getErrorType returns a string representation of the error type
func getErrorType(e error) string {
	if e == nil {
		return "nil"
	}
	return fmt.Sprintf("%T", e)
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
