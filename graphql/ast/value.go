package ast

import "fmt"

type Value interface {
	IsValue() bool
}

type IntValue struct {
	Value int64
}

func (v *IntValue) IsValue() bool {
	return true
}

type FloatValue struct {
	Value float64
}

func (v *FloatValue) IsValue() bool {
	return true
}

type StringValue struct {
	Value string
}

func (v *StringValue) IsValue() bool {
	return true
}

type BooleanValue struct {
	Value bool
}

func (v *BooleanValue) IsValue() bool {
	return true
}

type ListValue struct {
	Values []Value
}

func (v *ListValue) IsValue() bool {
	return true
}

type ObjectValue struct {
	Values map[string]Value
}

func (v *ObjectValue) IsValue() bool {
	return true
}

func ExtractValue(v Value) (interface{}, error) {
	switch v := v.(type) {
	case *IntValue:
		return v.Value, nil
	case *FloatValue:
		return v.Value, nil
	case *StringValue:
		return v.Value, nil
	case *BooleanValue:
		return v.Value, nil
	case *ListValue:
		values := make([]interface{}, len(v.Values))
		for i, value := range v.Values {
			extracted, err := ExtractValue(value)
			if err != nil {
				return nil, err
			}
			values[i] = extracted
		}
		return values, nil
	case *ObjectValue:
		values := make(map[string]interface{})
		for key, value := range v.Values {
			extracted, err := ExtractValue(value)
			if err != nil {
				return nil, err
			}
			values[key] = extracted
		}
		return values, nil
	default:
		return nil, fmt.Errorf("unsupported value type: %T", v)
	}
}
