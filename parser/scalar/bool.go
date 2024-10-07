package scalar

import (
	"fmt"
	"strconv"

	"github.com/light-speak/lighthouse/parser/ast"
)

type BooleanScalar struct{}

func (i *BooleanScalar) ParseValue(value string) (interface{}, error) {
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf("invalid boolean value: %s", value)
	}
	return boolValue, nil
}

func (i *BooleanScalar) Serialize(value interface{}) (string, error) {
	if boolValue, ok := value.(bool); ok {
		return strconv.FormatBool(boolValue), nil
	}
	return "", fmt.Errorf("value is not an boolean: %v", value)
}

func (i *BooleanScalar) ParseLiteral(value ast.Value) (interface{}, error) {
	switch v := value.(type) {
	case *ast.BooleanValue:
		return v.Value, nil
	default:
		return nil, fmt.Errorf("invalid literal for Boolean: %v", value)
	}
}
