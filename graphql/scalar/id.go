package scalar

import (
	"fmt"
	"strconv"

	"github.com/light-speak/lighthouse/errors"
)

type IDScalar struct{}

func (i *IDScalar) ParseValue(v interface{}, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
	switch v := v.(type) {
	case string:
		intValue, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("invalid integer value: %s, got %T", v, v),
				Locations: []*errors.GraphqlLocation{location},
			}
		}
		return intValue, nil
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case int:
		return int64(v), nil
	default:
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("invalid integer value: %v got %T", v, v),
			Locations: []*errors.GraphqlLocation{location},
		}
	}
}

func (i *IDScalar) Serialize(v interface{}, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
	switch v := v.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		intValue, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("value is not an integer: %v, got %T", v, v),
				Locations: []*errors.GraphqlLocation{location},
			}
		}
		return intValue, nil
	}
	return nil, &errors.GraphQLError{
		Message:   fmt.Sprintf("value is not an integer: %v, got %T", v, v),
		Locations: []*errors.GraphqlLocation{location},
	}
}

func (i *IDScalar) ParseLiteral(v interface{}, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
	switch v := v.(type) {
	case int64:
		return v, nil
	case string:
		intValue, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("invalid integer value: %s, got %T", v, v),
				Locations: []*errors.GraphqlLocation{location},
			}
		}
		return intValue, nil
	case float64:
		return int64(v), nil
	case int:
		return int64(v), nil
	}
	return nil, &errors.GraphQLError{
		Message:   fmt.Sprintf("invalid literal for Int: %v, got %T", v, v),
		Locations: []*errors.GraphqlLocation{location},
	}
}

func (i *IDScalar) GoType() string {
	return "int64"
}
