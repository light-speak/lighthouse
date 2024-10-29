package scalar

import (
	"fmt"
	"strconv"

	"github.com/light-speak/lighthouse/errors"
)

type FloatScalar struct{}

func (f *FloatScalar) ParseValue(v string, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
	floatValue, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("invalid float value: %s", v),
			Locations: []*errors.GraphqlLocation{location},
		}
	}
	return floatValue, nil
}

func (f *FloatScalar) Serialize(v interface{}, location *errors.GraphqlLocation) (string, errors.GraphqlErrorInterface) {
	if floatValue, ok := v.(float64); ok {
		return strconv.FormatFloat(floatValue, 'f', -1, 64), nil
	}
	return "", &errors.GraphQLError{
		Message:   fmt.Sprintf("value is not a float: %v", v),
		Locations: []*errors.GraphqlLocation{location},
	}
}

func (f *FloatScalar) ParseLiteral(v interface{}, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
	if vt, ok := v.(float64); ok {
		return vt, nil
	}
	return nil, &errors.GraphQLError{
		Message:   fmt.Sprintf("invalid literal for Float: %v", v),
		Locations: []*errors.GraphqlLocation{location},
	}
}

func (f *FloatScalar) GoType() string {
	return "float64"
}
