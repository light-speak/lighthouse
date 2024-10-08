package ast

import "github.com/light-speak/lighthouse/parser/value"

type EnumNode struct {
	Name        string
	Values      []*EnumValueNode
	Description string
	Directives  []*DirectiveNode
}

func (e *EnumNode) GetName() string {
	return e.Name
}

func (e *EnumNode) GetType() NodeType {
	return NodeTypeEnum
}

func (e *EnumNode) GetDescription() string {
	return e.Description
}

func (e *EnumNode) IsDeprecated() (bool, string) {
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

func (e *EnumNode) GetField(name string) *FieldNode {
	return nil
}

func (e *EnumNode) GetDirective(name string) *DirectiveNode {
	for _, directive := range e.Directives {
		if directive.Name == name {
			return directive
		}
	}
	return nil
}

func (e *EnumNode) GetArg(name string) *ArgumentNode {
	return nil
}

func (e *EnumNode) GetParent() Node {
	return nil
}

func (e *EnumNode) GetDirectives() []*DirectiveNode {
	return e.Directives
}

func (e *EnumNode) GetArgs() []*ArgumentNode {
	return nil
}
