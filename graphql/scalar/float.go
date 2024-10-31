package scalar

import (
	"fmt"
	"strconv"

	"github.com/light-speak/lighthouse/errors"
)

type FloatScalar struct{}

func (f *FloatScalar) ParseValue(v interface{}, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
	switch v := v.(type) {
	case string:
		floatValue, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("invalid float value: %s", v),
				Locations: []*errors.GraphqlLocation{location},
			}
		}
		return floatValue, nil
	case float64:
		return v, nil
	default:
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("invalid float value: %v", v),
			Locations: []*errors.GraphqlLocation{location},
		}
	}
}

func (f *FloatScalar) Serialize(v interface{}, location *errors.GraphqlLocation) (string, errors.GraphqlErrorInterface) {
	switch v := v.(type) {
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	default:
		return "", &errors.GraphQLError{
			Message:   fmt.Sprintf("value is not a float: %v", v),
			Locations: []*errors.GraphqlLocation{location},
		}
	}
}

func (f *FloatScalar) ParseLiteral(v interface{}, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
	switch v := v.(type) {
	case float64:
		return v, nil
	}
	return nil, &errors.GraphQLError{
		Message:   fmt.Sprintf("invalid literal for Float: %v", v),
		Locations: []*errors.GraphqlLocation{location},
	}
}

func (f *FloatScalar) GoType() string {
	return "float64"
}
