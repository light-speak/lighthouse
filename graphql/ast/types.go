package ast

// TypeNode represents a GraphQL type
type TypeNode struct {
	BaseNode
	IsResponse     bool
	Implements     []string
	ImplementTypes []Node
	Fields         []*FieldNode
	Go             string
}

func (t *TypeNode) GetNodeType() NodeType           { return NodeTypeType }
func (t *TypeNode) GetFields() []*FieldNode         { return t.Fields }
func (t *TypeNode) GetField(name string) *FieldNode { return findField(t.Fields, name) }
func (t *TypeNode) GoType() string {
	if t.Go != "" {
		return t.Go
	}
	return t.GetName()
}

// InputNode represents a GraphQL input type
type InputNode struct {
	BaseNode
	Fields []*FieldNode
}

func (i *InputNode) GetNodeType() NodeType           { return NodeTypeInput }
func (i *InputNode) GetFields() []*FieldNode         { return i.Fields }
func (i *InputNode) GetField(name string) *FieldNode { return findField(i.Fields, name) }

// InterfaceNode represents a GraphQL interface
type InterfaceNode struct {
	BaseNode
	Fields []*FieldNode
}

func (i *InterfaceNode) GetNodeType() NodeType           { return NodeTypeInterface }
func (i *InterfaceNode) GetFields() []*FieldNode         { return i.Fields }
func (i *InterfaceNode) GetField(name string) *FieldNode { return findField(i.Fields, name) }

// UnionNode represents a GraphQL union
type UnionNode struct {
	BaseNode
	Types     []string
	TypeNodes []Node
}

func (u *UnionNode) GetNodeType() NodeType { return NodeTypeUnion }

// ScalarNode represents a GraphQL scalar
type ScalarNode struct {
	BaseNode
	Scalar ScalarType
}

func (s *ScalarNode) GetNodeType() NodeType { return NodeTypeScalar }

type EnumNode struct {
	BaseNode
	Values []*EnumValueNode
}

func (e *EnumNode) GetNodeType() NodeType { return NodeTypeEnum }

type OperationType string

const (
	QueryOperation        OperationType = "query"
	MutationOperation     OperationType = "mutation"
	SubscriptionOperation OperationType = "subscription"
)

type OperationNode struct {
	BaseNode
	Type   OperationType
	Args   []*ArgumentNode
	Fields []*FieldNode
}

func (t *OperationNode) GetNodeType() NodeType           { return NodeTypeOperation }
func (t *OperationNode) GetArgs() []*ArgumentNode        { return t.Args }
func (t *OperationNode) GetFields() []*FieldNode         { return t.Fields }
func (t *OperationNode) GetField(name string) *FieldNode { return findField(t.Fields, name) }
