package ast

type ASTNode interface{}

type TypeNode struct {
	Name        string
	Fields      []FieldNode
	Description string
}

type FieldType string

const (
	FieldTypeString  FieldType = "String"
	FieldTypeInt     FieldType = "Int"
	FieldTypeFloat   FieldType = "Float"
	FieldTypeBoolean FieldType = "Boolean"
	FieldTypeID      FieldType = "ID"
)

type FieldNode struct {
	Name        string
	Type        FieldType
	Description string
	Args        []ArgumentNode
}

type InterfaceNode struct {
	Name        string
	Fields      []FieldNode
	Description string
}

type EnumNode struct {
	Name        string
	Values      []string
	Description string
}

type ExtendNode struct {
	Name        string
	Fields      []FieldNode
	Description string
}

type DirectiveNode struct {
	Name        string
	Args        []ArgumentNode
	Description string
}

type ArgumentNode struct {
	Name        string
	Type        FieldType
	Description string
}

type ScalarNode struct {
	Name        string
	Description string
}

type UnionNode struct {
	Name        string
	Types       []string
	Description string
}
