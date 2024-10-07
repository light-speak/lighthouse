package ast

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
