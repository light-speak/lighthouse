package scalar

import (
	"fmt"
	"strconv"

	"github.com/light-speak/lighthouse/errors"
)

type BooleanScalar struct{}

func (i *BooleanScalar) ParseValue(v string, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
	boolValue, err := strconv.ParseBool(v)
	if err != nil {
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("invalid boolean value: %s", v),
			Locations: []*errors.GraphqlLocation{location},
		}
	}
	return boolValue, nil
}

func (i *BooleanScalar) Serialize(v interface{}, location *errors.GraphqlLocation) (string, errors.GraphqlErrorInterface) {
	if boolValue, ok := v.(bool); ok {
		return strconv.FormatBool(boolValue), nil
	}
	return "", &errors.GraphQLError{
		Message:   fmt.Sprintf("value is not a boolean: %v", v),
		Locations: []*errors.GraphqlLocation{location},
	}
}

func (i *BooleanScalar) ParseLiteral(v interface{}, location *errors.GraphqlLocation	) (interface{}, errors.GraphqlErrorInterface) {
	if vt, ok := v.(bool); ok {
		return vt, nil
	}
	return nil, &errors.GraphQLError{
		Message:   fmt.Sprintf("invalid literal for Boolean: %v", v),
		Locations: []*errors.GraphqlLocation{location},
	}
}

func (i *BooleanScalar) GoType() string {
	return "bool"
}

func init() {
	
}
