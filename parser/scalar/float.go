package scalar

import (
	"fmt"
	"strconv"

	"github.com/light-speak/lighthouse/parser/ast"
)

type FloatScalar struct{}

func (f *FloatScalar) ParseValue(value string) (interface{}, error) {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid float value: %s", value)
	}
	return floatValue, nil
}

func (f *FloatScalar) Serialize(value interface{}) (string, error) {
	if floatValue, ok := value.(float64); ok {
		return strconv.FormatFloat(floatValue, 'f', -1, 64), nil
	}
	return "", fmt.Errorf("value is not a float: %v", value)
}

func (f *FloatScalar) ParseLiteral(value ast.Value) (interface{}, error) {
	switch v := value.(type) {
	case *ast.FloatValue:
		return v.Value, nil
	default:
		return nil, fmt.Errorf("invalid literal for Float: %v", value)
	}
}
