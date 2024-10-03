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
	GetElemType() *FieldType
	GetDefaultValue() string
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
)

type OperationType string

const (
	OperationTypeQuery        OperationType = "Query"
	OperationTypeMutation     OperationType = "Mutation"
	OperationTypeSubscription OperationType = "Subscription"
	OperationTypeEntity       OperationType = "Entity"
)

// TODO: 最后再确定一下是什么类型，先填入Name
type FieldType struct {
	Name      string
	IsScalar  bool
	IsEnum    bool
	IsUnion   bool
	IsInput   bool
	IsList    bool
	ElemType  *FieldType
	IsNonNull bool
}
