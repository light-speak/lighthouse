package ast

import (
	"github.com/light-speak/lighthouse/parser/value"
)

type ArgumentNode struct {
	Name         string
	Type         *FieldType
	Value        *ArgumentValue
	DefaultValue *ArgumentValue
	Description  string
	Directives   []*DirectiveNode
	Parent       Node
}

func (a *ArgumentNode) GetType() NodeType {
	return NodeTypeArgument
}

func (a *ArgumentNode) IsDeprecated() (bool, string) {
	directive := a.GetDirective("deprecated")
	if directive == nil {
		return false, ""
	}
	reason, err := value.ExtractValue(directive.GetArg("reason").Value.Value)
	if err != nil {
		return false, ""
	}
	return true, reason.(string)
}

func (a *ArgumentNode) GetField(name string) *FieldNode {
	return nil
}

func (a *ArgumentNode) GetDirective(name string) *DirectiveNode {
	for _, directive := range a.Directives {
		if directive.Name == name {
			return directive
		}
	}
	return nil
}

func (a *ArgumentNode) GetArg(name string) *ArgumentNode {
	return nil
}

func (a *ArgumentNode) GetParent() Node {
	return a.Parent
}

func (a *ArgumentNode) GetDirectives() []*DirectiveNode {
	return a.Directives
}

func (a *ArgumentNode) GetArgs() []*ArgumentNode {
	return nil
}
