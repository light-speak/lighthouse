package scalar

import (
	"fmt"
	"strconv"

	"github.com/light-speak/lighthouse/parser/ast"
)

type IntScalar struct{}

func (i *IntScalar) ParseValue(value string) (interface{}, error) {
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid integer value: %s", value)
	}
	return intValue, nil
}

func (i *IntScalar) Serialize(value interface{}) (string, error) {
	if intValue, ok := value.(int64); ok {
		return strconv.FormatInt(intValue, 10), nil
	}
	return "", fmt.Errorf("value is not an integer: %v", value)
}

func (i *IntScalar) ParseLiteral(value ast.Value) (interface{}, error) {
	switch v := value.(type) {
	case *ast.IntValue:
		return v.Value, nil
	default:
		return nil, fmt.Errorf("invalid literal for Int: %v", value)
	}
}
