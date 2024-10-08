package scalar

import (
	"fmt"
	"strconv"

	"github.com/light-speak/lighthouse/parser/value"
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
	return "", fmt.Errorf("value is not an boolean: %v", v)
}

func (i *BooleanScalar) ParseLiteral(v value.Value) (interface{}, error) {
	switch vt := v.(type) {
	case *value.BooleanValue:
		return vt.Value, nil
	default:
		return nil, fmt.Errorf("invalid literal for Boolean: %v", v)
	}
}

func init() {
	
}
