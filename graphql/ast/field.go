package ast

type FieldNode struct {
	BaseNode
	Type     *FieldType
	Value    Value
	Args     []*ArgumentNode
	Parent   Node
	Children []*FieldNode
	Fragment *FragmentNode
}

func (f *FieldNode) GetNodeType() NodeType { return NodeTypeField }
func (f *FieldNode) GetParent() Node       { return f.Parent }

// EnumValueNode represents a GraphQL enum value
type EnumValueNode struct {
	BaseNode
	Value  int8
	Parent Node
}

func (e *EnumValueNode) GetNodeType() NodeType { return NodeTypeEnumValue }
func (e *EnumValueNode) GetParent() Node       { return e.Parent }

func findField(fields []*FieldNode, name string) *FieldNode {
	for _, field := range fields {
		if field.Name == name {
			return field
		}
	}
	return nil
}

func findArg(args []*ArgumentNode, name string) *ArgumentNode {
	for _, arg := range args {
		if arg.Name == name {
			return arg
		}
	}
	return nil
}

type ArgumentNode struct {
	BaseNode
	Type         *FieldType
	Value        *ArgumentValue
	DefaultValue *ArgumentValue
	Parent       Node
}

func (a *ArgumentNode) GetNodeType() NodeType { return NodeTypeArgument }
func (a *ArgumentNode) GetParent() Node       { return a.Parent }
