package scalar

import (
	"fmt"
	"time"

	"github.com/light-speak/lighthouse/errors"
)

type DateTimeScalar struct{}

func (d *DateTimeScalar) ParseValue(v interface{}, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
	switch v := v.(type) {
	case string:
		t, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("invalid datetime value: %s", v),
				Locations: []*errors.GraphqlLocation{location},
			}
		}
		return t, nil
	case time.Time:
		return v, nil
	default:
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("invalid datetime value: %v", v),
			Locations: []*errors.GraphqlLocation{location},
		}
	}
}

func (d *DateTimeScalar) Serialize(v interface{}, location *errors.GraphqlLocation) (string, errors.GraphqlErrorInterface) {
	switch t := v.(type) {
	case time.Time:
		return t.Format("2006-01-02 15:04:05"), nil
	case string:
		parsedTime, err := time.Parse("2006-01-02T15:04:05Z07:00", t)
		if err != nil {
			return "", &errors.GraphQLError{
				Message:   fmt.Sprintf("invalid datetime string: %s", t),
				Locations: []*errors.GraphqlLocation{location},
			}
		}
		return parsedTime.Format("2006-01-02 15:04:05"), nil
	default:
		return "", &errors.GraphQLError{
			Message:   fmt.Sprintf("unsupported value type: %v, type is %T", v, v),
			Locations: []*errors.GraphqlLocation{location},
		}
	}
}

func (d *DateTimeScalar) ParseLiteral(v interface{}, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface) {
	switch v := v.(type) {
	case string:
		return d.ParseValue(v, location)
	}
	return nil, &errors.GraphQLError{
		Message:   fmt.Sprintf("invalid literal for DateTime: %v", v),
		Locations: []*errors.GraphqlLocation{location},
	}
}

func (d *DateTimeScalar) GoType() string {
	return "time.Time"
}
