package ast

import "github.com/light-speak/lighthouse/parser/value"

type EnumValueNode struct {
	Name        string
	Description string
	Directives  []*DirectiveNode
	Parent      Node
}

func (e *EnumValueNode) GetName() string {
	return e.Name
}

func (e *EnumValueNode) GetType() NodeType {
	return NodeTypeEnumValue
}

func (e *EnumValueNode) GetDescription() string {
	return e.Description
}

func (e *EnumValueNode) IsDeprecated() (bool, string) {
	directive := e.GetDirective("deprecated")
	if directive == nil {
		return false, ""
	}
	reason, err := value.ExtractValue(directive.GetArg("reason").Value.Value)
	if err != nil {
		return false, ""
	}
	return true, reason.(string)
}

func (e *EnumValueNode) GetField(name string) *FieldNode {
	return nil
}

func (e *EnumValueNode) GetDirective(name string) *DirectiveNode {
	for _, directive := range e.Directives {
		if directive.Name == name {
			return directive
		}
	}
	return nil
}

func (e *EnumValueNode) GetArg(name string) *ArgumentNode {
	return nil
}

func (e *EnumValueNode) GetParent() Node {
	return e.Parent
}

func (e *EnumValueNode) GetDirectives() []*DirectiveNode {
	return e.Directives
}

func (e *EnumValueNode) GetArgs() []*ArgumentNode {
	return nil
}
