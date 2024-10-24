package scalar

import (
	"fmt"
	"strconv"

	"github.com/light-speak/lighthouse/graphql/ast"
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

func (f *FloatScalar) ParseLiteral(v ast.Value) (interface{}, error) {
	if vt, ok := v.(*ast.FloatValue); ok {
		return vt.Value, nil
	}
	return nil, fmt.Errorf("invalid literal for Float: %v", v)
}

func (f *FloatScalar) GoType() string {
	return "float64"
}
