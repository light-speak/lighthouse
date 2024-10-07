package ast

type Node interface {
	GetName() string
	GetType() NodeType
	GetDescription() string
	GetImplements() []string
	GetFields() []FieldNode
	GetDirectives() []DirectiveNode
	GetArgs() []ArgumentNode
	IsDeprecated() bool
	GetDeprecationReason() string
	IsNonNull() bool
	IsList() bool
	HasField(name string) bool
	HasDirective(name string) bool
	GetDirective(name string) *DirectiveNode
	GetParent() Node
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
	Value    Value
	Type     *FieldType
	Children []*ArgumentValue
}
