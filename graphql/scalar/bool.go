package scalar

import (
	"fmt"
	"strconv"

	"github.com/light-speak/lighthouse/errors"
)

type BooleanScalar struct{}

func (i *BooleanScalar) ParseValue(v interface{}, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
	switch v := v.(type) {
	case string:
		boolValue, err := strconv.ParseBool(v)
		if err != nil {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("invalid boolean value: %s", v),
				Locations: []*errors.GraphqlLocation{location},
			}
		}
		return boolValue, nil
	case bool:
		return v, nil
	default:
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("invalid boolean value: %v", v),
			Locations: []*errors.GraphqlLocation{location},
		}
	}
}

func (i *BooleanScalar) Serialize(v interface{}, location *errors.GraphqlLocation) (string, errors.GraphqlErrorInterface) {
	switch v := v.(type) {
	case bool:
		return strconv.FormatBool(v), nil
	default:
		return "", &errors.GraphQLError{
			Message:   fmt.Sprintf("value is not a boolean: %v", v),
			Locations: []*errors.GraphqlLocation{location},
		}
	}
}

func (i *BooleanScalar) ParseLiteral(v interface{}, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
	switch v := v.(type) {
	case bool:
		return v, nil
	}
	return nil, &errors.GraphQLError{
		Message:   fmt.Sprintf("invalid literal for Boolean: %v", v),
		Locations: []*errors.GraphqlLocation{location},
	}
}

func (i *BooleanScalar) GoType() string {
	return "bool"
}
