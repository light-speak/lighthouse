package scalar

import (
	"fmt"
	"time"

	"github.com/light-speak/lighthouse/graphql/ast"
)

type DateTimeScalar struct{}

func (d *DateTimeScalar) ParseValue(v string) (interface{}, error) {
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return nil, fmt.Errorf("invalid datetime value: %s", v)
	}
	return t, nil
}

func (d *DateTimeScalar) Serialize(v interface{}) (string, error) {
	if t, ok := v.(time.Time); ok {
		return t.Format(time.RFC3339), nil
	}
	return "", fmt.Errorf("value is not a datetime: %v", v)
}

func (d *DateTimeScalar) ParseLiteral(v ast.Value) (interface{}, error) {
	if vt, ok := v.(*ast.StringValue); ok {
		return d.ParseValue(vt.Value)
	}
	return nil, fmt.Errorf("invalid literal for DateTime: %v", v)
}

func (d *DateTimeScalar) GoType() string {
	return "time.Time"
}