package ast

import "github.com/light-speak/lighthouse/parser/value"

type FragmentNode struct {
	Description string
	Name        string
	On          string
	Type        *TypeNode
	Directives  []*DirectiveNode
	Fields      []*FieldNode
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

func (f *FragmentNode) IsDeprecated() (bool, string) {
	directive := f.GetDirective("deprecated")
	if directive == nil {
		return false, ""
	}
	reason, err := value.ExtractValue(directive.GetArg("reason").Value.Value)
	if err != nil {
		return false, ""
	}
	return true, reason.(string)
}

func (f *FragmentNode) GetField(name string) *FieldNode {
	for _, field := range f.Fields {
		if field.Name == name {
			return field
		}
	}
	return nil
}

func (f *FragmentNode) GetDirective(name string) *DirectiveNode {
	for _, directive := range f.Directives {
		if directive.Name == name {
			return directive
		}
	}
	return nil
}

func (f *FragmentNode) GetArg(name string) *ArgumentNode {
	return nil
}

func (f *FragmentNode) GetParent() Node {
	return nil
}

func (f *FragmentNode) GetDirectives() []*DirectiveNode {
	return f.Directives
}

func (f *FragmentNode) GetArgs() []*ArgumentNode {
	return nil
}

func (f *FragmentNode) GetFields() []*FieldNode {
	return f.Fields
}
