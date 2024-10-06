package ast

type FragmentNode struct {
	Description   string
	Name          string
	OnType        string
	OperationType *TypeNode
	Directives    []DirectiveNode
	Fields        []FieldNode
}

func (f *FragmentNode) GetName() string {
	return f.Name
}

func (f *FragmentNode) GetType() NodeType {
	return NodeTypeFragment
}

func (f *FragmentNode) GetDescription() string {
	return f.Description
}

func (f *FragmentNode) GetImplements() []string {
	return []string{}
}

func (f *FragmentNode) GetFields() []FieldNode {
	return f.Fields
}

func (f *FragmentNode) GetDirectives() []DirectiveNode {
	return f.Directives
}

func (f *FragmentNode) GetArgs() []ArgumentNode {
	return []ArgumentNode{}
}

func (f *FragmentNode) IsDeprecated() bool {
	return false
}

func (f *FragmentNode) GetDeprecationReason() string {
	return ""
}

func (f *FragmentNode) IsNonNull() bool {
	return true
}

func (f *FragmentNode) IsList() bool {
	return false
}

func (f *FragmentNode) GetElemType() *FieldType {
	return nil
}

func (f *FragmentNode) GetDefaultValue() string {
	return ""
}

func (f *FragmentNode) HasField(name string) bool {
	for _, field := range f.Fields {
		if field.Name == name {
			return true
		}
	}
	return false
}

func (f *FragmentNode) HasDirective(name string) bool {
	for _, directive := range f.Directives {
		if directive.Name == name {
			return true
		}
	}
	return false
}

func (f *FragmentNode) GetDirective(name string) *DirectiveNode {
	for _, directive := range f.Directives {
		if directive.Name == name {
			return &directive
		}
	}
	return nil
}

func (f *FragmentNode) GetParent() Node {
	return nil
}
