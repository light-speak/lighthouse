package scalar

import (
	"fmt"
	"strconv"
)

type IDScalar struct{}

func (i *IDScalar) ParseValue(v string) (interface{}, error) {
	intValue, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid integer value: %s", v)
	}
	return intValue, nil
}

func (i *IDScalar) Serialize(v interface{}) (string, error) {
	if intValue, ok := v.(int64); ok {
		return strconv.FormatInt(intValue, 10), nil
	}
	return "", fmt.Errorf("value is not an integer: %v", v)
}

func (i *IDScalar) ParseLiteral(v interface{}) (interface{}, error) {
	if vt, ok := v.(int64); ok {
		return vt, nil
	}
	return nil, fmt.Errorf("invalid literal for Int: %v", v)
}

func (i *IDScalar) GoType() string {
	return "int64"
}
