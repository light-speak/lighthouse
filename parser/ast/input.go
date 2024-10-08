package ast

import "github.com/light-speak/lighthouse/parser/value"

type InputNode struct {
	Name        string
	Description string
	Fields      []*FieldNode
	Directives  []*DirectiveNode
}

func (i *InputNode) GetName() string {
	return i.Name
}

func (i *InputNode) GetType() NodeType {
	return NodeTypeInput
}

func (i *InputNode) GetDescription() string {
	return i.Description
}

func (i *InputNode) IsDeprecated() (bool, string) {
	directive := i.GetDirective("deprecated")
	if directive == nil {
		return false, ""
	}
	reason, err := value.ExtractValue(directive.GetArg("reason").Value.Value)
	if err != nil {
		return false, ""
	}
	return true, reason.(string)
}


func (i *InputNode) GetField(name string) *FieldNode {
	for _, field := range i.Fields {
		if field.Name == name {
			return field
		}
	}
	return nil
}


func (i *InputNode) GetDirective(name string) *DirectiveNode {
	for _, directive := range i.Directives {
		if directive.Name == name {
			return directive
		}
	}
	return nil
}

func (i *InputNode) GetArg(name string) *ArgumentNode {
	return nil
}

func (i *InputNode) GetParent() Node {
	return nil
}

func (i *InputNode) GetDirectives() []*DirectiveNode {
	return i.Directives
}

func (i *InputNode) GetArgs() []*ArgumentNode {
	return nil
}

func (i *InputNode) GetFields() []*FieldNode {
	return i.Fields
}
