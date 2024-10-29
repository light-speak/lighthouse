package scalar

import (
	"fmt"
	"strconv"

	"github.com/light-speak/lighthouse/errors"
)

type IDScalar struct{}

func (i *IDScalar) ParseValue(v string, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
	intValue, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("invalid integer value: %s", v),
			Locations: []*errors.GraphqlLocation{location},
		}
	}
	return intValue, nil
}

func (i *IDScalar) Serialize(v interface{}, location *errors.GraphqlLocation) (string, errors.GraphqlErrorInterface) {
	if intValue, ok := v.(int64); ok {
		return strconv.FormatInt(intValue, 10), nil
	}
	return "", &errors.GraphQLError{
		Message:   fmt.Sprintf("value is not an integer: %v", v),
		Locations: []*errors.GraphqlLocation{location},
	}
}

func (i *IDScalar) ParseLiteral(v interface{}, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
	if vt, ok := v.(int64); ok {
		return vt, nil
	}
	return nil, &errors.GraphQLError{
		Message:   fmt.Sprintf("invalid literal for Int: %v", v),
		Locations: []*errors.GraphqlLocation{location},
	}
}

func (i *IDScalar) GoType() string {
	return "int64"
}
