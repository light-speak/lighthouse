package lighterr

type ErrorCode int

const (
	// 内部错误
	ErrorCodeInternalError ErrorCode = iota
	// 输入错误
	ErrorCodeInvalidInput
	// 未找到
	ErrorCodeNotFound
	// 未授权
	ErrorCodeUnauthorized
	// 禁止访问
	ErrorCodeForbidden
	// 错误请求
	ErrorCodeBadRequest
	// 数据冲突
	ErrorCodeConflict
	// 服务不可用
	ErrorCodeServiceUnavailable
	// 请求过于频繁
	ErrorCodeTooManyRequests
	// 请求超时
	ErrorCodeRequestTimeout
	// 数据库错误
	ErrorCodeDatabaseError
	// 第三方服务错误
	ErrorCodeThirdPartyError
	// 参数验证失败
	ErrorCodeValidationFailed
	// 资源已存在
	ErrorCodeResourceExists
	// 操作失败
	ErrorCodeOperationFailed
)

type GraphQLError struct {
	Message string    `json:"message"` // 错误信息
	Code    ErrorCode `json:"code"`    // 错误码
	Err     error     `json:"-"`       // 原始错误
}

// Error 实现 error 接口
func (e *GraphQLError) Error() string {
	return e.Message
}

// NewInternalError 创建内部错误
func NewInternalError(message string, err ...error) *GraphQLError {
	return NewGraphQLError(message, ErrorCodeInternalError, err...)
}

// NewInvalidInputError 创建输入错误
func NewInvalidInputError(message string, err ...error) *GraphQLError {
	return NewGraphQLError(message, ErrorCodeInvalidInput, err...)
}

// NewNotFoundError 创建未找到错误
func NewNotFoundError(message string, err ...error) *GraphQLError {
	return NewGraphQLError(message, ErrorCodeNotFound, err...)
}

// NewUnauthorizedError 创建未授权错误
func NewUnauthorizedError(message string, err ...error) *GraphQLError {
	return NewGraphQLError(message, ErrorCodeUnauthorized, err...)
}

// NewForbiddenError 创建禁止访问错误
func NewForbiddenError(message string, err ...error) *GraphQLError {
	return NewGraphQLError(message, ErrorCodeForbidden, err...)
}

// NewBadRequestError 创建错误请求错误
func NewBadRequestError(message string, err ...error) *GraphQLError {
	return NewGraphQLError(message, ErrorCodeBadRequest, err...)
}

// NewConflictError 创建数据冲突错误
func NewConflictError(message string, err ...error) *GraphQLError {
	return NewGraphQLError(message, ErrorCodeConflict, err...)
}

// NewServiceUnavailableError 创建服务不可用错误
func NewServiceUnavailableError(message string, err ...error) *GraphQLError {
	return NewGraphQLError(message, ErrorCodeServiceUnavailable, err...)
}

// NewTooManyRequestsError 创建请求过于频繁错误
func NewTooManyRequestsError(message string, err ...error) *GraphQLError {
	return NewGraphQLError(message, ErrorCodeTooManyRequests, err...)
}

// NewRequestTimeoutError 创建请求超时错误
func NewRequestTimeoutError(message string, err ...error) *GraphQLError {
	return NewGraphQLError(message, ErrorCodeRequestTimeout, err...)
}

// NewDatabaseError 创建数据库错误
func NewDatabaseError(message string, err ...error) *GraphQLError {
	return NewGraphQLError(message, ErrorCodeDatabaseError, err...)
}

// NewThirdPartyError 创建第三方服务错误
func NewThirdPartyError(message string, err ...error) *GraphQLError {
	return NewGraphQLError(message, ErrorCodeThirdPartyError, err...)
}

// NewValidationFailedError 创建参数验证失败错误
func NewValidationFailedError(message string, err ...error) *GraphQLError {
	return NewGraphQLError(message, ErrorCodeValidationFailed, err...)
}

// NewResourceExistsError 创建资源已存在错误
func NewResourceExistsError(message string, err ...error) *GraphQLError {
	return NewGraphQLError(message, ErrorCodeResourceExists, err...)
}

// NewOperationFailedError 创建操作失败错误
func NewOperationFailedError(message string, err ...error) *GraphQLError {
	return NewGraphQLError(message, ErrorCodeOperationFailed, err...)
}

// NewGraphQLError 创建新的 GraphQL 错误
func NewGraphQLError(message string, code ErrorCode, err ...error) *GraphQLError {
	if len(err) > 0 {
		return &GraphQLError{
			Message: message,
			Code:    code,
			Err:     err[0],
		}
	}
	return &GraphQLError{
		Message: message,
		Code:    code,
	}
}
