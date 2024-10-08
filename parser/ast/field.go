package ast

import "github.com/light-speak/lighthouse/parser/value"

type FieldNode struct {
	Name        string
	Type        *FieldType
	Description string
	Args        []*ArgumentNode
	Directives  []*DirectiveNode
	Parent      Node
}

func (f *FieldNode) GetName() string {
	return f.Name
}

func (f *FieldNode) GetType() NodeType {
	return NodeTypeField
}

func (f *FieldNode) GetDescription() string {
	return f.Description
}

func (f *FieldNode) IsDeprecated() (bool, string) {
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

func (f *FieldNode) GetField(name string) *FieldNode {
	return nil
}

func (f *FieldNode) GetDirective(name string) *DirectiveNode {
	for _, directive := range f.Directives {
		if directive.Name == name {
			return directive
		}
	}
	return nil
}

func (f *FieldNode) GetArg(name string) *ArgumentNode {
	for _, arg := range f.Args {
		if arg.Name == name {
			return arg
		}
	}
	return nil
}

func (f *FieldNode) GetParent() Node {
	return f.Parent
}

func (f *FieldNode) GetDirectives() []*DirectiveNode {
	return f.Directives
}

func (f *FieldNode) GetArgs() []*ArgumentNode {
	return f.Args
}

func (f *FieldNode) GetFields() []*FieldNode {
	return nil
}
