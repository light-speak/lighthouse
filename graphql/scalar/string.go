package scalar

import (
	"fmt"
	"strconv"

	"github.com/light-speak/lighthouse/errors"
)

type StringScalar struct{}

func (s *StringScalar) ParseValue(v interface{}, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
	switch v := v.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("invalid string value: %v", v),
			Locations: []*errors.GraphqlLocation{location},
		}
	}
}

func (s *StringScalar) Serialize(v interface{}, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
	switch v := v.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	default:
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("value is not a string: %v", v),
			Locations: []*errors.GraphqlLocation{location},
		}
	}
}

func (s *StringScalar) ParseLiteral(v interface{}, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
	switch v := v.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	}

	return nil, &errors.GraphQLError{
		Message:   fmt.Sprintf("invalid literal for String: %v", v),
		Locations: []*errors.GraphqlLocation{location},
	}
}

func (s *StringScalar) GoType() string {
	return "string"
}
