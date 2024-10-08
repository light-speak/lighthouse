package scalar

import (
	"fmt"
	"strconv"

	"github.com/light-speak/lighthouse/graphql/ast"
)

type IntScalar struct{}

func (i *IntScalar) ParseValue(v string) (interface{}, error) {
	intValue, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid integer value: %s", v)
	}
	return intValue, nil
}

func (i *IntScalar) Serialize(v interface{}) (string, error) {
	if intValue, ok := v.(int64); ok {
		return strconv.FormatInt(intValue, 10), nil
	}
	return "", fmt.Errorf("value is not an integer: %v", v)
}

func (i *IntScalar) ParseLiteral(v ast.Value) (interface{}, error) {
	switch vt := v.(type) {
	case *ast.IntValue:
		return vt.Value, nil
	default:
		return nil, fmt.Errorf("invalid literal for Int: %v", v)
	}
}
