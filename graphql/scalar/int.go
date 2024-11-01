package scalar

import (
	"fmt"
	"strconv"

	"github.com/light-speak/lighthouse/errors"
)

type IntScalar struct{}

func (i *IntScalar) ParseValue(v interface{}, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
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
	default:
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("invalid integer value: %v, got %T", v, v),
			Locations: []*errors.GraphqlLocation{location},
		}
	}
}

func (i *IntScalar) Serialize(v interface{}, location *errors.GraphqlLocation) (string, errors.GraphqlErrorInterface) {
	switch v := v.(type) {
	case int64:
		return strconv.FormatInt(v, 10), nil
	default:
		return "", &errors.GraphQLError{
			Message:   fmt.Sprintf("value is not an integer: %v", v),
			Locations: []*errors.GraphqlLocation{location},
		}
	}
}

func (i *IntScalar) ParseLiteral(v interface{}, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
	switch v := v.(type) {
	case int64:
		return v, nil
	}
	return nil, &errors.GraphQLError{
		Message:   fmt.Sprintf("invalid literal for Int: %v", v),
		Locations: []*errors.GraphqlLocation{location},
	}
}

func (i *IntScalar) GoType() string {
	return "int64"
}
