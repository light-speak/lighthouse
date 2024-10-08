package ast

import "github.com/light-speak/lighthouse/parser/value"

type TypeNode struct {
	Name           string
	Implements     []string
	ImplementTypes []*TypeNode
	Fields         []*FieldNode
	Description    string
	OperationType  OperationType
	Directives     []*DirectiveNode
}

func (t *TypeNode) GetName() string {
	return t.Name
}

func (t *TypeNode) GetType() NodeType {
	return NodeTypeType
}

func (t *TypeNode) GetDescription() string {
	return t.Description
}

func (t *TypeNode) IsDeprecated() (bool, string) {
	directive := t.GetDirective("deprecated")
	if directive == nil {
		return false, ""
	}
	reason, err := value.ExtractValue(directive.GetArg("reason").Value.Value)
	if err != nil {
		return false, ""
	}
	return true, reason.(string)
}

func (t *TypeNode) GetField(name string) *FieldNode {
	for _, field := range t.Fields {
		if field.Name == name {
			return field
		}
	}
	return nil
}

func (t *TypeNode) GetDirective(name string) *DirectiveNode {
	for _, directive := range t.Directives {
		if directive.Name == name {
			return directive
		}
	}
	return nil
}

func (t *TypeNode) GetArg(name string) *ArgumentNode {
	return nil
}

func (t *TypeNode) GetParent() Node {
	return nil
}

func (t *TypeNode) GetDirectives() []*DirectiveNode {
	return t.Directives
}

func (t *TypeNode) GetArgs() []*ArgumentNode {
	return nil
}

func (t *TypeNode) GetFields() []*FieldNode {
	return t.Fields
}
