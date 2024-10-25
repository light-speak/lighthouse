package scalar

import (
	"fmt"
	"strconv"

	"github.com/light-speak/lighthouse/graphql/ast"
)

type BooleanScalar struct{}

func (i *BooleanScalar) ParseValue(v string) (interface{}, error) {
	boolValue, err := strconv.ParseBool(v)
	if err != nil {
		return nil, fmt.Errorf("invalid boolean value: %s", v)
	}
	return boolValue, nil
}

func (i *BooleanScalar) Serialize(v interface{}) (string, error) {
	if boolValue, ok := v.(bool); ok {
		return strconv.FormatBool(boolValue), nil
	}
	return "", fmt.Errorf("value is not a boolean: %v", v)
}

func (i *BooleanScalar) ParseLiteral(v ast.Value) (interface{}, error) {
	if vt, ok := v.(*ast.BooleanValue); ok {
		return vt.Value, nil
	}
	return nil, fmt.Errorf("invalid literal for Boolean: %v", v)
}

func (i *BooleanScalar) GoType() string {
	return "bool"
}

func init() {
	
}
