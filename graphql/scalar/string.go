package scalar

import (
	"fmt"
)

type StringScalar struct{}

func (s *StringScalar) ParseValue(v string) (interface{}, error) {
	return v, nil
}

func (s *StringScalar) Serialize(v interface{}) (string, error) {
	if stringValue, ok := v.(string); ok {
		return stringValue, nil
	}
	return "", fmt.Errorf("value is not a string: %v", v)
}

func (s *StringScalar) ParseLiteral(v interface{}) (interface{}, error) {
	if vt, ok := v.(string); ok {
		return vt, nil
	}
	return nil, fmt.Errorf("invalid literal for String: %v", v)
}

func (s *StringScalar) GoType() string {
	return "string"
}
