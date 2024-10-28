package scalar

import (
	"fmt"
	"time"
)

type DateTimeScalar struct{}

func (d *DateTimeScalar) ParseValue(v string) (interface{}, error) {
	t, err := time.Parse("2006-01-02 15:04:05", v)
	if err != nil {
		return nil, fmt.Errorf("invalid datetime value: %s", v)
	}
	return t, nil
}

func (d *DateTimeScalar) Serialize(v interface{}) (string, error) {
	switch t := v.(type) {
	case time.Time:
		return t.Format("2006-01-02 15:04:05"), nil
	case string:
		parsedTime, err := time.Parse("2006-01-02T15:04:05Z07:00", t)
		if err != nil {
			return "", fmt.Errorf("invalid datetime string: %s", t)
		}
		return parsedTime.Format("2006-01-02 15:04:05"), nil
	default:
		return "", fmt.Errorf("unsupported value type: %v, type is %T", v, v)
	}
}

func (d *DateTimeScalar) ParseLiteral(v interface{}) (interface{}, error) {
	if vt, ok := v.(string); ok {
		return d.ParseValue(vt)
	}
	return nil, fmt.Errorf("invalid literal for DateTime: %v", v)
}

func (d *DateTimeScalar) GoType() string {
	return "time.Time"
}
