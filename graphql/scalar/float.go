package scalar

import (
	"fmt"
	"strconv"
)

type FloatScalar struct{}

func (f *FloatScalar) ParseValue(v string) (interface{}, error) {
	floatValue, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid float value: %s", v)
	}
	return floatValue, nil
}

func (f *FloatScalar) Serialize(v interface{}) (string, error) {
	if floatValue, ok := v.(float64); ok {
		return strconv.FormatFloat(floatValue, 'f', -1, 64), nil
	}
	return "", fmt.Errorf("value is not a float: %v", v)
}

func (f *FloatScalar) ParseLiteral(v interface{}) (interface{}, error) {
	if vt, ok := v.(float64); ok {
		return vt, nil
	}
	return nil, fmt.Errorf("invalid literal for Float: %v", v)
}

func (f *FloatScalar) GoType() string {
	return "float64"
}
