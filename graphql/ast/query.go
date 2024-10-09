package ast

type FragmentNode struct {
	BaseNode
	On     string
	Type   Node
	Fields []*FieldNode
}

func (f *FragmentNode) GetNodeType() NodeType   { return NodeTypeFragment }
func (f *FragmentNode) GetFields() []*FieldNode { return f.Fields }
func (f *FragmentNode) GetField(name string) *FieldNode {
	return findField(f.Fields, name)
}
