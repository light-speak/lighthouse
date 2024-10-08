package ast

import "github.com/light-speak/lighthouse/parser/value"

type Node interface {
	GetName() string
	GetType() NodeType
	GetDescription() string
	IsDeprecated() (bool, string)
	GetField(name string) *FieldNode
	GetDirective(name string) *DirectiveNode
	GetArg(name string) *ArgumentNode
	GetParent() Node

	GetDirectives() []*DirectiveNode
	GetArgs() []*ArgumentNode
	GetFields() []*FieldNode
}

type NodeType string

const (
	NodeTypeType                NodeType = "Type"
	NodeTypeField               NodeType = "Field"
	NodeTypeArgument            NodeType = "Argument"
	NodeTypeDirective           NodeType = "Directive"
	NodeTypeDirectiveDefinition NodeType = "DirectiveDefinition"
	NodeTypeScalar              NodeType = "Scalar"
	NodeTypeUnion               NodeType = "Union"
	NodeTypeEnum                NodeType = "Enum"
	NodeTypeInterface           NodeType = "Interface"
	NodeTypeInput               NodeType = "Input"
	NodeTypeEnumValue           NodeType = "EnumValue"
	NodeTypeFragment            NodeType = "Fragment"
)

type OperationType string

const (
	OperationTypeQuery        OperationType = "Query"
	OperationTypeMutation     OperationType = "Mutation"
	OperationTypeSubscription OperationType = "Subscription"
	OperationTypeEntity       OperationType = "Entity"
)

type TypeCategory string

const (
	TypeCategoryScalar TypeCategory = "Scalar"
	TypeCategoryEnum   TypeCategory = "Enum"
	TypeCategoryInput  TypeCategory = "Input"
	TypeCategoryUnion  TypeCategory = "Union"
	TypeCategoryType   TypeCategory = "Type"
)

type FieldType struct {
	Name string
	Type Node

	TypeCategory TypeCategory

	IsList   bool
	ElemType *FieldType

	IsNonNull bool
}

type ArgumentValue struct {
	Value    value.Value
	Type     *FieldType
	Children []*ArgumentValue
}
