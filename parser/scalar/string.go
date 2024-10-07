package scalar

import (
	"fmt"

	"github.com/light-speak/lighthouse/parser/ast"
)

type StringScalar struct{}

func (s *StringScalar) ParseValue(value string) (interface{}, error) {
	return value, nil
}

func (s *StringScalar) Serialize(value interface{}) (string, error) {
	if stringValue, ok := value.(string); ok {
		return stringValue, nil
	}
	return "", fmt.Errorf("value is not a string: %v", value)
}

func (s *StringScalar) ParseLiteral(value ast.Value) (interface{}, error) {
	switch v := value.(type) {
	case *ast.StringValue:
		return v.Value, nil
	default:
		return nil, fmt.Errorf("invalid literal for String: %v", value)
	}
}
