package own

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/bytedance/sonic"
	"github.com/light-speak/lighthouse/logs"
	"github.com/light-speak/lighthouse/routers/auth"
)

func OwnDirective(ctx context.Context, obj interface{}, next graphql.Resolver, field string) (interface{}, error) {
	// Get field context
	fieldCtx := graphql.GetFieldContext(ctx)
	if fieldCtx == nil {
		return nil, errors.New("failed to get field context")
	}

	fieldType := fieldCtx.Field.Definition.Type
	userId := auth.GetCtxUserId(ctx)

	// Return early if user not logged in
	if userId == 0 {
		if fieldType.NonNull {
			return "", nil
		}
		return nil, nil
	}

	// Get result from next resolver
	res, err := next(ctx)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	// Convert object to map
	var data map[string]interface{}
	jsonBytes, err := sonic.Marshal(obj)
	if err != nil {
		return nil, err
	}
	if err = sonic.Unmarshal(jsonBytes, &data); err != nil {
		return nil, err
	}

	// Get owner ID from field
	ownerId, ok := data[field]
	if !ok {
		if fieldType.NonNull {
			return "", nil
		}
		return nil, nil
	}

	// Convert owner ID to uint
	var ownId uint
	switch v := ownerId.(type) {
	case int64:
		ownId = uint(v)
	case float64:
		ownId = uint(v)
	default:
		logs.Error().Msgf("invalid owner ID type: expected int64/float64, got %T", v)
		if fieldType.NonNull {
			return "", nil
		}
		return nil, nil
	}

	// Check ownership
	if ownId != userId {
		if fieldType.NonNull {
			return "", nil
		}
		return nil, nil
	}

	return res, nil
}
