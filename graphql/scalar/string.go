package scalar

import (
	"fmt"

	"github.com/light-speak/lighthouse/errors"
)

type StringScalar struct{}

func (s *StringScalar) ParseValue(v string, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
	return v, nil
}

func (s *StringScalar) Serialize(v interface{}, location *errors.GraphqlLocation) (string, errors.GraphqlErrorInterface) {
	if stringValue, ok := v.(string); ok {
		return stringValue, nil
	}
	return "", &errors.GraphQLError{
		Message:   fmt.Sprintf("value is not a string: %v", v),
		Locations: []*errors.GraphqlLocation{location},
	}
}

func (s *StringScalar) ParseLiteral(v interface{}, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
	if vt, ok := v.(string); ok {
		return vt, nil
	}
	return nil, &errors.GraphQLError{
		Message:   fmt.Sprintf("invalid literal for String: %v", v),
		Locations: []*errors.GraphqlLocation{location},
	}
}

func (s *StringScalar) GoType() string {
	return "string"
}
